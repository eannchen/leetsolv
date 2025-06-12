package usecase

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"leetsolv/core"
	"leetsolv/storage"
)

func ListQuestionsSummary(storage storage.Storage) ([]core.Question, []core.Question, int, error) {
	questions, err := storage.Load()
	if err != nil {
		return nil, nil, 0, err
	}

	now := time.Now().Truncate(24 * time.Hour) // Use only the date
	twoWeeksLater := now.AddDate(0, 0, 14)     // Add 14 days to the current date

	due := []core.Question{}
	upcoming := []core.Question{}

	for _, q := range questions {
		nextReviewDate := q.NextReview.Truncate(24 * time.Hour) // Truncate time
		if !nextReviewDate.After(now) {
			due = append(due, q)
		} else if nextReviewDate.Before(twoWeeksLater) {
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

func PaginatedListQuestions(storage storage.Storage, pageSize, page int) ([]core.Question, int, error) {
	questions, err := storage.Load()
	if err != nil {
		return nil, 0, err
	}

	// Sort questions by ID in descending order
	sort.Slice(questions, func(i, j int) bool {
		return questions[i].ID > questions[j].ID
	})

	totalQuestions := len(questions)
	if totalQuestions == 0 {
		return nil, 0, nil
	}

	// Calculate total pages
	totalPages := (totalQuestions + pageSize - 1) / pageSize

	// Ensure the requested page is within bounds
	if page < 0 || page >= totalPages {
		return nil, totalPages, fmt.Errorf("invalid page number")
	}

	// Get the questions for the current page
	start := page * pageSize
	end := start + pageSize
	if end > totalQuestions {
		end = totalQuestions
	}

	// Truncate NextReview to date only for display purposes
	for i := range questions[start:end] {
		questions[start:end][i].NextReview = questions[start:end][i].NextReview.Truncate(24 * time.Hour)
	}

	return questions[start:end], totalPages, nil
}

func UpsertQuestion(storage storage.Storage, scheduler core.Scheduler, url, note string, familiarity core.Familiarity, importance core.Importance) (*core.Question, error) {
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
			questions[i].Importance = importance
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
		q := scheduler.ScheduleNewQuestion(newID, url, note, familiarity, importance)
		questions = append(questions, *q)
		upsertedQuestion = q
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

	for _, q := range questions {
		if (isID && q.ID == id) || (!isID && q.URL == target) {
			deletedQuestion = &q
		} else {
			newQuestions = append(newQuestions, q)
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

// NormalizeLeetCodeURL validates and normalizes a LeetCode problem URL.
func NormalizeLeetCodeURL(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", errors.New("invalid URL format")
	}

	// Ensure the URL is from leetcode.com
	if parsedURL.Host != "leetcode.com" || !strings.HasPrefix(parsedURL.Path, "/problems/") {
		return "", errors.New("URL must be from leetcode.com/problems/")
	}

	// Extract the problem name from the path
	re := regexp.MustCompile(`^/problems/([^/]+)`)
	matches := re.FindStringSubmatch(parsedURL.Path)
	if len(matches) != 2 {
		return "", errors.New("invalid LeetCode problem URL format")
	}

	// Normalize the URL to "https://leetcode.com/problems/{question-name}/"
	normalizedURL := "https://leetcode.com/problems/" + matches[1] + "/"
	return normalizedURL, nil
}

func Undo(storage storage.Storage) error {
	return storage.Undo()
}
