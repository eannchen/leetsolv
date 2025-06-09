package core

import "time"

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
			HighImportance:     4, // NeetCode 150
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
			MediumImportance:   0.10,
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

	// Adjust EaseFactor if question is overdue
	overdueDays := int(now.Sub(q.NextReview).Hours() / 24)
	if overdueDays > 3 && q.Importance > LowImportance {
		penaltyFactor := 0.02 + float64(overdueDays-3)*0.01 // increase gradually
		q.EaseFactor -= penaltyFactor

		// Clamp ease factor within bounds
		s.clampEaseFactor(q)
	}

	// Get base interval based on importance
	baseInterval := s.baseIntervals[q.Importance]

	// Adjust scheduling based on familiarity
	if grade < Medium {
		// For low familiarity, reset to base interval and lower ease factor
		q.NextReview = now.AddDate(0, 0, baseInterval)
		q.EaseFactor = s.minEaseFactor
	} else {
		// Calculate next review interval using SM2 algorithm
		var intervalDays int
		if q.ReviewCount == 1 {
			// First review uses base interval, with familiarity boost for solved problems
			intervalDays = baseInterval
			if grade >= Easy {
				intervalDays += 3 // Add 3 days for Easy/VeryEasy to reflect prior solving
			}
		} else {
			// For subsequent reviews, calculate based on last review date and ease factor
			previousInterval := float64(now.Sub(q.LastReviewed).Hours()) / 24.0
			intervalDays = int(previousInterval * q.EaseFactor)
		}

		// Cap the interval to prevent overly long gaps
		if intervalDays > s.maxInterval {
			intervalDays = s.maxInterval
		}

		// Apply importance-based multiplier
		adjustedInterval := int(float64(intervalDays) * s.intervalMultipliers[q.Importance])

		// Ensure minimum interval of 1 day
		if adjustedInterval < 1 {
			adjustedInterval = 1
		}

		// Set next review date
		q.NextReview = now.AddDate(0, 0, adjustedInterval)

		// Update ease factor based on familiarity and importance
		easeBonus := s.easeAdjustments[q.Importance]
		penalty := float64(5-grade) * (0.04 + float64(5-grade)*0.008) // Penalty for lower grades
		q.EaseFactor += easeBonus - penalty

		// After 3+ reviews, ease keeps increasing even though recall is stable
		if q.ReviewCount >= 3 {
			q.EaseFactor += easeBonus * 0.5
		}

		// Clamp ease factor within bounds
		s.clampEaseFactor(q)
	}

	// Update last reviewed date and familiarity
	q.LastReviewed = now
	q.Familiarity = grade
}

// Clamp ease factor within bounds
func (s SM2Scheduler) clampEaseFactor(q *Question) {
	if q.EaseFactor < s.minEaseFactor {
		q.EaseFactor = s.minEaseFactor
	} else if q.EaseFactor > s.maxEaseFactor {
		q.EaseFactor = s.maxEaseFactor
	}
}
