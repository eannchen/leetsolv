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
	defaultEaseFactor float64
	minEaseFactor     float64
	maxEaseFactor     float64
	// Interval multipliers for importance levels
	intervalMultipliers map[Importance]float64
}

// NewSM2Scheduler creates a new scheduler with configured parameters
func NewSM2Scheduler() *SM2Scheduler {
	return &SM2Scheduler{
		baseIntervals: map[Importance]int{
			LowImportance:      10, // Others
			MediumImportance:   7,  // NeetCode 250
			HighImportance:     5,  // NeetCode 150
			CriticalImportance: 4,  // NeetCode 75
		},
		maxInterval:       90, // Cap at ~3 months to ensure retention
		defaultEaseFactor: 1.8,
		minEaseFactor:     1.5, // Lower bound for ease factor
		maxEaseFactor:     2.5, // Upper bound to prevent overly long intervals
		intervalMultipliers: map[Importance]float64{
			LowImportance:      1.2,
			MediumImportance:   1.0, // Standard
			HighImportance:     0.9, // Tighter intervals
			CriticalImportance: 0.9, // Tighter intervals
		},
	}
}

func (s SM2Scheduler) ScheduleNewQuestion(id int, url, note string, grade Familiarity, importance Importance) *Question {

	today := s.today()

	q := &Question{
		ID:           id,
		URL:          url,
		Note:         note,
		Familiarity:  grade,
		Importance:   importance,
		EaseFactor:   s.defaultEaseFactor,
		ReviewCount:  1,
		LastReviewed: today,
		CreatedAt:    today,
	}

	// Get base interval based on importance
	intervalDays := s.baseIntervals[importance]

	switch grade {
	case Easy:
		intervalDays += 3
		s.setEaseFactor(q, 2.0)
	case VeryEasy:
		intervalDays += 5
		s.setEaseFactor(q, 2.2)
	}

	s.setNextReview(q, today, intervalDays)
	return q
}

// Schedule updates the question's review schedule based on familiarity and importance
func (s SM2Scheduler) Schedule(q *Question, grade Familiarity) {
	// Increment review count
	q.ReviewCount++

	today := s.today()

	// Get base interval based on importance
	baseInterval := s.baseIntervals[q.Importance]

	// For low familiarity, reset to base interval and lower ease factor
	if grade < Medium {
		s.setNextReview(q, today, baseInterval)
		s.setEaseFactorWithPenalty(q, grade)
		q.LastReviewed = today
		q.Familiarity = grade
		return
	}

	// Adjust EaseFactor if question is overdue
	overdueDays := int(today.Sub(q.NextReview).Hours() / 24)
	if overdueDays > 3 && q.Importance > LowImportance && grade != VeryEasy {
		penaltyFactor := math.Min(float64(overdueDays-3)*0.01, 0.1) // max penalty: 13 days
		q.EaseFactor -= penaltyFactor
		s.secureEaseFactorBounds(q)
	}

	prevInterval := math.Min(float64(baseInterval), float64(q.NextReview.Sub(q.LastReviewed).Hours())/24.0)
	intervalDays := int(math.Round(
		prevInterval * q.EaseFactor * s.intervalMultipliers[q.Importance], // Apply importance-based multiplier
	))
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

func (s SM2Scheduler) setEaseFactor(q *Question, easeFactor float64) {
	q.EaseFactor = easeFactor
	s.secureEaseFactorBounds(q)
}

// Update ease factor based on familiarity and importance
func (s SM2Scheduler) setEaseFactorWithPenalty(q *Question, grade Familiarity) {
	easeAdjustments := map[Importance]float64{
		LowImportance:      0.5, // Penalty Tolerance: VeryHard
		MediumImportance:   0.3, // Penalty Tolerance: Hard
		HighImportance:     0.2, // Penalty Tolerance: Medium
		CriticalImportance: 0.1, // Penalty Tolerance: Medium
	}
	easeBonus := easeAdjustments[q.Importance]

	var penalty float64
	switch grade {
	case VeryHard:
		penalty = 0.5
	case Hard:
		penalty = 0.3
	case Medium:
		penalty = 0.1
	}

	q.EaseFactor += easeBonus - penalty
	// After 3+ reviews, add extra easeBonus if the question is not difficult
	if q.ReviewCount >= 3 && grade >= Medium {
		q.EaseFactor += easeBonus
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
