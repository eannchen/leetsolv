package core

import "time"

type Importance int

const (
	LowImportance Importance = iota
	MediumImportance
	HighImportance
	CriticalImportance
)

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
	Importance   Importance  `json:"importance"`
	LastReviewed time.Time   `json:"last_reviewed"`
	NextReview   time.Time   `json:"next_review"`
	ReviewCount  int         `json:"review_count"`
	EaseFactor   float64     `json:"ease_factor"`
	CreatedAt    time.Time   `json:"created_at"`
}

// ActionType defines the type of action performed.
type ActionType string

const (
	ActionAdd    ActionType = "add"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"
)

type Delta struct {
	Action     ActionType `json:"action"`
	QuestionID int        `json:"question_id"`
	OldState   *Question  `json:"old_state"`
	NewState   *Question  `json:"new_state"`
	CreatedAt  time.Time  `json:"created_at"`
}
