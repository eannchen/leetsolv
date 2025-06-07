package core

import "time"

type Familiarity int

const (
	VeryHard Familiarity = iota
	Hard
	Medium
	Easy
	VeryEasy
)

type Question struct {
	ID           int         `json:"id"`
	URL          string      `json:"url"`
	Note         string      `json:"note"`
	Familiarity  Familiarity `json:"familiarity"`
	LastReviewed time.Time   `json:"last_reviewed"`
	NextReview   time.Time   `json:"next_review"`
	ReviewCount  int         `json:"review_count"`
	EaseFactor   float64     `json:"ease_factor"`
	CreatedAt    time.Time   `json:"created_at"`
}
