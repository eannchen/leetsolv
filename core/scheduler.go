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
}

// NewSM2Scheduler creates a new scheduler with configured parameters
func NewSM2Scheduler() *SM2Scheduler {
	return &SM2Scheduler{
		baseIntervals: map[Importance]int{
			LowImportance:      5, // Others
			MediumImportance:   4, // NeetCode 250
			HighImportance:     3, // NeetCode 150
			CriticalImportance: 2, // NeetCode 75
		},
		maxInterval:   60,  // Cap at ~2 months to ensure retention
		minEaseFactor: 1.3, // Lower bound for ease factor
		maxEaseFactor: 2.5, // Upper bound to prevent overly long intervals
	}
}

// Schedule updates the question's review schedule based on familiarity and importance
func (s SM2Scheduler) Schedule(q *Question, grade Familiarity) {
	// Increment review count
	q.ReviewCount++

	// Get current date (truncate to day for consistency)
	now := time.Now().Truncate(24 * time.Hour)

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
			// First review uses base interval
			intervalDays = baseInterval
		} else {
			// Subsequent reviews use previous interval * ease factor
			previousInterval := float64(q.NextReview.Sub(q.LastReviewed).Hours()) / 24.0
			intervalDays = int(previousInterval * q.EaseFactor)
		}

		// Cap the interval to prevent overly long gaps
		if intervalDays > s.maxInterval {
			intervalDays = s.maxInterval
		}

		// Apply importance-based multiplier to tighten intervals for critical questions
		intervalMultiplier := map[Importance]float64{
			LowImportance:      1.2, // Slightly longer intervals
			MediumImportance:   1.0, // Standard
			HighImportance:     0.8, // Tighter intervals
			CriticalImportance: 0.6, // Tightest intervals
		}
		adjustedInterval := int(float64(intervalDays) * intervalMultiplier[q.Importance])

		// Ensure minimum interval of 1 day
		if adjustedInterval < 1 {
			adjustedInterval = 1
		}

		// Set next review date
		q.NextReview = now.AddDate(0, 0, adjustedInterval)

		// Update ease factor based on familiarity and importance
		easeAdjustment := map[Importance]float64{
			LowImportance:      0.05, // Slower ease growth for low priority
			MediumImportance:   0.08, // Moderate ease growth
			HighImportance:     0.10, // Faster ease growth
			CriticalImportance: 0.12, // Fastest ease growth
		}
		easeBonus := easeAdjustment[q.Importance]
		penalty := float64(5-grade) * (0.05 + float64(5-grade)*0.01) // Penalty for lower grades
		q.EaseFactor += easeBonus - penalty

		// Clamp ease factor within bounds
		if q.EaseFactor < s.minEaseFactor {
			q.EaseFactor = s.minEaseFactor
		} else if q.EaseFactor > s.maxEaseFactor {
			q.EaseFactor = s.maxEaseFactor
		}
	}

	// Update last reviewed date and familiarity
	q.LastReviewed = now
	q.Familiarity = grade
}
