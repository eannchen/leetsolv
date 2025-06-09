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
	// Ease factor adjustments for importance levels
	easeAdjustments map[Importance]float64
}

// NewSM2Scheduler creates a new scheduler with configured parameters
func NewSM2Scheduler() *SM2Scheduler {
	return &SM2Scheduler{
		baseIntervals: map[Importance]int{
			LowImportance:      7, // Others
			MediumImportance:   5, // NeetCode 250
			HighImportance:     3, // NeetCode 150
			CriticalImportance: 3, // NeetCode 75
		},
		maxInterval:       60, // Cap at ~2 months to ensure retention
		defaultEaseFactor: 1.3,
		minEaseFactor:     1.3, // Lower bound for ease factor
		maxEaseFactor:     2.5, // Upper bound to prevent overly long intervals
		intervalMultipliers: map[Importance]float64{
			LowImportance:      1.2,  // Slightly longer intervals
			MediumImportance:   1.0,  // Standard
			HighImportance:     0.75, // Tighter intervals
			CriticalImportance: 0.55, // Tightest intervals
		},
		easeAdjustments: map[Importance]float64{
			LowImportance:      0.13, // Grow faster, reviewed less often
			MediumImportance:   0.10, // Standard
			HighImportance:     0.07,
			CriticalImportance: 0.05, // Slow growth to keep reviews tight
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
		s.setEaseFactor(q, 1.8)
	case VeryEasy:
		intervalDays += 5
		s.setEaseFactor(q, 2.0)
	}

	s.setNextReview(q, today, intervalDays)
	return q
}

// Schedule updates the question's review schedule based on familiarity and importance
func (s SM2Scheduler) Schedule(q *Question, grade Familiarity) {
	// Increment review count
	q.ReviewCount++

	today := s.today()

	// Skip scheduling if reviewed today
	if q.LastReviewed == today {
		q.Familiarity = grade
		return
	}

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
	if overdueDays > 3 && q.Importance > LowImportance && q.Familiarity < VeryEasy {
		penaltyFactor := math.Min(0.02+float64(overdueDays-3)*0.01, 0.1) // increase gradually
		q.EaseFactor -= penaltyFactor
		s.secureEaseFactorBounds(q)
	}

	dayElapsed := float64(today.Sub(q.LastReviewed).Hours()) / 24.0
	if dayElapsed < 2 {
		dayElapsed = 2
	}
	intervalDays := int(math.Round(
		dayElapsed * q.EaseFactor * s.intervalMultipliers[q.Importance], // Apply importance-based multiplier
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
}

// Update ease factor based on familiarity and importance
func (s SM2Scheduler) setEaseFactorWithPenalty(q *Question, grade Familiarity) {
	easeBonus := s.easeAdjustments[q.Importance]
	penalty := math.Min(0.15, float64(5-grade)*0.035) // Penalty for lower grades
	q.EaseFactor += easeBonus - penalty
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
