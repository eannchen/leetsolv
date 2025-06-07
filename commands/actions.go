package commands

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"leetsolv/core"
	"leetsolv/storage"
)

func ListQuestionsSummary(storage storage.Storage) ([]core.Question, []core.Question, int, error) {
	questions, err := storage.Load()
	if err != nil {
		return nil, nil, 0, err
	}

	now := time.Now()
	twoWeeksLater := now.Add(14 * 24 * time.Hour)

	due := []core.Question{}
	upcoming := []core.Question{}

	for _, q := range questions {
		if !q.NextReview.After(now) {
			due = append(due, q)
		} else if q.NextReview.Before(twoWeeksLater) {
			upcoming = append(upcoming, q)
		}
	}

	// Sort upcoming questions by NextReview date
	sort.Slice(upcoming, func(i, j int) bool {
		return upcoming[i].NextReview.Before(upcoming[j].NextReview)
	})

	total := len(questions)
	return due, upcoming, total, nil
}

func UpsertQuestion(storage storage.Storage, scheduler core.Scheduler, url, note string, familiarity core.Familiarity) (*core.Question, error) {
	questions, err := storage.Load()
	if err != nil {
		return nil, err
	}

	var upsertedQuestion *core.Question
	found := false
	for i := range questions {
		if questions[i].URL == url {
			// Update existing question
			questions[i].Note = note
			questions[i].Familiarity = familiarity
			scheduler.Schedule(&questions[i], familiarity)
			upsertedQuestion = &questions[i]
			found = true
			break
		}
	}

	if !found {
		// Generate a new unique ID
		newID := 1
		for _, q := range questions {
			if q.ID >= newID {
				newID = q.ID + 1
			}
		}

		// Add new question
		now := time.Now()
		q := core.Question{
			ID:           newID,
			URL:          url,
			Note:         note,
			Familiarity:  familiarity,
			LastReviewed: now,
			NextReview:   now,
			ReviewCount:  0,
			EaseFactor:   2.5,
			CreatedAt:    now,
		}
		scheduler.Schedule(&q, familiarity)
		questions = append(questions, q)
		upsertedQuestion = &q
	}

	if err := storage.Save(questions); err != nil {
		return nil, err
	}
	return upsertedQuestion, nil
}

func DeleteQuestion(storage storage.Storage, target string) error {
	questions, err := storage.Load()
	if err != nil {
		return err
	}

	var newQuestions []core.Question
	var deletedQuestion *core.Question

	// Check if the target is an ID
	id, err := strconv.Atoi(target)
	isID := err == nil

	if target == "--last" {
		if len(questions) == 0 {
			return fmt.Errorf("no questions to delete")
		}
		last := questions[len(questions)-1]
		newQuestions = questions[:len(questions)-1]
		deletedQuestion = &last
	} else {
		for _, q := range questions {
			if (isID && q.ID == id) || (!isID && q.URL == target) {
				deletedQuestion = &q
			} else {
				newQuestions = append(newQuestions, q)
			}
		}
	}

	if deletedQuestion == nil {
		return fmt.Errorf("no matching question found to delete")
	}
	if err := storage.Save(newQuestions); err != nil {
		return err
	}
	fmt.Printf("Deleted: [%d] %s\n", deletedQuestion.ID, deletedQuestion.URL)
	return nil
}
