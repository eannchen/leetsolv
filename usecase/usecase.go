// Package usecase handles the business logic for the leetsolv application.
package usecase

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/eannchen/leetsolv/config"
	"github.com/eannchen/leetsolv/core"
	"github.com/eannchen/leetsolv/internal/clock"
	"github.com/eannchen/leetsolv/internal/errs"
	"github.com/eannchen/leetsolv/internal/logger"
	"github.com/eannchen/leetsolv/internal/rank"
	"github.com/eannchen/leetsolv/internal/tokenizer"
	"github.com/eannchen/leetsolv/internal/urlparser"
	"github.com/eannchen/leetsolv/storage"
)

// QuestionUseCase defines the interface for question use cases
type QuestionUseCase interface {
	ListQuestionsSummary() (QuestionsSummary, error)
	ListQuestionsOrderByDesc() ([]core.Question, error)
	GetQuestion(target string) (*core.Question, error)
	SearchQuestions(queries []string, filter *core.SearchFilter) ([]core.Question, error)
	UpsertQuestion(url, note string, familiarity core.Familiarity, importance core.Importance, memory core.MemoryUse) (*core.Delta, error)
	DeleteQuestion(target string) (*core.Question, error)
	Undo() error
	GetHistory() ([]core.Delta, error)
	GetSettings() error
	UpdateSetting(settingName string, value interface{}) error
	MigrateToUTC() (int, int, error)
}

// QuestionUseCaseImpl struct encapsulates dependencies for use cases
type QuestionUseCaseImpl struct {
	cfg       *config.Config
	logger    *logger.Logger
	Storage   storage.Storage
	Scheduler core.Scheduler
	Clock     clock.Clock
}

// NewQuestionUseCase creates a new QuestionUseCase instance
func NewQuestionUseCase(cfg *config.Config, logger *logger.Logger, storage storage.Storage, scheduler core.Scheduler, clock clock.Clock) *QuestionUseCaseImpl {
	return &QuestionUseCaseImpl{
		cfg:       cfg,
		logger:    logger,
		Storage:   storage,
		Scheduler: scheduler,
		Clock:     clock,
	}
}

type QuestionsSummary struct {
	TopDue        []core.Question // Top-K due questions (by score)
	TotalDue      int             // Total count of due questions
	TopUpcoming   []core.Question // Top-K upcoming questions (by NextReview, then score)
	TotalUpcoming int             // Total count of upcoming (within 1 day)
	Total         int             // Total number of questions in the store
}

func (u *QuestionUseCaseImpl) ListQuestionsSummary() (QuestionsSummary, error) {
	store, err := u.Storage.LoadQuestionStore()
	if err != nil {
		return QuestionsSummary{}, errs.WrapInternalError(err, "Failed to load question store")
	}

	today := u.Clock.Today()
	oneDayLater := u.Clock.AddDays(today, 1)

	var dueTotal int
	dueHeap := rank.NewTopKMinHeap(u.cfg.TopKDue)

	var upcomingTotal int
	upcomingHeap := rank.NewTopKMinHeap(u.cfg.TopKUpcoming)

	for _, q := range store.Questions {
		nextReviewDate := u.Clock.ToDate(q.NextReview)
		if !nextReviewDate.After(today) {
			dueHeap.Push(rank.HeapItem{
				Item:  q,
				Score: u.Scheduler.CalculatePriorityScore(q),
			})
			dueTotal++
		} else if !nextReviewDate.After(oneDayLater) {
			upcomingHeap.Push(rank.HeapItem{
				Item:  q,
				Score: u.Scheduler.CalculatePriorityScore(q),
			})
			upcomingTotal++
		}
	}

	// Pop items in reverse order to get the highest scores first
	due := make([]core.Question, dueHeap.Len())
	for i := len(due) - 1; i > -1; i-- {
		item, _ := dueHeap.Pop()
		due[i] = *(item.Item.(*core.Question))
	}

	// Pop items in reverse order to get the highest scores first
	upcoming := make([]core.Question, upcomingHeap.Len())
	for i := len(upcoming) - 1; i > -1; i-- {
		item, _ := upcomingHeap.Pop()
		upcoming[i] = *(item.Item.(*core.Question))
	}

	total := len(store.Questions)

	return QuestionsSummary{
		TopDue:        due,
		TotalDue:      dueTotal,
		TopUpcoming:   upcoming,
		TotalUpcoming: upcomingTotal,
		Total:         total,
	}, nil
}

func (u *QuestionUseCaseImpl) ListQuestionsOrderByDesc() ([]core.Question, error) {
	store, err := u.Storage.LoadQuestionStore()
	if err != nil {
		return nil, errs.WrapInternalError(err, "Failed to load question store")
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

func (u *QuestionUseCaseImpl) GetQuestion(target string) (*core.Question, error) {
	store, err := u.Storage.LoadQuestionStore()
	if err != nil {
		return nil, errs.WrapInternalError(err, "Failed to load question store")
	}
	return u.findQuestionByIDOrURL(store, target)
}

func (u *QuestionUseCaseImpl) SearchQuestions(queries []string, filter *core.SearchFilter) ([]core.Question, error) {
	store, err := u.Storage.LoadQuestionStore()
	if err != nil {
		return nil, errs.WrapInternalError(err, "Failed to load question store")
	}

	var questions []core.Question

	// If query is provided, search in trie
	if len(queries) > 0 {
		var idSets []map[int]struct{}
		for _, query := range queries {
			idSets = append(idSets, store.URLTrie.SearchPrefix(query), store.NoteTrie.SearchPrefix(query))
		}

		mergedSet := u.mergeIDSets(idSets)

		for id := range mergedSet {
			question, ok := store.Questions[id]
			if !ok {
				continue
			}
			if filter != nil && !u.matchesFilter(*question, *filter) {
				continue
			}
			questions = append(questions, *question)
		}
	} else {
		for _, question := range store.Questions {
			if filter != nil && !u.matchesFilter(*question, *filter) {
				continue
			}
			questions = append(questions, *question)
		}
	}

	return questions, nil
}

func (u *QuestionUseCaseImpl) mergeIDSets(idSets []map[int]struct{}) map[int]struct{} {
	if len(idSets) == 0 {
		return nil
	}

	mergedSet := make(map[int]struct{})
	for _, set := range idSets {
		for element := range set {
			mergedSet[element] = struct{}{}
		}
	}
	return mergedSet
}

// matchesFilter checks if a question matches the given filter criteria
func (u *QuestionUseCaseImpl) matchesFilter(question core.Question, filter core.SearchFilter) bool {
	// Filter by Familiarity
	if filter.Familiarity != nil && question.Familiarity != *filter.Familiarity {
		return false
	}

	// Filter by Importance
	if filter.Importance != nil && question.Importance != *filter.Importance {
		return false
	}

	// Filter by ReviewCount
	if filter.ReviewCount != nil && question.ReviewCount != *filter.ReviewCount {
		return false
	}

	// Filter by due date
	if filter.DueOnly && question.NextReview.After(u.Clock.Now()) {
		return false
	}

	return true
}

func (u *QuestionUseCaseImpl) UpsertQuestion(url, note string, familiarity core.Familiarity, importance core.Importance, memory core.MemoryUse) (*core.Delta, error) {
	u.logger.Info.Printf("Upserting question: URL=%s, Familiarity=%d, Importance=%d", url, familiarity, importance)

	u.Storage.Lock()
	defer u.Storage.Unlock()

	store, err := u.Storage.LoadQuestionStore()
	if err != nil {
		return nil, errs.WrapInternalError(err, "Failed to load question store")
	}

	foundQuestion, err := u.findQuestionByIDOrURL(store, url)
	if err != nil &&
		!errors.Is(err, errs.ErrQuestionNotFound) &&
		!errors.Is(err, errs.ErrNoQuestionsAvailable) {
		return nil, err
	}

	deltas, err := u.Storage.LoadDeltas()
	if err != nil {
		return nil, errs.WrapInternalError(err, "Failed to load deltas")
	}

	var delta *core.Delta
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
			UpdatedAt:    u.Clock.Now(),
			CreatedAt:    foundQuestion.CreatedAt,
		}
		u.Scheduler.Schedule(newState, memory)
		store.Questions[foundQuestion.ID] = newState

		// Update the note indices for search
		for _, word := range tokenizer.Tokenize(foundQuestion.Note) {
			store.NoteTrie.Delete(word, foundQuestion.ID)
		}
		for _, word := range tokenizer.Tokenize(newState.Note) {
			store.NoteTrie.Insert(word, newState.ID)
		}

		// Create a delta for the update
		delta = &core.Delta{
			Action:     core.ActionUpdate,
			QuestionID: foundQuestion.ID,
			OldState:   foundQuestion,
			NewState:   newState,
			CreatedAt:  u.Clock.Now(),
		}
		deltas = u.appendDelta(deltas, *delta)
	} else {
		// Create a new question
		store.MaxID++
		newState = &core.Question{
			ID:          store.MaxID,
			URL:         url,
			Note:        note,
			Familiarity: familiarity,
			Importance:  importance,
			UpdatedAt:   u.Clock.Now(),
			CreatedAt:   u.Clock.Now(),
		}
		newState = u.Scheduler.ScheduleNewQuestion(newState, memory)
		store.Questions[store.MaxID] = newState
		store.URLIndex[url] = store.MaxID

		// Create the URL and note indices for search
		questionName, err := u.extractProblemSlug(newState.URL)
		if err != nil {
			return nil, err
		}
		for _, word := range tokenizer.Tokenize(questionName) {
			store.URLTrie.Insert(word, newState.ID)
		}
		for _, word := range tokenizer.Tokenize(newState.Note) {
			store.NoteTrie.Insert(word, newState.ID)
		}

		// Create a delta for the new question
		delta = &core.Delta{
			Action:     core.ActionAdd,
			QuestionID: newState.ID,
			OldState:   nil,
			NewState:   newState,
			CreatedAt:  u.Clock.Now(),
		}
		deltas = u.appendDelta(deltas, *delta)
	}

	if err := u.Storage.SaveQuestionStore(store); err != nil {
		return nil, errs.WrapInternalError(err, "Failed to save question store")
	}
	if err := u.Storage.SaveDeltas(deltas); err != nil {
		u.logger.Error.Printf("Failed to save deltas: %v", err)
	}
	return delta, nil
}

func (u *QuestionUseCaseImpl) DeleteQuestion(target string) (*core.Question, error) {
	u.logger.Info.Printf("Deleting question: Target=%s", target)

	u.Storage.Lock()
	defer u.Storage.Unlock()

	store, err := u.Storage.LoadQuestionStore()
	if err != nil {
		return nil, errs.WrapInternalError(err, "Failed to load question store")
	}
	deletedQuestion, err := u.findQuestionByIDOrURL(store, target)
	if err != nil {
		return nil, err
	}

	deltas, err := u.Storage.LoadDeltas()
	if err != nil {
		return nil, errs.WrapInternalError(err, "Failed to load deltas")
	}

	// Delete the question from the store
	delete(store.Questions, deletedQuestion.ID)
	delete(store.URLIndex, deletedQuestion.URL)

	// Delete the question from the URL and note indices
	for _, word := range tokenizer.Tokenize(deletedQuestion.URL) {
		store.URLTrie.Delete(word, deletedQuestion.ID)
	}
	for _, word := range tokenizer.Tokenize(deletedQuestion.Note) {
		store.NoteTrie.Delete(word, deletedQuestion.ID)
	}

	if err := u.Storage.SaveQuestionStore(store); err != nil {
		return nil, errs.WrapInternalError(err, "Failed to save question store")
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
		u.logger.Error.Printf("Failed to save deltas: %v", err)
	}
	return deletedQuestion, nil
}

func (u *QuestionUseCaseImpl) Undo() error {
	u.logger.Info.Printf("Undoing last action")

	u.Storage.Lock()
	defer u.Storage.Unlock()

	deltas, err := u.Storage.LoadDeltas()
	if err != nil {
		return errs.WrapInternalError(err, "Failed to load deltas")
	}
	if len(deltas) == 0 {
		return errs.ErrNoActionsToUndo
	}

	store, err := u.Storage.LoadQuestionStore()
	if err != nil {
		return errs.WrapInternalError(err, "Failed to load question store")
	}

	// Get the last delta
	lastDelta := deltas[len(deltas)-1]

	var deltaError error
	switch lastDelta.Action {
	case core.ActionAdd:
		deltaError = u.undoAdd(store, lastDelta)
	case core.ActionUpdate:
		deltaError = u.undoUpdate(store, lastDelta)
	case core.ActionDelete:
		deltaError = u.undoDelete(store, lastDelta)
	}

	if deltaError != nil {
		return errs.WrapInternalError(deltaError, "Failed to undo last action")
	}

	// Save the updated questions
	if err := u.Storage.SaveQuestionStore(store); err != nil {
		return errs.WrapInternalError(err, "Failed to save question store")
	}

	// Remove the last delta only after successful undo
	deltas = deltas[:len(deltas)-1]
	if err := u.Storage.SaveDeltas(deltas); err != nil {
		u.logger.Error.Printf("Failed to save deltas: %v", err)
	}

	return nil
}

func (u *QuestionUseCaseImpl) undoAdd(store *storage.QuestionStore, delta core.Delta) error {
	if delta.NewState == nil {
		return errors.New("cannot undo add action with no new state")
	}

	delete(store.Questions, delta.NewState.ID)
	delete(store.URLIndex, delta.NewState.URL)
	for _, word := range tokenizer.Tokenize(delta.NewState.URL) {
		store.URLTrie.Delete(word, delta.NewState.ID)
	}
	for _, word := range tokenizer.Tokenize(delta.NewState.Note) {
		store.NoteTrie.Delete(word, delta.NewState.ID)
	}
	return nil
}

func (u *QuestionUseCaseImpl) undoUpdate(store *storage.QuestionStore, delta core.Delta) error {
	if delta.OldState == nil && delta.NewState == nil {
		return errors.New("cannot undo update action with no old or new state")
	}

	// Restore the previous state of the question
	store.Questions[delta.QuestionID] = delta.OldState

	// Delete the current state of the question from the trie
	for _, word := range tokenizer.Tokenize(delta.NewState.URL) {
		store.URLTrie.Delete(word, delta.NewState.ID)
	}
	for _, word := range tokenizer.Tokenize(delta.NewState.Note) {
		store.NoteTrie.Delete(word, delta.NewState.ID)
	}
	// Restore the previous state of the question to the trie
	for _, word := range tokenizer.Tokenize(delta.OldState.URL) {
		store.URLTrie.Insert(word, delta.OldState.ID)
	}
	for _, word := range tokenizer.Tokenize(delta.OldState.Note) {
		store.NoteTrie.Insert(word, delta.OldState.ID)
	}
	return nil
}

func (u *QuestionUseCaseImpl) undoDelete(store *storage.QuestionStore, delta core.Delta) error {
	if delta.OldState == nil {
		return errors.New("cannot undo delete action with no old state")
	}

	// Restore the previous state of the question
	store.Questions[delta.QuestionID] = delta.OldState
	store.URLIndex[delta.OldState.URL] = delta.QuestionID

	// Restore the previous state of the question to the trie
	for _, word := range tokenizer.Tokenize(delta.OldState.URL) {
		store.URLTrie.Insert(word, delta.OldState.ID)
	}
	for _, word := range tokenizer.Tokenize(delta.OldState.Note) {
		store.NoteTrie.Insert(word, delta.OldState.ID)
	}
	return nil
}

func (u *QuestionUseCaseImpl) GetHistory() ([]core.Delta, error) {
	deltas, err := u.Storage.LoadDeltas()
	if err != nil {
		return nil, errs.WrapInternalError(err, "Failed to load deltas")
	}

	// Copy the deltas to avoid modifying the cache
	reversedDeltas := make([]core.Delta, len(deltas))

	// Reverse the order to show most recent first
	L, R := 0, len(deltas)-1
	for L <= R {
		reversedDeltas[L], reversedDeltas[R] = deltas[R], deltas[L]
		L++
		R--
	}

	return reversedDeltas, nil
}

func (u *QuestionUseCaseImpl) appendDelta(deltas []core.Delta, delta core.Delta) []core.Delta {
	deltas = append(deltas, delta)

	maxDelta := u.cfg.MaxDelta
	if len(deltas) > maxDelta {
		// Remove the oldest delta if we exceed the maximum limit
		deltas = deltas[len(deltas)-maxDelta:]
	}
	return deltas
}

func (u *QuestionUseCaseImpl) findQuestionByIDOrURL(store *storage.QuestionStore, target string) (*core.Question, error) {
	if len(store.Questions) == 0 {
		return nil, errs.ErrNoQuestionsAvailable
	}

	// is target an ID or URL?
	id, err := strconv.Atoi(target)
	isID := err == nil

	var foundQuestion *core.Question
	if isID {
		foundQuestion = store.Questions[id]
	} else if id, ok := store.URLIndex[target]; ok {
		foundQuestion = store.Questions[id]
	} else {
		for _, q := range store.Questions {
			if strings.EqualFold(q.URL, target) {
				foundQuestion = q
				break
			}
		}
	}

	if foundQuestion == nil {
		return nil, errs.ErrQuestionNotFound
	}
	return foundQuestion, nil
}

// extractProblemSlug extracts the problem slug from any supported platform URL
func (u *QuestionUseCaseImpl) extractProblemSlug(inputURL string) (string, error) {
	parsed, err := urlparser.Parse(inputURL)
	if err != nil {
		return "", err
	}
	return parsed.ProblemSlug, nil
}

func (u *QuestionUseCaseImpl) GetSettings() error {
	// This method is a no-op since we can access config.Env() directly
	// It's kept for interface consistency and potential future use
	return nil
}

func (u *QuestionUseCaseImpl) UpdateSetting(settingName string, value any) error {
	// Use the registry-based approach
	if err := u.cfg.SetSettingValue(settingName, value); err != nil {
		return errs.WrapValidationError(err, fmt.Sprintf("Unknown setting: %s", settingName))
	}

	// Save the configuration
	if err := u.cfg.Save(); err != nil {
		return errs.WrapInternalError(err, "Failed to save settings")
	}

	return nil
}

// MigrateToUTC converts all timestamps in questions and deltas to UTC.
// This is needed for users upgrading from versions that stored local timezone.
// Returns the number of questions and deltas migrated.
func (u *QuestionUseCaseImpl) MigrateToUTC() (int, int, error) {
	u.Storage.Lock()
	defer u.Storage.Unlock()

	// Migrate questions
	store, err := u.Storage.LoadQuestionStore()
	if err != nil {
		return 0, 0, errs.WrapInternalError(err, "Failed to load questions")
	}

	questionsCount := 0
	for _, q := range store.Questions {
		q.LastReviewed = q.LastReviewed.UTC()
		q.NextReview = q.NextReview.UTC()
		q.UpdatedAt = q.UpdatedAt.UTC()
		q.CreatedAt = q.CreatedAt.UTC()
		questionsCount++
	}

	if err := u.Storage.SaveQuestionStore(store); err != nil {
		return 0, 0, errs.WrapInternalError(err, "Failed to save questions")
	}

	// Migrate deltas
	deltas, err := u.Storage.LoadDeltas()
	if err != nil {
		return questionsCount, 0, errs.WrapInternalError(err, "Failed to load deltas")
	}

	deltasCount := 0
	for i := range deltas {
		deltas[i].CreatedAt = deltas[i].CreatedAt.UTC()
		if deltas[i].OldState != nil {
			deltas[i].OldState.LastReviewed = deltas[i].OldState.LastReviewed.UTC()
			deltas[i].OldState.NextReview = deltas[i].OldState.NextReview.UTC()
			deltas[i].OldState.UpdatedAt = deltas[i].OldState.UpdatedAt.UTC()
			deltas[i].OldState.CreatedAt = deltas[i].OldState.CreatedAt.UTC()
		}
		if deltas[i].NewState != nil {
			deltas[i].NewState.LastReviewed = deltas[i].NewState.LastReviewed.UTC()
			deltas[i].NewState.NextReview = deltas[i].NewState.NextReview.UTC()
			deltas[i].NewState.UpdatedAt = deltas[i].NewState.UpdatedAt.UTC()
			deltas[i].NewState.CreatedAt = deltas[i].NewState.CreatedAt.UTC()
		}
		deltasCount++
	}

	if err := u.Storage.SaveDeltas(deltas); err != nil {
		return questionsCount, 0, errs.WrapInternalError(err, "Failed to save deltas")
	}

	return questionsCount, deltasCount, nil
}
