package usecase

import (
	"errors"
	"sort"
	"strconv"
	"time"

	"leetsolv/core"
	"leetsolv/logger"
	"leetsolv/storage"
)

// QuestionUseCase defines the interface for question use cases
type QuestionUseCase interface {
	ListQuestionsSummary() ([]core.Question, []core.Question, int, error)
	PaginatedListQuestions(pageSize, page int) ([]core.Question, int, error)
	GetQuestion(target string) (*core.Question, error)
	UpsertQuestion(url, note string, familiarity core.Familiarity, importance core.Importance) (*core.Question, error)
	DeleteQuestion(target string) (*core.Question, error)
	Undo() error
}

// QuestionUseCaseImpl struct encapsulates dependencies for use cases
type QuestionUseCaseImpl struct {
	Storage   storage.Storage
	Scheduler core.Scheduler
}

// NewQuestionUseCase creates a new QuestionUseCase instance
func NewQuestionUseCase(storage storage.Storage, scheduler core.Scheduler) *QuestionUseCaseImpl {
	return &QuestionUseCaseImpl{
		Storage:   storage,
		Scheduler: scheduler,
	}
}

func (u *QuestionUseCaseImpl) ListQuestionsSummary() ([]core.Question, []core.Question, int, error) {
	questions, err := u.Storage.Load()
	if err != nil {
		return nil, nil, 0, err
	}

	now := time.Now().Truncate(24 * time.Hour)
	threeDaysLater := now.AddDate(0, 0, 3)

	due := []core.Question{}
	upcoming := []core.Question{}

	for _, q := range questions {
		nextReviewDate := q.NextReview.Truncate(24 * time.Hour)
		if !nextReviewDate.After(now) {
			due = append(due, q)
		} else if nextReviewDate.Before(threeDaysLater) {
			upcoming = append(upcoming, q)
		}
	}

	sort.Slice(due, func(i, j int) bool {
		return u.Scheduler.CalculatePriorityScore(&due[i]) > u.Scheduler.CalculatePriorityScore(&due[j])
	})

	sort.Slice(upcoming, func(i, j int) bool {
		if upcoming[i].NextReview.Equal(upcoming[j].NextReview) {
			return u.Scheduler.CalculatePriorityScore(&upcoming[i]) > u.Scheduler.CalculatePriorityScore(&upcoming[j])
		}
		return upcoming[i].NextReview.Before(upcoming[j].NextReview)
	})

	total := len(questions)
	return due, upcoming, total, nil
}

func (u *QuestionUseCaseImpl) PaginatedListQuestions(pageSize, page int) ([]core.Question, int, error) {
	questions, err := u.Storage.Load()
	if err != nil {
		return nil, 0, err
	}

	sort.Slice(questions, func(i, j int) bool {
		return questions[i].ID > questions[j].ID
	})

	totalQuestions := len(questions)
	if totalQuestions == 0 {
		return nil, 0, nil
	}

	totalPages := (totalQuestions + pageSize - 1) / pageSize

	if page < 0 || page >= totalPages {
		return nil, totalPages, errors.New("invalid page number")
	}

	start := page * pageSize
	end := start + pageSize
	if end > totalQuestions {
		end = totalQuestions
	}

	return questions[start:end], totalPages, nil
}

func (u *QuestionUseCaseImpl) GetQuestion(target string) (*core.Question, error) {
	logger.Logger().Info.Println("Fetching a random question")

	questions, err := u.Storage.Load()
	if err != nil {
		return nil, err
	}

	if len(questions) == 0 {
		return nil, errors.New("no questions available")
	}

	// is target an ID or URL?
	id, err := strconv.Atoi(target)
	isID := err == nil

	// Find the question by ID or URL
	var foundQuestion *core.Question
	if isID {
		foundQuestion = u.findQuestionByID(questions, id)
	} else {
		for _, q := range questions {
			if q.URL == target {
				foundQuestion = &q
				break
			}
		}
	}

	if foundQuestion == nil {
		return nil, errors.New("question not found")
	}

	// Truncate NextReview to the start of the day
	foundQuestion.NextReview = foundQuestion.NextReview.Truncate(24 * time.Hour)

	return foundQuestion, nil
}

func (u *QuestionUseCaseImpl) UpsertQuestion(url, note string, familiarity core.Familiarity, importance core.Importance) (*core.Question, error) {
	logger.Logger().Info.Printf("Upserting question: URL=%s, Familiarity=%d, Importance=%d", url, familiarity, importance)

	questions, err := u.Storage.Load()
	if err != nil {
		return nil, err
	}

	var upsertedQuestion *core.Question
	found := false
	for i := range questions {
		if questions[i].URL == url {
			questions[i].Note = note
			questions[i].Familiarity = familiarity
			questions[i].Importance = importance
			u.Scheduler.Schedule(&questions[i], familiarity)
			upsertedQuestion = &questions[i]
			found = true
			break
		}
	}

	if !found {
		newID := 1
		for _, q := range questions {
			if q.ID >= newID {
				newID = q.ID + 1
			}
		}
		q := u.Scheduler.ScheduleNewQuestion(newID, url, note, familiarity, importance)
		questions = append(questions, *q)
		upsertedQuestion = q
	}

	if err := u.Storage.Save(questions); err != nil {
		return nil, err
	}
	return upsertedQuestion, nil
}

func (u *QuestionUseCaseImpl) DeleteQuestion(target string) (*core.Question, error) {
	logger.Logger().Info.Printf("Deleting question: Target=%s", target)

	questions, err := u.Storage.Load()
	if err != nil {
		return nil, err
	}

	var newQuestions []core.Question
	var deletedQuestion *core.Question

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
		return nil, errors.New("no matching question found to delete")
	}
	if err := u.Storage.Save(newQuestions); err != nil {
		return nil, err
	}
	return deletedQuestion, nil
}

func (u *QuestionUseCaseImpl) Undo() error {
	logger.Logger().Info.Printf("Undoing last action")
	return u.Storage.Undo()
}

func (u *QuestionUseCaseImpl) findQuestionByID(questions []core.Question, id int) *core.Question {
	// Binary search
	index := sort.Search(len(questions), func(i int) bool {
		return questions[i].ID >= id
	})
	if index < len(questions) && questions[index].ID == id {
		return &questions[index]
	}

	// Fallback to linear search
	for _, q := range questions {
		if q.ID == id {
			return &q
		}
	}
	return nil
}
