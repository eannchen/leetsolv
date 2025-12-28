// Package core implements the core models for the leetsolv application.
package core

import "time"

const MaxImportance = int(CriticalImportance) + 1
const MaxFamiliarity = int(VeryEasy) + 1

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

type MemoryUse int

const (
	MemoryReasoned MemoryUse = iota
	MemoryPartial
	MemoryFull
)

// Platform represents the source platform of a DSA problem
type Platform string

const (
	PlatformLeetCode   Platform = "leetcode"
	PlatformHackerRank Platform = "hackerrank"
)

func (p Platform) String() string {
	switch p {
	case PlatformLeetCode:
		return "LeetCode"
	case PlatformHackerRank:
		return "HackerRank"
	}
	return string(p)
}

// ParsedURL contains normalized URL info from any supported platform
type ParsedURL struct {
	Platform      Platform
	NormalizedURL string
	ProblemSlug   string
}

type QuestionMap map[int]*Question

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
	UpdatedAt    time.Time   `json:"updated_at"`
	CreatedAt    time.Time   `json:"created_at"`
}

// ActionType defines the type of action performed.
type ActionType string

const (
	ActionAdd    ActionType = "add"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"
)

func (a ActionType) String() string {
	switch a {
	case ActionAdd:
		return "Add"
	case ActionUpdate:
		return "Update"
	case ActionDelete:
		return "Delete"
	}
	return ""
}

func (a ActionType) PastTenseString() string {
	switch a {
	case ActionAdd:
		return "Added"
	case ActionUpdate:
		return "Updated"
	case ActionDelete:
		return "Deleted"
	}
	return ""
}

type Delta struct {
	Action     ActionType `json:"action"`
	QuestionID int        `json:"question_id"`
	OldState   *Question  `json:"old_state"`
	NewState   *Question  `json:"new_state"`
	CreatedAt  time.Time  `json:"created_at"`
}

// SearchFilter defines filtering criteria for question search
type SearchFilter struct {
	Familiarity *Familiarity `json:"familiarity,omitempty"`
	Importance  *Importance  `json:"importance,omitempty"`
	ReviewCount *int         `json:"review_count,omitempty"`
	DueOnly     bool         `json:"due_only,omitempty"`
}
