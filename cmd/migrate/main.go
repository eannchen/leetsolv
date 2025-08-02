package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"leetsolv/core"
	"leetsolv/storage"
)

// OldQuestion represents the structure from the old questions.json file
type OldQuestion struct {
	ID           int     `json:"id"`
	URL          string  `json:"url"`
	Note         string  `json:"note"`
	Familiarity  int     `json:"familiarity"`
	Importance   int     `json:"importance"`
	LastReviewed string  `json:"last_reviewed"`
	NextReview   string  `json:"next_review"`
	ReviewCount  int     `json:"review_count"`
	EaseFactor   float64 `json:"ease_factor"`
	CreatedAt    string  `json:"created_at"`
}

func main() {
	// Read the old questions.json file
	oldData, err := os.ReadFile("questions.json")
	if err != nil {
		log.Fatalf("Failed to read questions.json: %v", err)
	}

	// Parse the old data
	var oldQuestions []OldQuestion
	if err := json.Unmarshal(oldData, &oldQuestions); err != nil {
		log.Fatalf("Failed to parse old questions.json: %v", err)
	}

	// Create new QuestionStore
	store := &storage.QuestionStore{
		MaxID:     0,
		Questions: make(map[int]*core.Question),
		URLIndex:  make(map[string]int),
	}

	// Parse time layout
	timeLayout := "2006-01-02T15:04:05-07:00"

	// Convert each old question to new format
	for _, oldQ := range oldQuestions {
		// Parse timestamps
		lastReviewed, err := time.Parse(timeLayout, oldQ.LastReviewed)
		if err != nil {
			log.Printf("Warning: Failed to parse last_reviewed for question %d: %v", oldQ.ID, err)
			lastReviewed = time.Now()
		}

		nextReview, err := time.Parse(timeLayout, oldQ.NextReview)
		if err != nil {
			log.Printf("Warning: Failed to parse next_review for question %d: %v", oldQ.ID, err)
			nextReview = time.Now()
		}

		createdAt, err := time.Parse(timeLayout, oldQ.CreatedAt)
		if err != nil {
			log.Printf("Warning: Failed to parse created_at for question %d: %v", oldQ.ID, err)
			createdAt = time.Now()
		}

		// Convert familiarity (old: 0-4, new: VeryHard-VeryEasy)
		var familiarity core.Familiarity
		switch oldQ.Familiarity {
		case 0:
			familiarity = core.VeryHard
		case 1:
			familiarity = core.Hard
		case 2:
			familiarity = core.Medium
		case 3:
			familiarity = core.Easy
		case 4:
			familiarity = core.VeryEasy
		default:
			familiarity = core.Medium
		}

		// Convert importance (preserve original values: 1-3)
		var importance core.Importance
		switch oldQ.Importance {
		case 1:
			importance = core.MediumImportance // 1
		case 2:
			importance = core.HighImportance // 2
		case 3:
			importance = core.CriticalImportance // 3
		default:
			importance = core.MediumImportance
		}

		// Create new question
		newQuestion := &core.Question{
			ID:           oldQ.ID,
			URL:          oldQ.URL,
			Note:         oldQ.Note,
			Familiarity:  familiarity,
			Importance:   importance,
			LastReviewed: lastReviewed,
			NextReview:   nextReview,
			ReviewCount:  oldQ.ReviewCount,
			EaseFactor:   oldQ.EaseFactor,
			CreatedAt:    createdAt,
		}

		// Add to store
		store.Questions[oldQ.ID] = newQuestion
		store.URLIndex[oldQ.URL] = oldQ.ID

		// Update MaxID
		if oldQ.ID > store.MaxID {
			store.MaxID = oldQ.ID
		}
	}

	// Create file storage instance
	fileStorage := storage.NewFileStorage("questions_new.json", "deltas.json")

	// Save the migrated data
	if err := fileStorage.SaveQuestionStore(store); err != nil {
		log.Fatalf("Failed to save migrated data: %v", err)
	}

	fmt.Printf("Migration completed successfully!\n")
	fmt.Printf("Migrated %d questions\n", len(store.Questions))
	fmt.Printf("Max ID: %d\n", store.MaxID)
	fmt.Printf("Output saved to: questions_new.json\n")

	// Print some statistics
	fmt.Printf("\nStatistics:\n")
	familiarityCount := make(map[core.Familiarity]int)
	importanceCount := make(map[core.Importance]int)

	for _, q := range store.Questions {
		familiarityCount[q.Familiarity]++
		importanceCount[q.Importance]++
	}

	fmt.Printf("Familiarity distribution:\n")
	for f, count := range familiarityCount {
		fmt.Printf("  %v: %d\n", f, count)
	}

	fmt.Printf("Importance distribution:\n")
	for i, count := range importanceCount {
		fmt.Printf("  %v: %d\n", i, count)
	}
}
