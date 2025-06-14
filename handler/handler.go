package handler

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"leetsolv/config"
	"leetsolv/core"
	"leetsolv/usecase"
)

type Handler struct {
	QuestionUseCase usecase.QuestionUseCase
}

func NewHandler(questionUseCase usecase.QuestionUseCase) *Handler {
	return &Handler{
		QuestionUseCase: questionUseCase,
	}
}

func (h *Handler) HandleList(scanner *bufio.Scanner) {
	pageSize := config.Env().PageSize
	page := 0

	for {
		questions, totalPages, err := h.QuestionUseCase.PaginatedListQuestions(pageSize, page)
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		if len(questions) == 0 {
			fmt.Println("No questions available.")
			break
		}

		// Display the current page
		fmt.Printf("-- Page %d/%d --\n", page+1, totalPages)
		for _, q := range questions {
			fmt.Printf("[%d] %s (Next: %s)\n", q.ID, q.URL, q.NextReview.Format("2006-01-02")) // Date only
			fmt.Printf("   Note: %s\n", q.Note)
		}

		// Handle user input for pagination
		if page+1 == totalPages {
			fmt.Println("\nEnd of list.")
			break
		}

		fmt.Print("\nPress [Enter] for next page, [q] to quit: ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		if input == "q" {
			break
		}

		page++
	}
}

func (h *Handler) HandleStatus() {
	due, upcoming, total, err := h.QuestionUseCase.ListQuestionsSummary()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Total Questions: %d\n\n", total)

	fmt.Printf("Due Questions: %d\n", len(due))
	for _, q := range due {
		fmt.Printf("[%d] %s\n   Note: %s\n", q.ID, q.URL, q.Note)
	}

	fmt.Printf("\nUpcoming Questions (within 3 days): %d\n", len(upcoming))
	for _, q := range upcoming {
		fmt.Printf("[%d] %s (Next: %s)\n   Note: %s\n", q.ID, q.URL, q.NextReview.Format("2006-01-02"), q.Note)
	}
	fmt.Printf("\n")
}

func (h *Handler) HandleUpsert(scanner *bufio.Scanner) {
	rawURL := readLine(scanner, "URL: ")

	// Normalize and validate the URL
	url, err := h.QuestionUseCase.NormalizeLeetCodeURL(rawURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	note := readLine(scanner, "Note: ")

	fmt.Println("Familiarity:")
	fmt.Println("1. Struggled    - Solved, but barely. Needed heavy effort or help.")
	fmt.Println("2. Clumsy       - Solved with partial understanding, some errors.")
	fmt.Println("3. Decent       - Solved mostly right, but not smooth.")
	fmt.Println("4. Smooth       - Solved confidently and clearly.")
	fmt.Println("5. Fluent       - Solved perfectly and instantly.")
	famInput := readLine(scanner, "\nEnter a number (1-5): ")
	fam, err := strconv.Atoi(famInput)
	if err != nil || fam < 1 || fam > 5 {
		fmt.Println("Invalid familiarity level. Please enter a number between 1 and 5.")
		return
	}

	fmt.Printf("\n")

	fmt.Println("Importance:")
	fmt.Println("1. Low Importance")
	fmt.Println("2. Medium Importance")
	fmt.Println("3. High Importance")
	fmt.Println("4. Critical Importance")
	impInput := readLine(scanner, "\nEnter a number (1-4): ")
	imp, err := strconv.Atoi(impInput)
	if err != nil || imp < 1 || imp > 4 {
		fmt.Println("Invalid importance level. Please enter a number between 1 and 4.")
		return
	}

	// Adjust familiarity and importance to match enums
	familiarity := core.Familiarity(fam - 1)
	importance := core.Importance(imp - 1)

	// Call the updated UpsertQuestion function
	upsertedQuestion, err := h.QuestionUseCase.UpsertQuestion(url, note, familiarity, importance)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		// Display the upserted question
		fmt.Println("Question upserted:")
		fmt.Printf("[%d] %s\n", upsertedQuestion.ID, upsertedQuestion.URL)
		fmt.Printf("   Note: %s\n", upsertedQuestion.Note)
		fmt.Printf("   Familiarity: %d\n", upsertedQuestion.Familiarity+1)
		fmt.Printf("   Importance: %d\n", upsertedQuestion.Importance+1)
		fmt.Printf("   Last Reviewed: %s\n", upsertedQuestion.LastReviewed.Format("2006-01-02"))
		fmt.Printf("   Next Review: %s\n", upsertedQuestion.NextReview.Format("2006-01-02"))
		fmt.Printf("   Review Count: %d\n", upsertedQuestion.ReviewCount)
		fmt.Printf("   Ease Factor: %.2f\n", upsertedQuestion.EaseFactor)
		fmt.Printf("   Created At: %s\n", upsertedQuestion.CreatedAt.Format("2006-01-02"))
	}
	fmt.Printf("\n")
}

func (h *Handler) HandleDelete(scanner *bufio.Scanner) {
	input := readLine(scanner, "Enter ID or URL to delete the question: ")

	// Confirm before deleting
	confirm := strings.ToLower(readLine(scanner, "Do you want to delete the question? [y/N]: "))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("Cancelled.")
		fmt.Printf("\n")
		return
	}

	if err := h.QuestionUseCase.DeleteQuestion(input); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Question deleted.")
	}
	fmt.Printf("\n")
}

func (h *Handler) HandleUndo(scanner *bufio.Scanner) {
	// Confirm before undo
	confirm := strings.ToLower(readLine(scanner, "Do you want to undo the previous action? [y/N]: "))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("Cancelled.")
		fmt.Printf("\n")
		return
	}

	err := h.QuestionUseCase.Undo()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Undo successful.")
	}
}

func readLine(scanner *bufio.Scanner, prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}
