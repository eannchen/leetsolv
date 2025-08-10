package errs

type ErrorKind string

const (
	ValidationErrorKind ErrorKind = "VALIDATION ERROR" // Handler validation errors
	BusinessErrorKind   ErrorKind = "BUSINESS ERROR"   // Usecase business logic errors
	SystemErrorKind     ErrorKind = "SYSTEM ERROR"     // System/infrastructure errors
)

type CodedError struct {
	Err          error     // Technical error for debugging
	Kind         ErrorKind // Error category
	UserMsg      string    // User-friendly message for UI
	TechnicalMsg string    // Technical message for logs
}

func (e *CodedError) Error() string {
	// For debugging/logging - include technical details
	if e.TechnicalMsg != "" {
		if e.Err != nil {
			return string(e.Kind) + ": " + e.TechnicalMsg + " (" + e.Err.Error() + ")"
		}
		return string(e.Kind) + ": " + e.TechnicalMsg
	}
	if e.Err != nil {
		return string(e.Kind) + ": " + e.Err.Error()
	}
	return string(e.Kind)
}

// UserMessage returns the user-friendly message for UI display
func (e *CodedError) UserMessage() string {
	if e.UserMsg != "" {
		return e.UserMsg
	}
	// Fallback to technical message if no user message provided
	if e.Err != nil {
		return e.Err.Error()
	}
	return string(e.Kind)
}

func (e *CodedError) Unwrap() error {
	return e.Err
}

func WrapInternalError(err error, msg string) error {
	return &CodedError{
		Err:          err,
		Kind:         SystemErrorKind,
		TechnicalMsg: msg,
	}
}

func WrapValidationError(err error, msg string) error {
	return &CodedError{
		Err:     err,
		Kind:    ValidationErrorKind,
		UserMsg: msg,
	}
}

func WrapBusinessError(err error, msg string) error {
	return &CodedError{
		Err:     err,
		Kind:    BusinessErrorKind,
		UserMsg: msg,
	}
}
