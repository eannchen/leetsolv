package core

import (
	"math"
	"time"
)

type Scheduler interface {
	ScheduleNewQuestion(id int, url, note string, grade Familiarity, importance Importance) *Question
	Schedule(q *Question, grade Familiarity)
}

// SM2Scheduler implements the spaced repetition scheduling logic
type SM2Scheduler struct {
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
func NewSM2Scheduler() *SM2Scheduler {
	return &SM2Scheduler{
		baseIntervals: map[Importance]int{
			LowImportance:      8, // Faster growth, so start more spaced
			MediumImportance:   6, // Balanced
			HighImportance:     5, // Slightly tighter
			CriticalImportance: 4, // Tightest
		},
		maxInterval:   90,  // Cap at ~3 months to ensure retention
		minEaseFactor: 1.3, // Lower bound for ease factor
		maxEaseFactor: 2.6, // Upper bound to prevent overly long intervals
	}
}

func (s SM2Scheduler) ScheduleNewQuestion(id int, url, note string, grade Familiarity, importance Importance) *Question {
	today := s.today()

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
		CreatedAt:    today,
	}

	intervalDays := s.baseIntervals[importance]

	// Small tweaks to interval for early grading signal
	switch grade {
	case Easy:
		intervalDays += 2
	case VeryEasy:
		intervalDays += 3
	case VeryHard:
		intervalDays -= 1
	}

	s.setNextReview(q, today, intervalDays)
	return q
}

// Schedule updates the question's review schedule based on familiarity and importance
func (s SM2Scheduler) Schedule(q *Question, grade Familiarity) {
	q.ReviewCount++
	today := s.today()

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
	overdueDays := int(today.Sub(q.NextReview).Hours() / 24)
	if overdueDays > 3 && q.Importance > LowImportance && grade < VeryEasy {
		penaltyFactor := math.Min(float64(overdueDays-3)*0.01, 0.1)
		q.EaseFactor -= penaltyFactor
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

// Get current date (truncate to day for consistency)
func (s SM2Scheduler) today() time.Time {
	return time.Now().Truncate(24 * time.Hour)
}

func (s SM2Scheduler) setNextReview(q *Question, now time.Time, intervalDays int) {
	if intervalDays < 1 {
		intervalDays = 1
	} else if intervalDays > s.maxInterval {
		intervalDays = s.maxInterval
	}
	q.NextReview = now.AddDate(0, 0, intervalDays)
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
		VeryEasy: -0.05, // Negative penalty = small bonus
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
