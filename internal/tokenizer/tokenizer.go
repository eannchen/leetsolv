// Package tokenizer implements the tokenizer for the leetsolv application.
package tokenizer

import (
	"strings"
	"unicode"
)

func Tokenize(text string) []string {
	text = strings.ToLower(text)
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	return words
}
