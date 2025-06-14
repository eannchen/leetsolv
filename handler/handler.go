package handler

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"leetsolv/config"
	"leetsolv/core"
	"leetsolv/usecase"
)

type Handler interface {
	HandleList(scanner *bufio.Scanner)
	HandleStatus()
	HandleUpsert(scanner *bufio.Scanner)
	HandleDelete(scanner *bufio.Scanner)
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
	pageSize := config.Env().PageSize
	page := 0

	for {
		questions, totalPages, err := h.QuestionUseCase.PaginatedListQuestions(pageSize, page)
		if err != nil {
			h.IO.Println("Error:", err)
			break
		}

		if len(questions) == 0 {
			h.IO.Println("No questions available.")
			break
		}

		// Display the current page
		h.IO.Printf("-- Page %d/%d --\n", page+1, totalPages)
		for _, q := range questions {
			h.IO.Printf("[%d] %s (Next: %s)\n", q.ID, q.URL, q.NextReview.Format("2006-01-02")) // Date only
			h.IO.Printf("   Note: %s\n", q.Note)
		}

		// Handle user input for pagination
		if page+1 == totalPages {
			h.IO.Println("\nEnd of list.")
			break
		}

		h.IO.Println("\nPress [Enter] for next page, [q] to quit: ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		if input == "q" {
			break
		}

		page++
	}
}

func (h *HandlerImpl) HandleStatus() {
	due, upcoming, total, err := h.QuestionUseCase.ListQuestionsSummary()
	if err != nil {
		h.IO.Println("Error:", err)
		return
	}

	h.IO.Printf("Total Questions: %d\n\n", total)

	h.IO.Printf("Due Questions: %d\n", len(due))
	for _, q := range due {
		h.IO.Printf("[%d] %s\n   Note: %s\n", q.ID, q.URL, q.Note)
	}

	h.IO.Printf("\nUpcoming Questions (within 3 days): %d\n", len(upcoming))
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
		h.IO.Println("Error:", err)
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
		h.IO.Println("Invalid familiarity level. Please enter a number between 1 and 5.")
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
		h.IO.Println("Invalid importance level. Please enter a number between 1 and 4.")
		return
	}

	// Call the updated UpsertQuestion function
	upsertedQuestion, err := h.QuestionUseCase.UpsertQuestion(url, note, familiarity, importance)
	if err != nil {
		h.IO.Println("Error:", err)
	} else {
		// Display the upserted question
		h.IO.Println("Question upserted:")
		h.IO.Printf("[%d] %s\n", upsertedQuestion.ID, upsertedQuestion.URL)
		h.IO.Printf("   Note: %s\n", upsertedQuestion.Note)
		h.IO.Printf("   Familiarity: %d\n", upsertedQuestion.Familiarity+1)
		h.IO.Printf("   Importance: %d\n", upsertedQuestion.Importance+1)
		h.IO.Printf("   Last Reviewed: %s\n", upsertedQuestion.LastReviewed.Format("2006-01-02"))
		h.IO.Printf("   Next Review: %s\n", upsertedQuestion.NextReview.Format("2006-01-02"))
		h.IO.Printf("   Review Count: %d\n", upsertedQuestion.ReviewCount)
		h.IO.Printf("   Ease Factor: %.2f\n", upsertedQuestion.EaseFactor)
		h.IO.Printf("   Created At: %s\n", upsertedQuestion.CreatedAt.Format("2006-01-02"))
	}
	h.IO.Printf("\n")
}

func (h *HandlerImpl) validateFamiliarity(input string) (core.Familiarity, error) {
	fam, err := strconv.Atoi(input)
	if err != nil || fam < 1 || fam > 5 {
		return 0, fmt.Errorf("invalid familiarity level: %d", fam)
	}
	return core.Familiarity(fam - 1), nil
}

func (h *HandlerImpl) validateImportance(input string) (core.Importance, error) {
	imp, err := strconv.Atoi(input)
	if err != nil || imp < 1 || imp > 4 {
		return 0, fmt.Errorf("invalid importance level: %d", imp)
	}
	return core.Importance(imp - 1), nil
}

func (h *HandlerImpl) normalizeLeetCodeURL(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", errors.New("invalid URL format")
	}

	if parsedURL.Host != "leetcode.com" || !strings.HasPrefix(parsedURL.Path, "/problems/") {
		return "", errors.New("URL must be from leetcode.com/problems/")
	}

	re := regexp.MustCompile(`^/problems/([^/]+)`)
	matches := re.FindStringSubmatch(parsedURL.Path)
	if len(matches) != 2 {
		return "", errors.New("invalid LeetCode problem URL format")
	}

	normalizedURL := "https://leetcode.com/problems/" + matches[1] + "/"
	return normalizedURL, nil
}

func (h *HandlerImpl) HandleDelete(scanner *bufio.Scanner) {
	input := h.IO.ReadLine(scanner, "Enter ID or URL to delete the question: ")

	// Confirm before deleting
	confirm := strings.ToLower(h.IO.ReadLine(scanner, "Do you want to delete the question? [y/N]: "))
	if confirm != "y" && confirm != "yes" {
		h.IO.Println("Cancelled.")
		h.IO.Printf("\n")
		return
	}

	if err := h.QuestionUseCase.DeleteQuestion(input); err != nil {
		h.IO.Println("Error:", err)
	} else {
		h.IO.Println("Question deleted.")
	}
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
		h.IO.Println("Error:", err)
	} else {
		h.IO.Println("Undo successful.")
	}
}
