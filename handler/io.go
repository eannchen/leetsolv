package handler

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/eannchen/leetsolv/core"
	"github.com/eannchen/leetsolv/internal/clock"
	"github.com/eannchen/leetsolv/internal/errs"
)

const (
	ColorSuccess    = ColorGreen
	ColorCancel     = ColorGray
	ColorWarning    = ColorYellow
	ColorError      = ColorRed
	ColorAnnotation = ColorLightGray + Italic

	ColorHeader       = ColorBlue
	ColorQuestionURL  = ColorBlue
	ColorLogo         = ColorOrange
	ColorStatTotal    = ColorBlue
	ColorStatDueTotal = ColorYellow
)

const (
	ColorReset     = "\033[0m"
	ColorRed       = "\033[31m"
	ColorGreen     = "\033[32m"
	ColorYellow    = "\033[33m"
	ColorOrange    = "\033[38;5;208m"
	ColorGray      = "\033[38;5;245m"
	ColorLightGray = "\033[38;5;245m"
	ColorBlue      = "\033[34m"
	Bold           = "\033[1m"
	Italic         = "\033[3m"
)

type IOHandler interface {
	Println(a ...interface{})
	Printf(format string, a ...interface{})
	PrintlnColored(color string, a ...interface{})
	PrintfColored(color string, format string, a ...interface{})
	ReadLine(scanner *bufio.Scanner, prompt string) string
	PrintQuestionBrief(q *core.Question)
	PrintQuestionDetail(question *core.Question)
	PrintQuestionUpsertDetail(delta *core.Delta)
	PrintSuccess(message string)
	PrintError(err error)
	PrintCancel(message string)
}

type IOHandlerImpl struct {
	Reader io.Reader
	Writer io.Writer
	Clock  clock.Clock
}

func NewIOHandler(clock clock.Clock) *IOHandlerImpl {
	return &IOHandlerImpl{
		Reader: os.Stdin,
		Writer: os.Stdout,
		Clock:  clock,
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
	return strings.TrimSpace(scanner.Text())
}

func (ioh *IOHandlerImpl) PrintQuestionBrief(q *core.Question) {
	ioh.PrintfColored(ColorQuestionURL, "[%d] %s (Due: %s)\n", q.ID, q.URL, q.NextReview.Format("2006-01-02"))
	if q.Note == "" {
		ioh.Printf(" ↳ Note: (none)\n")
	} else {
		ioh.Printf(" ↳ Note: %s\n", q.Note)
	}
}

func (ioh *IOHandlerImpl) PrintQuestionDetail(question *core.Question) {
	ioh.PrintfColored(ColorQuestionURL, "[%d] %s\n", question.ID, question.URL)
	if question.Note == "" {
		ioh.Printf(" ↳ Note: (none)\n")
	} else {
		ioh.Printf(" ↳ Note: %s\n", question.Note)
	}
	ioh.Printf("   Familiarity: %d/%d\n", question.Familiarity+1, core.MaxFamiliarity)
	ioh.Printf("   Importance: %d/%d\n", question.Importance+1, core.MaxImportance)
	ioh.Printf("   Last Reviewed: %s\n", question.LastReviewed.Format("2006-01-02"))
	if question.NextReview.Before(ioh.Clock.Now()) {
		ioh.PrintfColored(ColorWarning, "   Next Review: %s (Due)\n", question.NextReview.Format("2006-01-02"))
	} else {
		ioh.Printf("   Next Review: %s\n", question.NextReview.Format("2006-01-02"))
	}
	ioh.Printf("   Review Count: %d\n", question.ReviewCount)
	ioh.Printf("   Ease Factor: %.2f\n", question.EaseFactor)
	ioh.Printf("   Created At: %s\n", question.CreatedAt.Format("2006-01-02"))
	ioh.Printf("\n")
}

func (ioh *IOHandlerImpl) PrintQuestionUpsertDetail(delta *core.Delta) {

	newState := delta.NewState
	oldState := delta.OldState

	ioh.PrintfColored(ColorQuestionURL, "[%d] %s\n", newState.ID, newState.URL)
	if newState.Note == "" {
		ioh.Printf(" ↳ Note: (none)\n")
	} else {
		ioh.Printf(" ↳ Note: %s\n", newState.Note)
	}

	if oldState == nil {
		ioh.Printf("   Familiarity: %d/%d\n", newState.Familiarity+1, core.MaxFamiliarity)
		ioh.Printf("   Importance: %d/%d\n", newState.Importance+1, core.MaxImportance)
		ioh.Printf("   Last Reviewed: %s\n", newState.LastReviewed.Format("2006-01-02"))
		ioh.Printf("   Next Review: %s\n", newState.NextReview.Format("2006-01-02"))
		ioh.Printf("   Review Count: %d\n", newState.ReviewCount)
		ioh.Printf("   Ease Factor: %.2f\n", newState.EaseFactor)
		ioh.Printf("   Created At: %s\n", newState.CreatedAt.Format("2006-01-02"))
		ioh.Printf("\n")
	} else {
		if oldState.Familiarity != newState.Familiarity {
			ioh.Printf("   Familiarity: %d → %d (Max: %d)\n", oldState.Familiarity+1, newState.Familiarity+1, core.MaxFamiliarity)
		} else {
			ioh.Printf("   Familiarity: %d/%d\n", newState.Familiarity+1, core.MaxFamiliarity)
		}
		if oldState.Importance != newState.Importance {
			ioh.Printf("   Importance: %d → %d (Max: %d)\n", oldState.Importance+1, newState.Importance+1, core.MaxImportance)
		} else {
			ioh.Printf("   Importance: %d/%d\n", newState.Importance+1, core.MaxImportance)
		}
		if oldState.LastReviewed != newState.LastReviewed {
			ioh.Printf("   Last Reviewed: %s → %s\n", oldState.LastReviewed.Format("2006-01-02"), newState.LastReviewed.Format("2006-01-02"))
		} else {
			ioh.Printf("   Last Reviewed: %s\n", newState.LastReviewed.Format("2006-01-02"))
		}
		if oldState.NextReview != newState.NextReview {
			ioh.Printf("   Next Review: %s → %s\n", oldState.NextReview.Format("2006-01-02"), newState.NextReview.Format("2006-01-02"))
		} else {
			ioh.Printf("   Next Review: %s\n", newState.NextReview.Format("2006-01-02"))
		}
		if oldState.ReviewCount != newState.ReviewCount {
			ioh.Printf("   Review Count: %d → %d\n", oldState.ReviewCount, newState.ReviewCount)
		} else {
			ioh.Printf("   Review Count: %d\n", newState.ReviewCount)
		}
		if oldState.EaseFactor != newState.EaseFactor {
			ioh.Printf("   Ease Factor: %.2f → %.2f\n", oldState.EaseFactor, newState.EaseFactor)
		} else {
			ioh.Printf("   Ease Factor: %.2f\n", newState.EaseFactor)
		}
		ioh.Printf("   Created At: %s\n", newState.CreatedAt.Format("2006-01-02"))
		ioh.Printf("\n")
	}

}

func (ioh *IOHandlerImpl) PrintCancel(message string) {
	ioh.PrintlnColored(ColorCancel, "[i] "+message)
}

func (ioh *IOHandlerImpl) PrintSuccess(message string) {
	ioh.PrintlnColored(ColorSuccess, "[✔] "+message)
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
			ioh.PrintlnColored(ColorWarning, "[!] "+codedErr.UserMessage())
			return
		case errs.BusinessErrorKind:
			// Business errors - show in yellow with user-friendly message
			ioh.PrintlnColored(ColorWarning, "[!] "+codedErr.UserMessage())
			return
		case errs.SystemErrorKind:
			// System errors - show in red with technical details
			ioh.PrintlnColored(ColorError, "[✘] "+codedErr.Error())
			return
		}
	}

	ioh.PrintlnColored(ColorError, "[✘] Error: "+err.Error())
}
