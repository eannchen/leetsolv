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
	HandleUnknown(command string)
	HandleHelp()
	HandleQuit()
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
		h.IO.PrintfColored(ColorHeader, "-- Page %d/%d --\n", page+1, totalPages)
		for _, q := range paginatedQuestions {
			h.IO.PrintQuestionBrief(&q)
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
		if target == "" {
			h.IO.PrintError(errs.ErrInvalidEmptyInput)
			return
		}
	}
	_, err := strconv.Atoi(target)
	isID := err == nil
	if !isID {
		target, err = h.normalizeLeetCodeURL(target)
		if err != nil {
			h.IO.PrintError(err)
			return
		}
	}

	question, err := h.QuestionUseCase.GetQuestion(target)
	if err != nil {
		h.IO.PrintError(err)
		return
	}

	h.IO.PrintQuestionDetail(question)
}

func (h *HandlerImpl) HandleStatus() {
	summary, err := h.QuestionUseCase.ListQuestionsSummary()
	if err != nil {
		h.IO.PrintError(err)
		return
	}

	h.IO.PrintlnColored(ColorHeader, "───────────── Question Status ─────────────")
	h.IO.PrintfColored(ColorStatTotal, "Total Questions: %d\n", summary.Total)
	h.IO.Printf("\n")

	if len(summary.TopDue) > 0 {
		h.IO.PrintlnColored(ColorHeader, "-- Due Questions --")
		if summary.TotalDue > len(summary.TopDue) {
			h.IO.PrintfColored(ColorStatDue, "Total Due: %d (showing top %d)\n", summary.TotalDue, len(summary.TopDue))
		} else {
			h.IO.PrintfColored(ColorStatDue, "Total Due: %d\n", summary.TotalDue)
		}
		for _, q := range summary.TopDue {
			h.IO.PrintQuestionBrief(&q)
		}
	}
	h.IO.Printf("\n")

	h.IO.PrintlnColored(ColorHeader, "-- Upcoming Questions --")
	if summary.TotalUpcoming > len(summary.TopUpcoming) {
		h.IO.PrintfColored(ColorStatUpcoming, "Total Upcoming: %d (showing top %d)\n", summary.TotalUpcoming, len(summary.TopUpcoming))
	} else {
		h.IO.PrintfColored(ColorStatUpcoming, "Total Upcoming: %d\n", summary.TotalUpcoming)
	}
	for _, q := range summary.TopUpcoming {
		h.IO.PrintQuestionBrief(&q)
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
	h.IO.Println("4. Smooth       - Solved smoothly and clearly.")
	h.IO.Println("5. Fluent       - Solved perfectly and confidently.")
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
		h.IO.PrintSuccess("Question Upserted")
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
		h.IO.PrintCancel("Cancelled")
		h.IO.Printf("\n")
		return
	}

	_, err := h.QuestionUseCase.DeleteQuestion(target)
	if err != nil {
		h.IO.PrintError(err)
	} else {
		h.IO.PrintSuccess("Question Deleted")
	}
	h.IO.Printf("\n")
}

func (h *HandlerImpl) HandleUndo(scanner *bufio.Scanner) {
	// Confirm before undo
	confirm := strings.ToLower(h.IO.ReadLine(scanner, "Do you want to undo the previous action? [y/N]: "))
	if confirm != "y" && confirm != "yes" {
		h.IO.PrintCancel("Cancelled")
		h.IO.Printf("\n")
		return
	}

	err := h.QuestionUseCase.Undo()
	if err != nil {
		h.IO.PrintError(err)
	} else {
		h.IO.PrintSuccess("Undo successful")
	}
	h.IO.Printf("\n")
}

func (h *HandlerImpl) HandleUnknown(command string) {
	h.IO.PrintfColored(ColorWarning, "Unknown command: '%s'\n", command)
	h.IO.PrintfColored(ColorWarning, "Available commands: status, list, detail, upsert, remove, undo, help, quit\n")
	h.IO.PrintfColored(ColorWarning, "Type 'help' or 'h' for more information\n")
}

func (h *HandlerImpl) HandleHelp() {
	h.IO.Println("╭───────────────────────────────────────────────────╮")
	h.IO.Println("│                                                   │")
	h.IO.Println("│                                                   │")
	h.IO.Println("│    ░▒▓   LeetSolv — CLI SRS for LeetCode   ▓▒░    │")
	h.IO.Println("│                                                   │")
	h.IO.Println("│                                                   │")
	h.IO.Println("╰───────────────────────────────────────────────────╯")
	h.IO.Println("\nAvailable Commands:")
	h.IO.Println("  status/stat         - Show question status (total, due, upcoming)")
	h.IO.Println("  list/ls             - List all questions with pagination")
	h.IO.Println("  detail/get [id|url] - Get details of a question by ID or URL")
	h.IO.Println("  upsert/add          - Add or update a question")
	h.IO.Println("  remove/rm [id|url]  - Delete a question by ID or URL")
	h.IO.Println("  undo/back           - Undo the last action")
	h.IO.Println("  help/h              - Show this help message")
	h.IO.Println("  quit/q/exit         - Exit the application")
	h.IO.Println("\nTips:")
	h.IO.Println("  • Commands are case-insensitive")
	h.IO.Println("  • Press Enter to continue pagination")
	h.IO.Printf("\n")
}

func (h *HandlerImpl) HandleQuit() {
	h.IO.Println("Goodbye!")
}
