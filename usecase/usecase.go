package usecase

import (
	"errors"
	"sort"
	"strconv"

	"leetsolv/config"
	"leetsolv/core"
	"leetsolv/internal/clock"
	"leetsolv/internal/errs"
	"leetsolv/internal/logger"
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
	Clock     clock.Clock
}

// NewQuestionUseCase creates a new QuestionUseCase instance
func NewQuestionUseCase(storage storage.Storage, scheduler core.Scheduler, clock clock.Clock) *QuestionUseCaseImpl {
	return &QuestionUseCaseImpl{
		Storage:   storage,
		Scheduler: scheduler,
		Clock:     clock,
	}
}

func (u *QuestionUseCaseImpl) ListQuestionsSummary() ([]core.Question, []core.Question, int, error) {
	store, err := u.Storage.LoadQuestionStore()
	if err != nil {
		return nil, nil, 0, errs.WrapInternalError(err, "failed to load question store")
	}

	today := u.Clock.Today()
	oneDayLater := u.Clock.AddDays(today, 1)

	due := []core.Question{}
	upcoming := []core.Question{}

	for _, q := range store.Questions {
		nextReviewDate := u.Clock.ToDate(q.NextReview)
		if !nextReviewDate.After(today) {
			due = append(due, *q)
		} else if !nextReviewDate.After(oneDayLater) {
			upcoming = append(upcoming, *q)
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

	total := len(store.Questions)
	return due, upcoming, total, nil
}

func (u *QuestionUseCaseImpl) ListQuestionsOrderByDesc() ([]core.Question, error) {
	store, err := u.Storage.LoadQuestionStore()
	if err != nil {
		return nil, errs.WrapInternalError(err, "failed to load question store")
	}
	questions := make([]core.Question, 0, len(store.Questions))
	for _, q := range store.Questions {
		questions = append(questions, *q)
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
		return nil, totalPages, errs.Err400InvalidPageNumber
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
	store, err := u.Storage.LoadQuestionStore()
	if err != nil {
		return nil, errs.WrapInternalError(err, "failed to load question store")
	}
	return u.findQuestionByIDOrURL(store, target)
}

func (u *QuestionUseCaseImpl) UpsertQuestion(url, note string, familiarity core.Familiarity, importance core.Importance) (*core.Question, error) {
	logger.Logger().Info.Printf("Upserting question: URL=%s, Familiarity=%d, Importance=%d", url, familiarity, importance)

	u.Storage.Lock()
	defer u.Storage.Unlock()

	store, err := u.Storage.LoadQuestionStore()
	if err != nil {
		return nil, errs.WrapInternalError(err, "failed to load question store")
	}

	foundQuestion, err := u.findQuestionByIDOrURL(store, url)
	if err != nil &&
		!errors.Is(err, errs.Err400QuestionNotFound) &&
		!errors.Is(err, errs.Err400NoQuestionsAvailable) {
		return nil, err
	}

	deltas, err := u.Storage.LoadDeltas()
	if err != nil {
		return nil, errs.WrapInternalError(err, "failed to load deltas")
	}

	var newState *core.Question

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
		store.Questions[foundQuestion.ID] = newState

		// Create a delta for the update
		deltas = u.appendDelta(deltas, core.Delta{
			Action:     core.ActionUpdate,
			QuestionID: foundQuestion.ID,
			OldState:   foundQuestion,
			NewState:   newState,
			CreatedAt:  u.Clock.Now(),
		})
	} else {
		// Create a new question
		store.MaxID++
		newState = u.Scheduler.ScheduleNewQuestion(store.MaxID, url, note, familiarity, importance)
		store.Questions[store.MaxID] = newState
		store.URLIndex[url] = store.MaxID

		// Create a delta for the new question
		deltas = u.appendDelta(deltas, core.Delta{
			Action:     core.ActionAdd,
			QuestionID: newState.ID,
			OldState:   nil,
			NewState:   newState,
			CreatedAt:  u.Clock.Now(),
		})
	}

	if err := u.Storage.SaveQuestionStore(store); err != nil {
		return nil, errs.WrapInternalError(err, "failed to save question store")
	}
	if err := u.Storage.SaveDeltas(deltas); err != nil {
		logger.Logger().Error.Printf("Failed to save deltas: %v", err)
	}
	return newState, nil
}

func (u *QuestionUseCaseImpl) DeleteQuestion(target string) (*core.Question, error) {
	logger.Logger().Info.Printf("Deleting question: Target=%s", target)

	u.Storage.Lock()
	defer u.Storage.Unlock()

	store, err := u.Storage.LoadQuestionStore()
	if err != nil {
		return nil, errs.WrapInternalError(err, "failed to load question store")
	}
	deletedQuestion, err := u.findQuestionByIDOrURL(store, target)
	if err != nil {
		return nil, err
	}

	deltas, err := u.Storage.LoadDeltas()
	if err != nil {
		return nil, errs.WrapInternalError(err, "failed to load deltas")
	}

	delete(store.Questions, deletedQuestion.ID)
	delete(store.URLIndex, deletedQuestion.URL)

	if err := u.Storage.SaveQuestionStore(store); err != nil {
		return nil, errs.WrapInternalError(err, "failed to save question store")
	}

	// Create a delta for the deletion
	deltas = u.appendDelta(deltas, core.Delta{
		Action:     core.ActionDelete,
		QuestionID: deletedQuestion.ID,
		OldState:   deletedQuestion,
		NewState:   nil,
		CreatedAt:  u.Clock.Now(),
	})
	if err := u.Storage.SaveDeltas(deltas); err != nil {
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
		return errs.WrapInternalError(err, "failed to load deltas")
	}
	if len(deltas) == 0 {
		return errs.Err400NoActionsToUndo
	}

	store, err := u.Storage.LoadQuestionStore()
	if err != nil {
		return errs.WrapInternalError(err, "failed to load question store")
	}

	// Get the last delta
	lastDelta := deltas[len(deltas)-1]

	var deltaError error

	switch lastDelta.Action {
	case core.ActionAdd:
		// Remove the last added question
		if lastDelta.NewState == nil {
			deltaError = errors.New("cannot undo add action with no new state")
		} else {
			delete(store.Questions, lastDelta.NewState.ID)
			delete(store.URLIndex, lastDelta.NewState.URL)
		}
	case core.ActionUpdate, core.ActionDelete:
		// Restore the previous state of the question
		if lastDelta.OldState == nil {
			deltaError = errors.New("cannot undo update/delete action with no old state")
		} else {
			store.Questions[lastDelta.QuestionID] = lastDelta.OldState
		}
	}

	if deltaError != nil {
		return errs.WrapInternalError(deltaError, "failed to undo last action")
	}

	// Save the updated questions
	if err := u.Storage.SaveQuestionStore(store); err != nil {
		return errs.WrapInternalError(err, "failed to save question store")
	}

	// Remove the last delta only after successful undo
	deltas = deltas[:len(deltas)-1]
	if err := u.Storage.SaveDeltas(deltas); err != nil {
		logger.Logger().Error.Printf("Failed to save deltas: %v", err)
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

func (u *QuestionUseCaseImpl) findQuestionByIDOrURL(store *storage.QuestionStore, target string) (*core.Question, error) {
	if len(store.Questions) == 0 {
		return nil, errs.Err400NoQuestionsAvailable
	}

	// is target an ID or URL?
	id, err := strconv.Atoi(target)
	isID := err == nil

	var foundQuestion *core.Question
	if isID {
		foundQuestion, _ = store.Questions[id]
	} else if id, ok := store.URLIndex[target]; ok {
		foundQuestion, _ = store.Questions[id]
	} else {
		for _, q := range store.Questions {
			if q.URL == target {
				foundQuestion = q
			}
			break
		}
	}

	if foundQuestion == nil {
		return nil, errs.Err400QuestionNotFound
	}
	return foundQuestion, nil
}
