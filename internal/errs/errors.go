package errs

import "errors"

// 4xx – Client/Input errors
var (
	Err400QuestionNotFound     = WrapError(InputErrorKind, errors.New("question not found"), "question not found")
	Err400NoQuestionsAvailable = WrapError(InputErrorKind, errors.New("no questions available"), "no questions available")
	Err400InvalidPageNumber    = WrapError(InputErrorKind, errors.New("invalid page number"), "invalid page number")
	Err400NoActionsToUndo      = WrapError(InputErrorKind, errors.New("no actions to undo"), "no actions to undo")
)

// Handler validation errors
var (
	ErrInvalidURLFormat         = WrapError(InputErrorKind, errors.New("invalid URL format"), "Please provide a valid URL")
	ErrInvalidLeetCodeURL       = WrapError(InputErrorKind, errors.New("URL must be from leetcode.com/problems/"), "Please provide a valid LeetCode problem URL")
	ErrInvalidLeetCodeURLFormat = WrapError(InputErrorKind, errors.New("invalid LeetCode problem URL format"), "Please provide a valid LeetCode problem URL format")
	ErrInvalidFamiliarityLevel  = WrapError(InputErrorKind, errors.New("invalid familiarity level"), "Please enter a familiarity level between 1 and 5")
	ErrInvalidImportanceLevel   = WrapError(InputErrorKind, errors.New("invalid importance level"), "Please enter an importance level between 1 and 4")
)

// // 5xx – System/Internal errors
// var (
// 	Err500FailedToLoadDetails = errors.New("failed to load details")
// 	Err500DatabaseUnavailable = errors.New("database unavailable")
// )
