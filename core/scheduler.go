package core

import "time"

type Scheduler interface {
	Schedule(q *Question, grade Familiarity)
}

type SimpleScheduler struct{}

func (s SimpleScheduler) Schedule(q *Question, grade Familiarity) {
	// Increment the review count every time the question is reviewed
	q.ReviewCount++

	// Adjust the scheduling logic based on familiarity
	if grade < Medium {
		// Reset for low familiarity
		q.NextReview = time.Now().Add(24 * time.Hour)
		q.EaseFactor = 2.5
	} else {
		// Calculate the next review interval for higher familiarity
		if q.ReviewCount == 1 {
			q.NextReview = time.Now().Add(24 * time.Hour)
		} else if q.ReviewCount == 2 {
			q.NextReview = time.Now().Add(6 * 24 * time.Hour)
		} else {
			// Use the SuperMemo 2 (SM2) algorithm for scheduling
			interval := float64((q.NextReview.Sub(q.LastReviewed).Hours())/24) * q.EaseFactor
			q.NextReview = time.Now().Add(time.Duration(interval*24) * time.Hour)
		}

		// Update the EaseFactor based on the familiarity grade
		q.EaseFactor += 0.1 - float64(5-grade)*(0.08+float64(5-grade)*0.02)
		if q.EaseFactor < 1.3 {
			q.EaseFactor = 1.3 // Minimum EaseFactor
		}
	}

	// Update the last reviewed time and familiarity
	q.LastReviewed = time.Now()
	q.Familiarity = grade
}
