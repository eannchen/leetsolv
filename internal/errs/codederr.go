package errs

type ErrorKind string

const (
	InputErrorKind  ErrorKind = "INPUT"
	SystemErrorKind ErrorKind = "SYSTEM"
)

type CodedError struct {
	Err     error
	Kind    ErrorKind
	Message string
}

func (e *CodedError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *CodedError) Unwrap() error {
	return e.Err
}

func WrapError(kind ErrorKind, err error, msg string) error {
	return &CodedError{
		Err:     err,
		Kind:    kind,
		Message: msg,
	}
}

func WrapInternalError(err error, msg string) error {
	return WrapError(SystemErrorKind, err, msg)
}

func WrapClientError(err error, msg string) error {
	return WrapError(InputErrorKind, err, msg)
}
