package handler

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"leetsolv/core"
	"leetsolv/internal/errs"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	Bold        = "\033[1m"
)

type IOHandler interface {
	Println(a ...interface{})
	Printf(format string, a ...interface{})
	PrintlnColored(color string, a ...interface{})
	PrintfColored(color string, format string, a ...interface{})
	ReadLine(scanner *bufio.Scanner, prompt string) string
	PrintQuestionDetail(question *core.Question)
	PrintError(err error)
}

type IOHandlerImpl struct {
	Reader io.Reader
	Writer io.Writer
}

func NewIOHandler() *IOHandlerImpl {
	return &IOHandlerImpl{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}
}

func (ioh *IOHandlerImpl) Println(a ...interface{}) {
	fmt.Fprintln(ioh.Writer, a...)
}

func (ioh *IOHandlerImpl) Printf(format string, a ...interface{}) {
	fmt.Fprintf(ioh.Writer, format, a...)
}

func (ioh *IOHandlerImpl) PrintlnColored(color string, a ...interface{}) {
	fmt.Fprint(ioh.Writer, color)
	ioh.Println(a...)
	fmt.Fprint(ioh.Writer, ColorReset)
}

func (ioh *IOHandlerImpl) PrintfColored(color string, format string, a ...interface{}) {
	fmt.Fprint(ioh.Writer, color)
	ioh.Printf(format, a...)
	fmt.Fprint(ioh.Writer, ColorReset)
}

func (ioh *IOHandlerImpl) ReadLine(scanner *bufio.Scanner, prompt string) string {
	ioh.Printf("%s", prompt)
	scanner.Scan()
	return scanner.Text()
}

func (ioh *IOHandlerImpl) PrintQuestionDetail(question *core.Question) {
	ioh.Printf("[%d] %s\n", question.ID, question.URL)
	ioh.Printf("   Note: %s\n", question.Note)
	ioh.Printf("   Familiarity: %d\n", question.Familiarity+1)
	ioh.Printf("   Importance: %d\n", question.Importance+1)
	ioh.Printf("   Last Reviewed: %s\n", question.LastReviewed.Format("2006-01-02"))
	ioh.Printf("   Next Review: %s\n", question.NextReview.Format("2006-01-02"))
	ioh.Printf("   Review Count: %d\n", question.ReviewCount)
	ioh.Printf("   Ease Factor: %.2f\n", question.EaseFactor)
	ioh.Printf("   Created At: %s\n", question.CreatedAt.Format("2006-01-02"))
	ioh.Printf("\n")
}

func (ioh *IOHandlerImpl) PrintError(err error) {
	if err == nil {
		return
	}

	// Check if it's a coded error
	var codedErr *errs.CodedError
	if errors.As(err, &codedErr) {
		switch codedErr.Kind {
		case errs.InputErrorKind:
			// Input errors - show in yellow with user-friendly message
			ioh.PrintlnColored(ColorYellow, "⚠️ "+codedErr.Error())
			return
		case errs.SystemErrorKind:
			// System errors - show in red with technical details
			ioh.PrintlnColored(ColorRed, "❌ "+codedErr.Error())
			return
		}
	}

	// Handle non-coded errors (handler validation errors)
	// Check for specific error types using errors.Is
	switch {
	case errors.Is(err, errs.ErrInvalidURLFormat):
		ioh.PrintlnColored(ColorYellow, "⚠️  Please provide a valid URL")
	case errors.Is(err, errs.ErrInvalidLeetCodeURL):
		ioh.PrintlnColored(ColorYellow, "⚠️  Please provide a valid LeetCode problem URL")
	case errors.Is(err, errs.ErrInvalidLeetCodeURLFormat):
		ioh.PrintlnColored(ColorYellow, "⚠️  Please provide a valid LeetCode problem URL format")
	case errors.Is(err, errs.ErrInvalidFamiliarityLevel):
		ioh.PrintlnColored(ColorYellow, "⚠️  Please enter a familiarity level between 1 and 5")
	case errors.Is(err, errs.ErrInvalidImportanceLevel):
		ioh.PrintlnColored(ColorYellow, "⚠️  Please enter an importance level between 1 and 4")
	case errors.Is(err, errs.Err400QuestionNotFound):
		ioh.PrintlnColored(ColorYellow, "⚠️  Question not found. Please check the ID or URL")
	case errors.Is(err, errs.Err400NoQuestionsAvailable):
		ioh.PrintlnColored(ColorYellow, "ℹ️  No questions available yet")
	case errors.Is(err, errs.Err400InvalidPageNumber):
		ioh.PrintlnColored(ColorYellow, "⚠️  Invalid page number")
	case errors.Is(err, errs.Err400NoActionsToUndo):
		ioh.PrintlnColored(ColorYellow, "ℹ️  No actions to undo")
	default:
		// Fallback for unknown errors - show in red
		ioh.PrintlnColored(ColorRed, "❌ Error: "+err.Error())
	}
}
