package handler

import (
	"bufio"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"leetsolv/config"
	"leetsolv/core"
	"leetsolv/internal/errs"
	"leetsolv/usecase"
)

type Handler interface {
	HandleList(scanner *bufio.Scanner)
	HandleGet(scanner *bufio.Scanner, target string)
	HandleStatus()
	HandleUpsert(scanner *bufio.Scanner)
	HandleDelete(scanner *bufio.Scanner, target string)
	HandleUndo(scanner *bufio.Scanner)
}

type HandlerImpl struct {
	QuestionUseCase usecase.QuestionUseCase
	IO              IOHandler
}

func NewHandler(IOHandler IOHandler, questionUseCase usecase.QuestionUseCase) *HandlerImpl {
	return &HandlerImpl{
		QuestionUseCase: questionUseCase,
		IO:              IOHandler,
	}
}

func (h *HandlerImpl) HandleList(scanner *bufio.Scanner) {

	questions, err := h.QuestionUseCase.ListQuestionsOrderByDesc()
	if err != nil {
		h.IO.PrintError(err)
		return
	}
	if len(questions) == 0 {
		h.IO.Println("No questions available.")
		return
	}

	pageSize := config.Env().PageSize
	page := 0

	for {
		paginatedQuestions, totalPages, err := h.QuestionUseCase.PaginateQuestions(questions, pageSize, page)
		if err != nil {
			h.IO.PrintError(err)
			return
		}

		// Display the current page
		h.IO.PrintfColored(ColorCyan, "-- Page %d/%d --\n", page+1, totalPages)
		for _, q := range paginatedQuestions {
			h.IO.Printf("[%d] %s (Next: %s)\n", q.ID, q.URL, q.NextReview.Format("2006-01-02")) // Date only
			if q.Note == "" {
				h.IO.Printf("   Note: (none)\n")
			} else {
				h.IO.Printf("   Note: %s\n", q.Note)
			}
		}

		// Handle user input for pagination
		if page+1 == totalPages {
			h.IO.Println("\nEnd of list.\n")
			break
		}

		h.IO.Println("\n--- Navigation ---")
		h.IO.Println("[Enter] Next Page    [q] Quit")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		if input == "q" {
			break
		}

		page++
	}
}

func (h *HandlerImpl) HandleGet(scanner *bufio.Scanner, target string) {
	if target == "" {
		target = h.IO.ReadLine(scanner, "Enter ID or URL to get the question details: ")
	}

	question, err := h.QuestionUseCase.GetQuestion(target)
	if err != nil {
		h.IO.PrintError(err)
		return
	}

	h.IO.PrintQuestionDetail(question)
}

func (h *HandlerImpl) HandleStatus() {
	due, upcoming, total, err := h.QuestionUseCase.ListQuestionsSummary()
	if err != nil {
		h.IO.PrintError(err)
		return
	}

	h.IO.PrintlnColored(ColorCyan, "========== Question Status ==========")
	h.IO.Printf("Total Questions: %d\n\n", total)

	if len(due) > 0 {
		h.IO.PrintlnColored(ColorCyan, "---------- Due Questions ----------")
		for _, q := range due {
			h.IO.Printf("[%d] %s\n   Note: %s\n", q.ID, q.URL, q.Note)
		}
	}
	h.IO.Printf("\n")

	h.IO.PrintlnColored(ColorCyan, "---------- Upcoming Questions (within a day) ----------")
	for _, q := range upcoming {
		h.IO.Printf("[%d] %s (Next: %s)\n   Note: %s\n", q.ID, q.URL, q.NextReview.Format("2006-01-02"), q.Note)
	}
	h.IO.Printf("\n")
}

func (h *HandlerImpl) HandleUpsert(scanner *bufio.Scanner) {
	rawURL := h.IO.ReadLine(scanner, "URL: ")

	// Normalize and validate the URL
	url, err := h.normalizeLeetCodeURL(rawURL)
	if err != nil {
		h.IO.PrintError(err)
		return
	}

	note := h.IO.ReadLine(scanner, "Note: ")

	h.IO.Println("Familiarity:")
	h.IO.Println("1. Struggled    - Solved, but barely. Needed heavy effort or help.")
	h.IO.Println("2. Clumsy       - Solved with partial understanding, some errors.")
	h.IO.Println("3. Decent       - Solved mostly right, but not smooth.")
	h.IO.Println("4. Smooth       - Solved confidently and clearly.")
	h.IO.Println("5. Fluent       - Solved perfectly and instantly.")
	famInput := h.IO.ReadLine(scanner, "\nEnter a number (1-5): ")
	familiarity, err := h.validateFamiliarity(famInput)
	if err != nil {
		h.IO.PrintError(err)
		return
	}

	h.IO.Printf("\n")

	h.IO.Println("Importance:")
	h.IO.Println("1. Low Importance")
	h.IO.Println("2. Medium Importance")
	h.IO.Println("3. High Importance")
	h.IO.Println("4. Critical Importance")
	impInput := h.IO.ReadLine(scanner, "\nEnter a number (1-4): ")
	importance, err := h.validateImportance(impInput)
	if err != nil {
		h.IO.PrintError(err)
		return
	}

	// Call the updated UpsertQuestion function
	upsertedQuestion, err := h.QuestionUseCase.UpsertQuestion(url, note, familiarity, importance)
	if err != nil {
		h.IO.PrintError(err)
	} else {
		// Display the upserted question
		h.IO.Printf("\n")
		h.IO.PrintlnColored(ColorGreen, "[✔] Question upserted:")
		h.IO.PrintQuestionDetail(upsertedQuestion)
	}
	h.IO.Printf("\n")
}

func (h *HandlerImpl) validateFamiliarity(input string) (core.Familiarity, error) {
	fam, err := strconv.Atoi(input)
	if err != nil || fam < 1 || fam > 5 {
		return 0, errs.ErrInvalidFamiliarityLevel
	}
	return core.Familiarity(fam - 1), nil
}

func (h *HandlerImpl) validateImportance(input string) (core.Importance, error) {
	imp, err := strconv.Atoi(input)
	if err != nil || imp < 1 || imp > 4 {
		return 0, errs.ErrInvalidImportanceLevel
	}
	return core.Importance(imp - 1), nil
}

func (h *HandlerImpl) normalizeLeetCodeURL(inputURL string) (string, error) {
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

	normalizedURL := "https://leetcode.com/problems/" + matches[1] + "/"
	return normalizedURL, nil
}

func (h *HandlerImpl) HandleDelete(scanner *bufio.Scanner, target string) {
	if target == "" {
		target = h.IO.ReadLine(scanner, "Enter ID or URL to delete the question: ")
	}

	// Confirm before deleting
	confirm := strings.ToLower(h.IO.ReadLine(scanner, "Do you want to delete the question? [y/N]: "))
	if confirm != "y" && confirm != "yes" {
		h.IO.Println("Cancelled.")
		h.IO.Printf("\n")
		return
	}

	deletedQuestion, err := h.QuestionUseCase.DeleteQuestion(target)
	if err != nil {
		h.IO.PrintError(err)
		return
	}

	h.IO.Printf("Question Deleted: [%d] %s\n", deletedQuestion.ID, deletedQuestion.URL)
	h.IO.Printf("\n")
}

func (h *HandlerImpl) HandleUndo(scanner *bufio.Scanner) {
	// Confirm before undo
	confirm := strings.ToLower(h.IO.ReadLine(scanner, "Do you want to undo the previous action? [y/N]: "))
	if confirm != "y" && confirm != "yes" {
		h.IO.Println("Cancelled.")
		h.IO.Printf("\n")
		return
	}

	err := h.QuestionUseCase.Undo()
	if err != nil {
		h.IO.PrintError(err)
	} else {
		h.IO.Println("Undo successful.")
	}
}
