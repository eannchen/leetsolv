package core

import (
	"math"
	"math/rand/v2"
	"time"

	"leetsolv/config"
	"leetsolv/internal/clock"
)

type Scheduler interface {
	ScheduleNewQuestion(q *Question, memory MemoryUse) *Question
	Schedule(q *Question, memory MemoryUse)
	CalculatePriorityScore(q *Question) float64
}

// SM2Scheduler implements the spaced repetition scheduling logic
type SM2Scheduler struct {
	Clock clock.Clock
	// Base intervals for each importance level (in days)
	baseIntervals map[Importance]int
	// Memory multipliers for each memory use level
	memoryMultipliers map[MemoryUse]float64
	// Maximum interval to prevent overly long gaps (in days)
	maxInterval int
	// Starting ease factors for each importance level
	startEaseFactors map[Importance]float64
	// Minimum and maximum ease factors
	minEaseFactor float64
	maxEaseFactor float64
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
		memoryMultipliers: map[MemoryUse]float64{
			MemoryReasoned: 1.00, // don't change
			MemoryPartial:  1.10, // give more forgetting time
			MemoryFull:     1.25, // give even more forgetting time
		},
		maxInterval: 90, // 90 days is the maximum interval
		startEaseFactors: map[Importance]float64{
			LowImportance:      2.0,
			MediumImportance:   1.9,
			HighImportance:     1.8,
			CriticalImportance: 1.7,
		},
		minEaseFactor: 1.3, // Lower bound for ease factor
		maxEaseFactor: 2.6, // Upper bound to prevent overly long intervals
	}
}

func (s SM2Scheduler) ScheduleNewQuestion(q *Question, memory MemoryUse) *Question {
	today := s.Clock.Today()

	q.EaseFactor = s.startEaseFactors[q.Importance]
	q.ReviewCount = 1
	q.LastReviewed = today

	intervalDays := s.baseIntervals[q.Importance]

	// Small tweaks to interval for early grading signal
	switch q.Familiarity {
	case VeryEasy:
		intervalDays += 7
	case Easy:
		intervalDays += 5
	case Medium:
		intervalDays += 2
	}

	intervalDays *= int(math.Round(float64(intervalDays) * s.memoryMultipliers[memory]))

	s.setNextReview(q, today, intervalDays)
	return q
}

// Schedule updates the question's review schedule based on familiarity and importance
func (s SM2Scheduler) Schedule(q *Question, memory MemoryUse) {
	q.ReviewCount++
	today := s.Clock.Today()

	baseInterval := s.baseIntervals[q.Importance]

	// Reset if still struggling
	if q.Familiarity == VeryHard {
		s.setNextReview(q, today, baseInterval)
		s.setEaseFactorWithPenalty(q)
		q.LastReviewed = today
		return
	}

	// Penalty for being overdue
	if config.Env().OverduePenalty {
		overdueLimit := config.Env().OverdueLimit
		overdueDays := int(today.Sub(q.NextReview).Hours() / 24)
		if overdueDays > overdueLimit && q.Importance > LowImportance && q.Familiarity < VeryEasy {
			penaltyFactor := math.Min(float64(overdueDays-overdueLimit)*0.01, 0.1)
			q.EaseFactor -= penaltyFactor
		}
	}

	// Growth based on last interval × EaseFactor
	prevInterval := int(q.NextReview.Sub(q.LastReviewed).Hours() / 24)
	if prevInterval < 1 {
		prevInterval = baseInterval // fallback
	}

	intervalDays := int(math.Round(float64(prevInterval) * q.EaseFactor * s.memoryMultipliers[memory]))

	s.setNextReview(q, today, intervalDays)
	s.setEaseFactorWithPenalty(q)
	s.setEaseFactorWithMemoryPenalty(q, memory)
	q.LastReviewed = today
}

func (s SM2Scheduler) setNextReview(q *Question, date time.Time, intervalDays int) {

	// Randomize interval to avoid overfitting to the same interval
	if config.Env().RandomizeInterval {
		// Randomize between -1 and 2 days
		intervalDays += rand.IntN(4) - 1
	}

	if intervalDays < 1 {
		intervalDays = 1
	} else if intervalDays > s.maxInterval {
		intervalDays = s.maxInterval
	}
	q.NextReview = s.Clock.AddDays(date, intervalDays)
}

// Update ease factor based on familiarity and importance
func (s SM2Scheduler) setEaseFactorWithPenalty(q *Question) {
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
		Easy:     -0.05,
		VeryEasy: -0.15, // Negative penalty = small bonus
	}

	bonus := importanceEaseBonus[q.Importance]
	penalty := familiarityPenalty[q.Familiarity]

	// Apply core adjustment
	q.EaseFactor += bonus - penalty

	// Encourage stability if consistently good
	if q.ReviewCount >= 3 && q.Familiarity >= Medium {
		q.EaseFactor += bonus * 0.5 // Smaller additive bonus
	}

	s.secureEaseFactorBounds(q)
}

func (s SM2Scheduler) setEaseFactorWithMemoryPenalty(q *Question, memory MemoryUse) {
	memoryPenalty := map[MemoryUse]float64{
		MemoryReasoned: 0.00,  // neutral
		MemoryPartial:  -0.02, // lower EF slightly for partial recall
		MemoryFull:     -0.05, // lower EF for brittle recall
	}

	penalty := memoryPenalty[memory]
	q.EaseFactor += penalty

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

	// Get weights from config
	cfg := config.Env()
	importanceWeight := cfg.ImportanceWeight       // Prioritizes designated importance
	overdueWeight := cfg.OverdueWeight             // Prioritizes items past their due date
	familiarityWeight := cfg.FamiliarityWeight     // Prioritizes historically difficult items
	reviewPenaltyWeight := cfg.ReviewPenaltyWeight // De-prioritizes questions seen many times (prevents leeching)
	easePenaltyWeight := cfg.EasePenaltyWeight     // De-prioritizes "easier" questions to focus on struggles

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
