package core

import (
	"math"
	"time"
)

type Scheduler interface {
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
		maxInterval:   60,  // Cap at ~2 months to ensure retention
		minEaseFactor: 1.3, // Lower bound for ease factor
		maxEaseFactor: 2.5, // Upper bound to prevent overly long intervals
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

// Schedule updates the question's review schedule based on familiarity and importance
func (s SM2Scheduler) Schedule(q *Question, grade Familiarity) {
	// Increment review count
	q.ReviewCount++

	// Get current date (truncate to day for consistency)
	now := time.Now().Truncate(24 * time.Hour)

	// Skip scheduling if reviewed today
	if !q.LastReviewed.IsZero() && now.Sub(q.LastReviewed) < 24*time.Hour {
		return
	}

	// Set default ease factor when new question is added
	if q.LastReviewed.IsZero() || q.EaseFactor == 0 {
		s.setDefaultEaseFactor(q)
	}

	// Get base interval based on importance
	baseInterval := s.baseIntervals[q.Importance]

	// First review uses base interval, with interval boost for high familiarity
	if q.ReviewCount == 1 {
		intervalDays := baseInterval
		if grade == Easy {
			intervalDays += 3 // Add 3 days for Easy
			s.setEaseFactor(q, 1.8)
		} else if grade == VeryEasy {
			intervalDays += 5 // Add 5 days for VeryEasy
			s.setEaseFactor(q, 2)
		}
		s.setNextReview(q, now, intervalDays)
		q.LastReviewed = now
		q.Familiarity = grade
		return
	}

	// Subsequent reviews

	// For low familiarity, reset to base interval and lower ease factor
	if grade < Medium {
		s.setNextReview(q, now, baseInterval)
		s.setEaseFactorWithPenalty(q, grade)
		q.LastReviewed = now
		q.Familiarity = grade
		return
	}

	// Adjust EaseFactor if question is overdue
	if !q.NextReview.IsZero() {
		overdueDays := int(now.Sub(q.NextReview).Hours() / 24)
		if overdueDays > 3 && q.Importance > LowImportance && q.Familiarity < VeryEasy {
			penaltyFactor := math.Min(0.02+float64(overdueDays-3)*0.01, 0.1) // increase gradually
			q.EaseFactor -= penaltyFactor
			s.secureEaseFactorBounds(q)
		}
	}

	dayElapsed := float64(now.Sub(q.LastReviewed).Hours()) / 24.0
	if q.LastReviewed.IsZero() || dayElapsed < 2 {
		dayElapsed = 2
	}
	intervalDays := int(math.Round(
		dayElapsed * q.EaseFactor * s.intervalMultipliers[q.Importance], // Apply importance-based multiplier
	))
	s.setNextReview(q, now, intervalDays)
	s.setEaseFactorWithPenalty(q, grade)
	q.LastReviewed = now
	q.Familiarity = grade
}

func (s SM2Scheduler) setNextReview(q *Question, now time.Time, intervalDays int) {
	if intervalDays < 1 {
		intervalDays = 1
	} else if intervalDays > s.maxInterval {
		intervalDays = s.maxInterval
	}
	q.NextReview = now.AddDate(0, 0, intervalDays)
}

func (s SM2Scheduler) setDefaultEaseFactor(q *Question) {
	q.EaseFactor = s.minEaseFactor
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
