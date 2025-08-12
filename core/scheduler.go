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

	// Interval settings (in days)
	maxInterval       int
	baseIntervals     map[Importance]int
	memoryMultipliers map[MemoryUse]float64

	// Ease Factor settings
	minEaseFactor          float64
	maxEaseFactor          float64
	startEaseFactors       map[Importance]float64
	importanceEaseBonus    map[Importance]float64
	familiarityEasePenalty map[Familiarity]float64
	memoryEasePenalty      map[MemoryUse]float64

	// Due Priority List settings
	importanceWeight    float64
	overdueWeight       float64
	familiarityWeight   float64
	reviewPenaltyWeight float64
	easePenaltyWeight   float64
}

func NewSM2Scheduler(clock clock.Clock) *SM2Scheduler {
	return &SM2Scheduler{
		Clock: clock,

		// Interval settings (in days)
		maxInterval: 90,
		baseIntervals: map[Importance]int{
			LowImportance:      8,
			MediumImportance:   6,
			HighImportance:     5,
			CriticalImportance: 4,
		},
		memoryMultipliers: map[MemoryUse]float64{
			MemoryReasoned: 1.00, // don't change
			MemoryPartial:  1.10, // give more forgetting time
			MemoryFull:     1.25, // give even more forgetting time
		},

		// Ease Factor settings
		minEaseFactor: 1.3,
		maxEaseFactor: 2.6,
		startEaseFactors: map[Importance]float64{
			LowImportance:      2.0,
			MediumImportance:   1.9,
			HighImportance:     1.8,
			CriticalImportance: 1.7,
		},
		importanceEaseBonus: map[Importance]float64{
			LowImportance:      0.15,
			MediumImportance:   0.10,
			HighImportance:     0.05,
			CriticalImportance: 0.03,
		},
		familiarityEasePenalty: map[Familiarity]float64{
			VeryHard: -0.40,
			Hard:     -0.25,
			Medium:   -0.10,
			Easy:     0.05,
			VeryEasy: 0.15,
		},
		memoryEasePenalty: map[MemoryUse]float64{
			MemoryReasoned: 0.00,
			MemoryPartial:  -0.02,
			MemoryFull:     -0.05,
		},

		// Due Priority List settings
		importanceWeight:    config.Env().ImportanceWeight,    // Prioritizes designated importance
		overdueWeight:       config.Env().OverdueWeight,       // Prioritizes items past their due date
		familiarityWeight:   config.Env().FamiliarityWeight,   // Prioritizes historically difficult items
		reviewPenaltyWeight: config.Env().ReviewPenaltyWeight, // De-prioritizes questions seen many times (prevents leeching)
		easePenaltyWeight:   config.Env().EasePenaltyWeight,   // De-prioritizes "easier" questions to focus on struggles
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
	intervalDays = int(math.Round(float64(intervalDays) * s.memoryMultipliers[memory]))

	s.setNextReview(q, today, intervalDays)
	return q
}

func (s SM2Scheduler) Schedule(q *Question, memory MemoryUse) {
	q.ReviewCount++
	today := s.Clock.Today()

	baseInterval := s.baseIntervals[q.Importance]

	// Reset if still struggling
	if q.Familiarity == VeryHard {
		s.setNextReview(q, today, baseInterval)
		s.setEaseFactor(q, memory)
		q.LastReviewed = today
		return
	}

	// Penalty for being overdue
	if config.Env().OverduePenalty {
		s.setEaseFactorOverduePenalty(q)
	}

	// Growth based on last interval × EaseFactor × MemoryUse
	prevIntervalDays := int(q.NextReview.Sub(q.LastReviewed).Hours() / 24)
	if prevIntervalDays < 1 {
		prevIntervalDays = baseInterval // fallback
	}
	intervalDays := int(math.Round(float64(prevIntervalDays) * q.EaseFactor * s.memoryMultipliers[memory]))

	s.setNextReview(q, today, intervalDays)
	s.setEaseFactor(q, memory)
	q.LastReviewed = today
}

func (s SM2Scheduler) setNextReview(q *Question, date time.Time, intervalDays int) {

	// Randomize interval to avoid over-fitting to a specific date
	if config.Env().RandomizeInterval {
		// Randomize between -1 and 2 days
		intervalDays += rand.IntN(4) - 1
	}

	// Secure bounds
	if intervalDays < 1 {
		intervalDays = 1
	} else if intervalDays > s.maxInterval {
		intervalDays = s.maxInterval
	}

	q.NextReview = s.Clock.AddDays(date, intervalDays)
}

func (s SM2Scheduler) setEaseFactor(q *Question, memory MemoryUse) {
	bonus := s.importanceEaseBonus[q.Importance]
	penalty := s.familiarityEasePenalty[q.Familiarity]
	memoryPenalty := s.memoryEasePenalty[memory]

	// Apply core adjustments
	q.EaseFactor += bonus
	q.EaseFactor += penalty
	q.EaseFactor += memoryPenalty

	// Encourage stability if consistently good
	if q.ReviewCount >= 3 && q.Familiarity >= Medium && memory == MemoryReasoned {
		q.EaseFactor += bonus * 0.5
	}

	// Secure bounds
	if q.EaseFactor < s.minEaseFactor {
		q.EaseFactor = s.minEaseFactor
	} else if q.EaseFactor > s.maxEaseFactor {
		q.EaseFactor = s.maxEaseFactor
	}
}

func (s SM2Scheduler) setEaseFactorOverduePenalty(q *Question) {
	today := s.Clock.Today()

	overdueLimit := config.Env().OverdueLimit
	overdueDays := int(today.Sub(q.NextReview).Hours() / 24)
	if overdueDays > overdueLimit && q.Importance > LowImportance && q.Familiarity < VeryEasy {
		penaltyFactor := math.Min(float64(overdueDays-overdueLimit)*0.01, 0.1)
		q.EaseFactor -= penaltyFactor
	}
}

func (s SM2Scheduler) CalculatePriorityScore(q *Question) float64 {
	today := s.Clock.Today()

	// Compute overdue days (at least 0)
	overdueDays := int(today.Sub(q.NextReview).Hours() / 24)
	if overdueDays < 0 {
		overdueDays = 0
	}

	// Invert Familiarity (VeryEasy = 0, VeryHard = 4)
	famScore := 4 - int(q.Familiarity)

	score := s.importanceWeight*float64(q.Importance) +
		s.overdueWeight*float64(overdueDays) +
		s.familiarityWeight*float64(famScore) +
		s.reviewPenaltyWeight*float64(q.ReviewCount) +
		s.easePenaltyWeight*q.EaseFactor

	return score
}
