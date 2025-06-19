package usecase

import (
	"errors"
	"sort"
	"strconv"
	"time"

	"leetsolv/config"
	"leetsolv/core"
	"leetsolv/logger"
	"leetsolv/storage"
)

// QuestionUseCase defines the interface for question use cases
type QuestionUseCase interface {
	ListQuestionsSummary() ([]core.Question, []core.Question, int, error)
	ListQuestionsOrderByDesc() ([]core.Question, error)
	PaginateQuestions(questions []core.Question, pageSize, page int) ([]core.Question, int, error)
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
	questions, err := u.Storage.LoadQuestions()
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

func (u *QuestionUseCaseImpl) ListQuestionsOrderByDesc() ([]core.Question, error) {
	questions, err := u.Storage.LoadQuestions()
	if err != nil {
		return nil, err
	}
	sort.Slice(questions, func(i, j int) bool {
		return questions[i].ID > questions[j].ID
	})
	return questions, nil
}

func (u *QuestionUseCaseImpl) PaginateQuestions(questions []core.Question, pageSize, page int) ([]core.Question, int, error) {
	totalQuestions := len(questions)
	if totalQuestions == 0 {
		return nil, 0, nil
	}

	// Round up to get total pages needed; ensures partial last page is counted
	totalPages := (totalQuestions + pageSize - 1) / pageSize

	if page < 0 || page >= totalPages {
		return nil, totalPages, errors.New("invalid page number")
	}

	// 0-index-based page
	start := page * pageSize
	end := start + pageSize
	if end > totalQuestions {
		end = totalQuestions
	}
	return questions[start:end], totalPages, nil
}

func (u *QuestionUseCaseImpl) GetQuestion(target string) (*core.Question, error) {

	questions, err := u.Storage.LoadQuestions()
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
		_, foundQuestion = u.findQuestionByID(questions, id)
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
	return foundQuestion, nil
}

func (u *QuestionUseCaseImpl) UpsertQuestion(url, note string, familiarity core.Familiarity, importance core.Importance) (*core.Question, error) {
	logger.Logger().Info.Printf("Upserting question: URL=%s, Familiarity=%d, Importance=%d", url, familiarity, importance)

	u.Storage.Lock()
	defer u.Storage.Unlock()

	questions, err := u.Storage.LoadQuestions()
	if err != nil {
		return nil, err
	}
	deltas, err := u.Storage.LoadDeltas()
	if err != nil {
		return nil, err
	}

	var newState *core.Question

	var foundQuestionIndex int
	var foundQuestion *core.Question
	for i, q := range questions {
		if q.URL == url {
			foundQuestion = &q
			foundQuestionIndex = i
			break
		}
	}

	if foundQuestion != nil {
		// Update existing question
		newState = &core.Question{
			ID:           foundQuestion.ID,
			URL:          url,
			Note:         note,
			Familiarity:  familiarity,
			Importance:   importance,
			LastReviewed: foundQuestion.LastReviewed,
			NextReview:   foundQuestion.NextReview,
			ReviewCount:  foundQuestion.ReviewCount,
			EaseFactor:   foundQuestion.EaseFactor,
			CreatedAt:    foundQuestion.CreatedAt,
		}
		u.Scheduler.Schedule(newState, familiarity)
		questions[foundQuestionIndex] = *newState

		// Create a delta for the update
		deltas = u.appendDelta(deltas, core.Delta{
			Action:     core.ActionUpdate,
			QuestionID: foundQuestion.ID,
			OldState:   foundQuestion,
			NewState:   newState,
			CreatedAt:  time.Now(),
		})
	} else {
		// Create a new question
		newID := 1
		for _, q := range questions {
			if q.ID >= newID {
				newID = q.ID + 1
			}
		}
		newState = u.Scheduler.ScheduleNewQuestion(newID, url, note, familiarity, importance)
		questions = append(questions, *newState)

		// Create a delta for the new question
		deltas = u.appendDelta(deltas, core.Delta{
			Action:     core.ActionAdd,
			QuestionID: newState.ID,
			OldState:   nil,
			NewState:   newState,
			CreatedAt:  time.Now(),
		})
	}

	if err := u.Storage.SaveQuestions(questions); err != nil {
		return nil, err
	}
	if err := u.Storage.SaveDeltas(deltas); err != nil {
		// TODO: Tell user that the delta was not saved, undo will not work
		logger.Logger().Error.Printf("Failed to save deltas: %v", err)
	}
	return newState, nil
}

func (u *QuestionUseCaseImpl) DeleteQuestion(target string) (*core.Question, error) {
	logger.Logger().Info.Printf("Deleting question: Target=%s", target)

	u.Storage.Lock()
	defer u.Storage.Unlock()

	questions, err := u.Storage.LoadQuestions()
	if err != nil {
		return nil, err
	}

	var deletedQuestion *core.Question
	var deletedIndex int

	id, err := strconv.Atoi(target)
	isID := err == nil
	if isID {
		deletedIndex, deletedQuestion = u.findQuestionByID(questions, id)
	} else {
		for i, q := range questions {
			if q.URL == target {
				deletedIndex = i
				deletedQuestion = &q
				break
			}
		}
	}
	if deletedQuestion == nil {
		return nil, errors.New("no matching question found to delete")
	}

	questions = append(questions[:deletedIndex], questions[deletedIndex+1:]...)
	if err := u.Storage.SaveQuestions(questions); err != nil {
		return nil, err
	}

	// Create a delta for the deletion
	deltas, err := u.Storage.LoadDeltas()
	if err != nil {
		// TODO: Tell user that the delta was not saved, undo will not work
		logger.Logger().Error.Printf("Failed to load deltas for deletion: %v", err)
		return deletedQuestion, nil
	}
	deltas = u.appendDelta(deltas, core.Delta{
		Action:     core.ActionDelete,
		QuestionID: deletedQuestion.ID,
		OldState:   deletedQuestion,
		NewState:   nil,
		CreatedAt:  time.Now(),
	})
	if err := u.Storage.SaveDeltas(deltas); err != nil {
		// TODO: Tell user that the delta was not saved, undo will not work
		logger.Logger().Error.Printf("Failed to save deltas: %v", err)
	}
	return deletedQuestion, nil
}

func (u *QuestionUseCaseImpl) Undo() error {
	logger.Logger().Info.Printf("Undoing last action")

	u.Storage.Lock()
	defer u.Storage.Unlock()

	deltas, err := u.Storage.LoadDeltas()
	if err != nil {
		return err
	}
	if len(deltas) == 0 {
		return errors.New("no actions to undo")
	}

	// Get the last delta
	lastDelta := deltas[len(deltas)-1]

	questions, err := u.Storage.LoadQuestions()
	if err != nil {
		return err
	}

	var deltaError error

	switch lastDelta.Action {
	case core.ActionAdd:
		// Remove the last added question
		if lastDelta.NewState == nil {
			deltaError = errors.New("cannot undo add action with no new state")
		} else {
			questions, deltaError = u.deleteQuestionByID(questions, lastDelta.NewState.ID)
		}
	case core.ActionUpdate:
		// Restore the previous state of the question
		if lastDelta.OldState == nil {
			deltaError = errors.New("cannot undo update action with no old state")
		} else {
			questions, deltaError = u.updateQuestionByID(questions, lastDelta.QuestionID, lastDelta.OldState)
		}
	case core.ActionDelete:
		// restore the deleted question
		if lastDelta.OldState == nil {
			deltaError = errors.New("cannot undo delete action with no old state")
		} else {
			questions = append(questions, *lastDelta.OldState)
		}
	}

	if deltaError != nil {
		logger.Logger().Error.Printf("Undo failed: %v", deltaError)
		return deltaError
	}

	// Save the updated questions
	if err := u.Storage.SaveQuestions(questions); err != nil {
		logger.Logger().Error.Printf("Failed to save questions: %v", err)
		return err
	}

	// Remove the last delta only after successful undo
	deltas = deltas[:len(deltas)-1]
	if err := u.Storage.SaveDeltas(deltas); err != nil {
		logger.Logger().Error.Printf("Failed to save deltas: %v", err)
		return err
	}

	return nil
}

func (u *QuestionUseCaseImpl) appendDelta(deltas []core.Delta, delta core.Delta) []core.Delta {
	deltas = append(deltas, delta)

	maxDelta := config.Env().MaxDelta
	if len(deltas) > maxDelta {
		// Remove the oldest delta if we exceed the maximum limit
		deltas = deltas[len(deltas)-maxDelta:]
	}
	return deltas
}

func (u *QuestionUseCaseImpl) findQuestionByID(questions []core.Question, id int) (index int, question *core.Question) {
	// Binary search
	L, R := 0, len(questions)-1
	for L <= R {
		mid := L + (R-L)/2
		if questions[mid].ID == id {
			return mid, &questions[mid]
		} else if questions[mid].ID < id {
			L = mid + 1
		} else {
			R = mid - 1
		}
	}

	// Fallback to linear search
	for i, q := range questions {
		if q.ID == id {
			return i, &q
		}
	}
	return -1, nil
}

func (u *QuestionUseCaseImpl) updateQuestionByID(questions []core.Question, id int, newState *core.Question) ([]core.Question, error) {
	// Binary search
	L, R := 0, len(questions)-1
	for L <= R {
		mid := L + (R-L)/2
		if questions[mid].ID == id {
			questions[mid] = *newState
			return questions, nil
		} else if questions[mid].ID < id {
			L = mid + 1
		} else {
			R = mid - 1
		}
	}

	// Fallback to linear search
	for i, q := range questions {
		if q.ID == id {
			questions[i] = *newState
			return questions, nil
		}
	}
	return questions, errors.New("question not found")
}

func (u *QuestionUseCaseImpl) deleteQuestionByID(questions []core.Question, id int) ([]core.Question, error) {
	// Binary search
	L, R := 0, len(questions)-1
	for L <= R {
		mid := L + (R-L)/2
		if questions[mid].ID == id {
			return append(questions[:mid], questions[mid+1:]...), nil
		} else if questions[mid].ID < id {
			L = mid + 1
		} else {
			R = mid - 1
		}
	}

	// Fallback to linear search
	for i, q := range questions {
		if q.ID == id {
			return append(questions[:i], questions[i+1:]...), nil
		}
	}
	return questions, errors.New("question not found")
}
