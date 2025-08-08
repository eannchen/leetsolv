package core

import (
	"math"
	"time"

	"leetsolv/config"
	"leetsolv/internal/clock"
)

type Scheduler interface {
	ScheduleNewQuestion(id int, url, note string, grade Familiarity, importance Importance) *Question
	Schedule(q *Question, grade Familiarity)
	CalculatePriorityScore(q *Question) float64
}

// SM2Scheduler implements the spaced repetition scheduling logic
type SM2Scheduler struct {
	Clock clock.Clock
	// Base intervals for each importance level (in days)
	baseIntervals map[Importance]int
	// Maximum interval to prevent overly long gaps (in days)
	maxInterval int
	// Minimum and maximum ease factors
	minEaseFactor float64
	maxEaseFactor float64
	// Interval multipliers for importance levels
	intervalMultipliers map[Importance]float64
}

// NewSM2Scheduler creates a new scheduler with configured parameters
func NewSM2Scheduler(clock clock.Clock) *SM2Scheduler {
	return &SM2Scheduler{
		Clock: clock,
		baseIntervals: map[Importance]int{
			LowImportance:      8, // Faster growth, so start more spaced
			MediumImportance:   6, // Balanced
			HighImportance:     5, // Slightly tighter
			CriticalImportance: 4, // Tightest
		},
		maxInterval:   45,  // Ensure at least 2-3 reviews within 90 days
		minEaseFactor: 1.3, // Lower bound for ease factor
		maxEaseFactor: 2.6, // Upper bound to prevent overly long intervals
	}
}

func (s SM2Scheduler) ScheduleNewQuestion(id int, url, note string, grade Familiarity, importance Importance) *Question {
	today := s.Clock.Today()

	// Dynamic default EaseFactor based on importance
	startingEase := map[Importance]float64{
		LowImportance:      2.0,
		MediumImportance:   1.9,
		HighImportance:     1.8,
		CriticalImportance: 1.7,
	}[importance]

	q := &Question{
		ID:           id,
		URL:          url,
		Note:         note,
		Familiarity:  grade,
		Importance:   importance,
		EaseFactor:   startingEase,
		ReviewCount:  1,
		LastReviewed: today,
		CreatedAt:    s.Clock.Now(),
	}

	intervalDays := s.baseIntervals[importance]

	// Small tweaks to interval for early grading signal
	switch grade {
	case Easy:
		intervalDays += 2
	case VeryEasy:
		intervalDays += 5
	case VeryHard:
		intervalDays -= 1
	}

	s.setNextReview(q, today, intervalDays)
	return q
}

// Schedule updates the question's review schedule based on familiarity and importance
func (s SM2Scheduler) Schedule(q *Question, grade Familiarity) {
	q.ReviewCount++
	today := s.Clock.Today()

	baseInterval := s.baseIntervals[q.Importance]

	// Reset if still struggling
	if grade < Hard {
		s.setNextReview(q, today, baseInterval)
		s.setEaseFactorWithPenalty(q, grade)
		q.LastReviewed = today
		q.Familiarity = grade
		return
	}

	// Penalty for being overdue
	if config.Env().OverduePenalty {
		overdueLimit := config.Env().OverdueLimit
		overdueDays := int(today.Sub(q.NextReview).Hours() / 24)
		if overdueDays > overdueLimit && q.Importance > LowImportance && grade < VeryEasy {
			penaltyFactor := math.Min(float64(overdueDays-overdueLimit)*0.01, 0.1)
			q.EaseFactor -= penaltyFactor
		}
	}

	// Growth based on last interval × EaseFactor
	prevInterval := int(q.NextReview.Sub(q.LastReviewed).Hours() / 24)
	if prevInterval < 1 {
		prevInterval = baseInterval // fallback
	}

	intervalDays := int(math.Round(float64(prevInterval) * q.EaseFactor))

	s.setNextReview(q, today, intervalDays)
	s.setEaseFactorWithPenalty(q, grade)
	q.LastReviewed = today
	q.Familiarity = grade
}

func (s SM2Scheduler) setNextReview(q *Question, date time.Time, intervalDays int) {
	if intervalDays < 1 {
		intervalDays = 1
	} else if intervalDays > s.maxInterval {
		intervalDays = s.maxInterval
	}
	q.NextReview = s.Clock.AddDays(date, intervalDays)
}

// Update ease factor based on familiarity and importance
func (s SM2Scheduler) setEaseFactorWithPenalty(q *Question, grade Familiarity) {
	// How forgiving each importance level is
	importanceEaseBonus := map[Importance]float64{
		LowImportance:      0.15, // More aggressive boost
		MediumImportance:   0.10,
		HighImportance:     0.05,
		CriticalImportance: 0.03, // Tightest boost
	}

	// Penalties based on recall difficulty
	familiarityPenalty := map[Familiarity]float64{
		VeryHard: 0.40,
		Hard:     0.25,
		Medium:   0.10,
		Easy:     0.00,
		VeryEasy: -0.10, // Negative penalty = small bonus
	}

	bonus := importanceEaseBonus[q.Importance]
	penalty := familiarityPenalty[grade]

	// Apply core adjustment
	q.EaseFactor += bonus - penalty

	// Encourage stability if consistently good
	if q.ReviewCount >= 3 && grade >= Medium {
		q.EaseFactor += bonus * 0.5 // Smaller additive bonus
	}

	s.secureEaseFactorBounds(q)
}

// Secure ease factor within bounds
func (s SM2Scheduler) secureEaseFactorBounds(q *Question) {
	if q.EaseFactor < s.minEaseFactor {
		q.EaseFactor = s.minEaseFactor
	} else if q.EaseFactor > s.maxEaseFactor {
		q.EaseFactor = s.maxEaseFactor
	}
}

func (s SM2Scheduler) CalculatePriorityScore(q *Question) float64 {
	today := s.Clock.Today()

	// Constants: Tuned for prioritizing the most critical items.
	const (
		importanceWeight    = 2.0  // Prioritizes designated importance
		overdueWeight       = 1.0  // Prioritizes items past their due date
		familiarityWeight   = 2.5  // Prioritizes historically difficult items
		reviewPenaltyWeight = -0.5 // De-prioritizes questions seen many times (prevents leeching)
		easePenaltyWeight   = -1.0 // De-prioritizes "easier" questions to focus on struggles
	)

	// Compute overdue days (at least 0)
	overdueDays := int(today.Sub(q.NextReview).Hours() / 24)
	if overdueDays < 0 {
		overdueDays = 0
	}

	// Invert Familiarity (VeryEasy = 0, VeryHard = 4)
	// A higher score for harder questions.
	famScore := 4 - int(q.Familiarity)

	score := importanceWeight*float64(q.Importance) +
		overdueWeight*float64(overdueDays) +
		familiarityWeight*float64(famScore) +
		reviewPenaltyWeight*float64(q.ReviewCount) +
		easePenaltyWeight*q.EaseFactor

	return score
}
