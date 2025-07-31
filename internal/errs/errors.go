package errs

import "errors"

// 4xx – Client/Input errors
var (
	Err400QuestionNotFound     = WrapError(InputErrorKind, errors.New("question not found"), "question not found")
	Err400NoQuestionsAvailable = WrapError(InputErrorKind, errors.New("no questions available"), "no questions available")
	Err400InvalidPageNumber    = WrapError(InputErrorKind, errors.New("invalid page number"), "invalid page number")
	Err400NoActionsToUndo      = WrapError(InputErrorKind, errors.New("no actions to undo"), "no actions to undo")
)

// // 5xx – System/Internal errors
// var (
// 	Err500FailedToLoadDetails = errors.New("failed to load details")
// 	Err500DatabaseUnavailable = errors.New("database unavailable")
// )
