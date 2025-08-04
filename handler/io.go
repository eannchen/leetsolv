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
	ioh.Printf("   Familiarity: %d/%d\n", question.Familiarity+1, core.VeryEasy+1)
	ioh.Printf("   Importance: %d/%d\n", question.Importance+1, core.CriticalImportance+1)
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
		case errs.ValidationErrorKind:
			// Validation errors - show in yellow with user-friendly message
			ioh.PrintlnColored(ColorYellow, "⚠️ "+codedErr.UserMessage())
			return
		case errs.BusinessErrorKind:
			// Business errors - show in yellow with user-friendly message
			ioh.PrintlnColored(ColorYellow, "⚠️ "+codedErr.UserMessage())
			return
		case errs.SystemErrorKind:
			// System errors - show in red with technical details
			ioh.PrintlnColored(ColorRed, "❌ "+codedErr.Error())
			return
		}
	}

	ioh.PrintlnColored(ColorRed, "❌ Error: "+err.Error())
}
