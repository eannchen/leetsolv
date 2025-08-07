package usecase

import (
	"errors"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"leetsolv/config"
	"leetsolv/core"
	"leetsolv/internal/clock"
	"leetsolv/internal/errs"
	"leetsolv/internal/logger"
	"leetsolv/internal/rank"
	"leetsolv/internal/tokenizer"
	"leetsolv/storage"
)

// QuestionUseCase defines the interface for question use cases
type QuestionUseCase interface {
	ListQuestionsSummary() (QuestionsSummary, error)
	ListQuestionsOrderByDesc() ([]core.Question, error)
	PaginateQuestions(questions []core.Question, pageSize, page int) ([]core.Question, int, error)
	GetQuestion(target string) (*core.Question, error)
	SearchQuestions(query string) ([]core.Question, error)
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
	dueHeap := rank.NewTopKMinHeap(config.Env().TopKDue)

	var upcomingTotal int
	upcomingHeap := rank.NewTopKMinHeap(config.Env().TopKUpcoming)

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

func (u *QuestionUseCaseImpl) PaginateQuestions(questions []core.Question, pageSize, page int) ([]core.Question, int, error) {
	totalQuestions := len(questions)
	if totalQuestions == 0 {
		return nil, 0, nil
	}

	// Round up to get total pages needed; ensures partial last page is counted
	totalPages := (totalQuestions + pageSize - 1) / pageSize

	if page < 0 || page >= totalPages {
		return nil, totalPages, errs.ErrInvalidPageNumber
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
		return nil, errs.WrapInternalError(err, "Failed to load question store")
	}
	return u.findQuestionByIDOrURL(store, target)
}

func (u *QuestionUseCaseImpl) SearchQuestions(query string) ([]core.Question, error) {
	store, err := u.Storage.LoadQuestionStore()
	if err != nil {
		return nil, errs.WrapInternalError(err, "Failed to load question store")
	}

	idSet1 := store.URLTrie.SearchPrefix(query)
	idSet2 := store.NoteTrie.SearchPrefix(query)

	// Merge the two sets
	if len(idSet1) < len(idSet2) { // Determine the larger set
		idSet1, idSet2 = idSet2, idSet1
	}
	for id := range idSet2 { // Add all IDs from the smaller set to the larger set
		idSet1[id] = struct{}{}
	}

	questions := make([]core.Question, 0, len(idSet1))
	for id := range idSet1 {
		if question, ok := store.Questions[id]; ok {
			questions = append(questions, *question)
		}
	}
	return questions, nil
}

func (u *QuestionUseCaseImpl) UpsertQuestion(url, note string, familiarity core.Familiarity, importance core.Importance) (*core.Question, error) {
	logger.Logger().Info.Printf("Upserting question: URL=%s, Familiarity=%d, Importance=%d", url, familiarity, importance)

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

		// Update the note indices for search
		for _, word := range tokenizer.Tokenize(foundQuestion.Note) {
			store.NoteTrie.Delete(word, foundQuestion.ID)
		}
		for _, word := range tokenizer.Tokenize(newState.Note) {
			store.NoteTrie.Insert(word, newState.ID)
		}

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

		// Create the URL and note indices for search
		questionName, err := u.extractLeetCodeQuestionName(newState.URL)
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
		deltas = u.appendDelta(deltas, core.Delta{
			Action:     core.ActionAdd,
			QuestionID: newState.ID,
			OldState:   nil,
			NewState:   newState,
			CreatedAt:  u.Clock.Now(),
		})
	}

	if err := u.Storage.SaveQuestionStore(store); err != nil {
		return nil, errs.WrapInternalError(err, "Failed to save question store")
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
		// Remove the last added question
		if lastDelta.NewState == nil {
			deltaError = errors.New("cannot undo add action with no new state")
		} else {
			delete(store.Questions, lastDelta.NewState.ID)
			delete(store.URLIndex, lastDelta.NewState.URL)
			for _, word := range tokenizer.Tokenize(lastDelta.NewState.URL) {
				store.URLTrie.Delete(word, lastDelta.NewState.ID)
			}
			for _, word := range tokenizer.Tokenize(lastDelta.NewState.Note) {
				store.NoteTrie.Delete(word, lastDelta.NewState.ID)
			}
		}
	case core.ActionUpdate:
		if lastDelta.OldState == nil {
			deltaError = errors.New("cannot undo update action with no old state")
		} else {
			// Restore the previous state of the question
			store.Questions[lastDelta.QuestionID] = lastDelta.OldState

			// Delete the current state of the question from the trie
			for _, word := range tokenizer.Tokenize(lastDelta.NewState.URL) {
				store.URLTrie.Delete(word, lastDelta.NewState.ID)
			}
			for _, word := range tokenizer.Tokenize(lastDelta.NewState.Note) {
				store.NoteTrie.Delete(word, lastDelta.NewState.ID)
			}
			// Restore the previous state of the question to the trie
			for _, word := range tokenizer.Tokenize(lastDelta.OldState.URL) {
				store.URLTrie.Insert(word, lastDelta.OldState.ID)
			}
			for _, word := range tokenizer.Tokenize(lastDelta.OldState.Note) {
				store.NoteTrie.Insert(word, lastDelta.OldState.ID)
			}
		}
	case core.ActionDelete:
		// Restore the previous state of the question
		if lastDelta.OldState == nil {
			deltaError = errors.New("cannot undo delete action with no old state")
		} else {
			// Restore the previous state of the question
			store.Questions[lastDelta.QuestionID] = lastDelta.OldState

			// Restore the previous state of the question to the trie
			for _, word := range tokenizer.Tokenize(lastDelta.OldState.URL) {
				store.URLTrie.Insert(word, lastDelta.OldState.ID)
			}
			for _, word := range tokenizer.Tokenize(lastDelta.OldState.Note) {
				store.NoteTrie.Insert(word, lastDelta.OldState.ID)
			}
		}
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
		return nil, errs.ErrNoQuestionsAvailable
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
		return nil, errs.ErrQuestionNotFound
	}
	return foundQuestion, nil
}

func (u *QuestionUseCaseImpl) extractLeetCodeQuestionName(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", errs.ErrInvalidURLFormat
	}

	if parsedURL.Host != "leetcode.com" || !strings.HasPrefix(parsedURL.Path, "/problems/") {
		return "", errs.ErrInvalidLeetCodeURL
	}

	re := regexp.MustCompile(`^/problems/([^/]+)`)
	matches := re.FindStringSubmatch(parsedURL.Path)
	if len(matches) != 2 {
		return "", errs.ErrInvalidLeetCodeURLFormat
	}
	questionName := strings.TrimSpace(matches[1])
	return questionName, nil
}
