package handler

import (
	"bufio"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"leetsolv/config"
	"leetsolv/core"
	"leetsolv/internal/errs"
	"leetsolv/internal/tokenizer"
	"leetsolv/usecase"
)

type Handler interface {
	HandleList(scanner *bufio.Scanner)
	HandleSearch(scanner *bufio.Scanner, args []string)
	HandleGet(scanner *bufio.Scanner, target string)
	HandleStatus()
	HandleUpsert(scanner *bufio.Scanner, rawURL string)
	HandleDelete(scanner *bufio.Scanner, target string)
	HandleUndo(scanner *bufio.Scanner)
	HandleHistory()
	HandleUnknown(command string)
	HandleHelp()
	HandleClear()
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
		h.IO.PrintError(errs.ErrNoQuestionsAvailable)
		return
	}

	h.paginateQuestions(scanner, questions)
}

func (h *HandlerImpl) HandleSearch(scanner *bufio.Scanner, args []string) {
	if len(args) == 0 {
		args = strings.Fields(h.IO.ReadLine(scanner, "Enter search query (or press Enter to search all): "))
	}

	targets, filterArgs := h.parseSearchQueries(args)

	filter, err := h.parseFilterArgs(filterArgs)
	if err != nil {
		h.IO.PrintError(err)
		return
	}

	questions, err := h.QuestionUseCase.SearchQuestions(targets, filter)
	if err != nil {
		h.IO.PrintError(err)
		return
	}
	if len(questions) == 0 {
		h.IO.PrintError(errs.ErrNoQuestionsAvailable)
		return
	}

	h.paginateQuestions(scanner, questions)
}

func (h *HandlerImpl) parseSearchQueries(args []string) ([]string, []string) {
	var targets []string
	var filterArgs []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			filterArgs = append(filterArgs, arg)
		} else {
			targets = append(targets, tokenizer.Tokenize(arg)...)
		}
	}

	return targets, filterArgs
}

// parseFilterArgs parses command line arguments for filter criteria
func (h *HandlerImpl) parseFilterArgs(args []string) (*core.SearchFilter, error) {
	filter := &core.SearchFilter{}

	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "--familiarity="):
			val := strings.TrimPrefix(arg, "--familiarity=")
			familiarity, err := h.validateFamiliarity(val)
			if err != nil {
				return nil, err
			}
			filter.Familiarity = &familiarity

		case strings.HasPrefix(arg, "--importance="):
			val := strings.TrimPrefix(arg, "--importance=")
			importance, err := h.validateImportance(val)
			if err != nil {
				return nil, err
			}
			filter.Importance = &importance

		case strings.HasPrefix(arg, "--review-count="):
			val := strings.TrimPrefix(arg, "--review-count=")
			reviewCount, err := strconv.Atoi(val)
			if err != nil {
				return nil, errs.ErrInvalidReviewCount
			}
			filter.ReviewCount = &reviewCount

		case arg == "--due-only":
			filter.DueOnly = true

		default:
			// Skip unknown arguments
			continue
		}
	}

	return filter, nil
}

func (h *HandlerImpl) paginateQuestions(scanner *bufio.Scanner, questions []core.Question) {
	page := 0

	for {
		paginatedQuestions, totalPages, err := h.getQuestionsPage(questions, page)
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

func (h *HandlerImpl) getQuestionsPage(questions []core.Question, page int) ([]core.Question, int, error) {
	totalQuestions := len(questions)
	if totalQuestions == 0 {
		return nil, 0, nil
	}

	pageSize := config.Env().PageSize

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

	h.IO.PrintlnColored(ColorHeader, "-- Due Questions --")
	if summary.TotalDue == 0 {
		h.IO.PrintfColored(ColorStatTotal, "Total Due: 0\n")
	} else if summary.TotalDue > len(summary.TopDue) {
		h.IO.PrintfColored(ColorStatDueTotal, "Total Due: %d  (showing top %d by priority)\n", summary.TotalDue, len(summary.TopDue))
	} else {
		h.IO.PrintfColored(ColorStatDueTotal, "Total Due: %d  (in priority order)\n", summary.TotalDue)
	}
	for _, q := range summary.TopDue {
		h.IO.PrintQuestionBrief(&q)
	}
	h.IO.Printf("\n")

	h.IO.PrintlnColored(ColorHeader, "-- Upcoming Questions --")
	if summary.TotalUpcoming == 0 {
		h.IO.PrintfColored(ColorStatTotal, "Total Upcoming: 0\n")
	} else if summary.TotalUpcoming > len(summary.TopUpcoming) {
		h.IO.PrintfColored(ColorStatTotal, "Total Upcoming: %d  (showing top %d by priority)\n", summary.TotalUpcoming, len(summary.TopUpcoming))
	} else {
		h.IO.PrintfColored(ColorStatTotal, "Total Upcoming: %d  (in priority order)\n", summary.TotalUpcoming)
	}
	for _, q := range summary.TopUpcoming {
		h.IO.PrintQuestionBrief(&q)
	}

	if summary.TotalDue != 0 || summary.TotalUpcoming != 0 {
		h.IO.Printf("\n")
		h.IO.Printf("\n")
		cfg := config.Env()
		h.IO.PrintfColored(ColorAnnotation, "* Priority Scoring Formula = (%.1f×Importance)+(%.1f×Overdue Days)+(%.1f×Difficulty)+(%.1f×Review Count)+(%.1f×Ease Factor)\n",
			cfg.ImportanceWeight, cfg.OverdueWeight, cfg.FamiliarityWeight, cfg.ReviewPenaltyWeight, cfg.EasePenaltyWeight)
	}

	h.IO.Printf("\n")
}

func (h *HandlerImpl) HandleUpsert(scanner *bufio.Scanner, rawURL string) {
	if rawURL == "" {
		h.IO.Println("Provided URL will be normalized to a canonical form to match existing data.")
		rawURL = h.IO.ReadLine(scanner, "URL: ")
	}

	// Normalize and validate the URL
	url, err := h.normalizeLeetCodeURL(rawURL)
	if err != nil {
		h.IO.PrintError(err)
		return
	}
	h.IO.PrintfColored(ColorGreen, "Using normalized URL: %s\n", url)

	h.IO.Printf("\n")

	note := h.IO.ReadLine(scanner, "Note: ")

	h.IO.Printf("\n")

	h.IO.Println("Familiarity:")
	h.IO.Println("1. Struggled - Solved, but barely; needed heavy effort or help.")
	h.IO.Println("2. Clumsy    - Solved with partial understanding or recurring mistakes.")
	h.IO.Println("3. Decent    - Solved mostly right, but with uncertainty or slow spots.")
	h.IO.Println("4. Smooth    - Solved cleanly with clear reasoning, minor pauses, and no real confusion.")
	h.IO.Println("5. Fluent    - Solved confidently with no hesitation.")
	famInput := h.IO.ReadLine(scanner, "\nEnter a number (1-5): ")
	familiarity, err := h.validateFamiliarity(famInput)
	if err != nil {
		h.IO.PrintError(err)
		return
	}

	h.IO.Printf("\n")

	memory := core.MemoryReasoned
	if familiarity >= core.Easy {
		h.IO.Println("Memory Use:")
		h.IO.Println("1. Reasoned - Solved purely from reasoning.")
		h.IO.Println("2. Partial  - Recalled some solution fragments, but still reasoned through the rest.")
		h.IO.Println("3. Full     - Solved mainly from memory of the full approach or exact steps.")
		h.IO.PrintlnColored(ColorAnnotation, "When you report that you solved the problem from memory, the scheduler interprets that as weaker learning.")
		memoryInput := h.IO.ReadLine(scanner, "\nEnter a number (1-3): ")
		memory, err = h.validateMemoryUse(memoryInput)
		if err != nil {
			h.IO.PrintError(err)
			return
		}
		h.IO.Printf("\n")
	}

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
	delta, err := h.QuestionUseCase.UpsertQuestion(url, note, familiarity, importance, memory)
	if err != nil {
		h.IO.PrintError(err)
	} else {
		// Display the upserted question
		h.IO.Printf("\n")
		h.IO.PrintSuccess(fmt.Sprintf("Question %s", delta.Action.PastTenseString()))
		h.IO.PrintQuestionUpsertDetail(delta)
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

func (h *HandlerImpl) validateMemoryUse(input string) (core.MemoryUse, error) {
	memory, err := strconv.Atoi(input)
	if err != nil || memory < 1 || memory > 3 {
		return 0, errs.ErrInvalidMemoryUseLevel
	}
	return core.MemoryUse(memory - 1), nil
}

func (h *HandlerImpl) normalizeLeetCodeURL(inputURL string) (string, error) {
	parsedURL, err := url.Parse(strings.TrimSpace(inputURL))
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

	normalizedURL := "https://leetcode.com/problems/" + questionName + "/"
	return normalizedURL, nil
}

func (h *HandlerImpl) HandleDelete(scanner *bufio.Scanner, target string) {
	if target == "" {
		target = h.IO.ReadLine(scanner, "Enter ID or URL to delete the question: ")
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

	// Confirm before deleting
	confirm := strings.ToLower(h.IO.ReadLine(scanner, "Do you want to delete the question? [y/N]: "))
	if confirm != "y" && confirm != "yes" {
		h.IO.PrintCancel("Cancelled")
		h.IO.Printf("\n")
		return
	}

	_, err = h.QuestionUseCase.DeleteQuestion(target)
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

func (h *HandlerImpl) HandleHistory() {
	deltas, err := h.QuestionUseCase.GetHistory()
	if err != nil {
		h.IO.PrintError(err)
		return
	}

	if len(deltas) == 0 {
		h.IO.Println("No history available.")
		return
	}

	formatWithStrId := "%-6s %-9s %-35s %-22s %s"
	formatWithIntId := "%-6d %-9s %-35s %-22s %s"

	h.IO.PrintlnColored(ColorHeader, "──────────────────────────────────── Action History ────────────────────────────────────")
	h.IO.PrintfColored(ColorHeader, formatWithStrId, "# ID", "Action", "Question", "Changes", "When")
	h.IO.Printf("\n")
	for _, delta := range deltas {
		// Extract question name from URL
		var questionName string
		if delta.NewState != nil {
			questionName = h.extractQuestionNameFromURL(delta.NewState.URL)
		} else if delta.OldState != nil {
			questionName = h.extractQuestionNameFromURL(delta.OldState.URL)
		} else {
			questionName = "unknown"
		}

		// Prepare changes for update actions
		var changeList []string
		if delta.Action == core.ActionUpdate && delta.OldState != nil && delta.NewState != nil {
			changeList = h.getChanges(delta.OldState, delta.NewState)
		}

		// Format the time
		timeDesc := h.formatTimeAgo(delta.CreatedAt)

		// Print entry. If multiple changes, print them on separate aligned lines.
		if len(changeList) == 0 {
			entry := fmt.Sprintf(formatWithIntId, delta.QuestionID, delta.Action.String(), questionName, "", timeDesc)
			h.IO.Println(entry)
		} else {
			// First line with the first change and time
			first := fmt.Sprintf(formatWithIntId, delta.QuestionID, delta.Action.String(), questionName, changeList[0], timeDesc)
			h.IO.Println(first)
			// Continuation lines: only the Change column filled
			for i := 1; i < len(changeList); i++ {
				cont := fmt.Sprintf(formatWithStrId, "", "", "", changeList[i], "")
				h.IO.Println(cont)
			}
		}
	}
	h.IO.Printf("\n")
}

func (h *HandlerImpl) HandleUnknown(command string) {
	h.IO.PrintfColored(ColorWarning, "Unknown command: '%s'\n", command)
	h.IO.PrintfColored(ColorWarning, "Type 'help' or 'h' for more information\n")
}

func (h *HandlerImpl) HandleHelp() {
	h.IO.PrintlnColored(ColorLogo, "╭───────────────────────────────────────────────────╮")
	h.IO.PrintlnColored(ColorLogo, "│                                                   │")
	h.IO.PrintlnColored(ColorLogo, "│                                                   │")
	h.IO.PrintlnColored(ColorLogo, "│    ░▒▓   LeetSolv — CLI SRS for LeetCode   ▓▒░    │")
	h.IO.PrintlnColored(ColorLogo, "│                                                   │")
	h.IO.PrintlnColored(ColorLogo, "│                                                   │")
	h.IO.PrintlnColored(ColorLogo, "╰───────────────────────────────────────────────────╯")
	h.IO.PrintfColored(ColorHeader, "\nAvailable Commands:\n")
	h.IO.Println("  status/stat                   - Show question status (total, due, upcoming)")
	h.IO.Println("  list/ls                       - List all questions with pagination")
	h.IO.Println("  search/s [queries] [filters]  - Search questions on URL or note with optional filters")
	h.IO.Println("                                   Filters: --familiarity=1-5, --importance=1-4, --review-count=N, --due-only")
	h.IO.Println("  detail/get [id|url]           - Get details of a question by ID or URL")
	h.IO.Println("  upsert/add                    - Add or update a question")
	h.IO.Println("  remove/rm/delete/del [id|url] - Delete a question by ID or URL")
	h.IO.Println("  undo/back                     - Undo the last action")
	h.IO.Println("  history/hist/log              - Show action history")
	h.IO.Println("  help/h                        - Show this help message")
	h.IO.Println("  clear/cls                     - Clear the screen")
	h.IO.Println("  quit/q/exit                   - Exit the application")
	h.IO.PrintfColored(ColorHeader, "\nTips:\n")
	h.IO.Println("  • Commands are case-insensitive")
	h.IO.Println("  • Press Enter to continue pagination")
	h.IO.Printf("\n")
}

func (h *HandlerImpl) HandleClear() {
	h.IO.Println("\033[H\033[2J") // Clear screen
	h.HandleHelp()
}

func (h *HandlerImpl) HandleQuit() {
	h.IO.Println("Goodbye!")
}

// Helper methods for history formatting
func (h *HandlerImpl) extractQuestionNameFromURL(url string) string {
	// Extract the question name from LeetCode URL
	// e.g., "https://leetcode.com/problems/two-sum/" -> "two-sum"
	if strings.Contains(url, "/problems/") {
		parts := strings.Split(url, "/problems/")
		if len(parts) > 1 {
			questionPart := parts[1]
			// Remove trailing slash if present
			questionPart = strings.TrimSuffix(questionPart, "/")
			return questionPart
		}
	}
	return "unknown"
}

func (h *HandlerImpl) getChanges(oldState, newState *core.Question) []string {
	var changes []string

	if oldState.Importance != newState.Importance {
		changes = append(changes, fmt.Sprintf("Importance: %d → %d", oldState.Importance+1, newState.Importance+1))
	}

	if oldState.Familiarity != newState.Familiarity {
		changes = append(changes, fmt.Sprintf("Familiarity: %d → %d", oldState.Familiarity+1, newState.Familiarity+1))
	}

	return changes
}

func (h *HandlerImpl) formatTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}
