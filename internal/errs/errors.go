package errs

import "errors"

// Business logic errors (from usecase)
var (
	ErrQuestionNotFound     = WrapBusinessError(errors.New("question not found"), "Question not found. Please check the ID or URL")
	ErrNoQuestionsAvailable = WrapBusinessError(errors.New("no questions available"), "No questions available yet")
	ErrInvalidPageNumber    = WrapBusinessError(errors.New("invalid page number"), "Invalid page number")
	ErrNoActionsToUndo      = WrapBusinessError(errors.New("no actions to undo"), "No actions to undo")
)

// Validation errors (from handler)
var (
	ErrInvalidURLFormat         = WrapValidationError(errors.New("invalid URL format"), "Please provide a valid URL")
	ErrInvalidLeetCodeURL       = WrapValidationError(errors.New("URL must be from leetcode.com/problems/"), "Please provide a valid LeetCode problem URL")
	ErrInvalidLeetCodeURLFormat = WrapValidationError(errors.New("invalid LeetCode problem URL format"), "Please provide a valid LeetCode problem URL format")
	ErrInvalidFamiliarityLevel  = WrapValidationError(errors.New("invalid familiarity level"), "Please enter a familiarity level between 1 and 5")
	ErrInvalidImportanceLevel   = WrapValidationError(errors.New("invalid importance level"), "Please enter an importance level between 1 and 4")
)
