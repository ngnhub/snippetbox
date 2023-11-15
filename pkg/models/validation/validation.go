package validation

import (
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"
)

func NotEmpty(val string) string {
	var message string
	if strings.TrimSpace(val) == "" {
		message = "This field cannot be blank"
	}
	return message
}

func NoLongerThan(val string, length int) string {
	var message string
	if utf8.RuneCountInString(val) > 100 {
		message = fmt.Sprintf("This field cannot be longer than %d", length)
	}
	return message
}

func NotContainsIn(val string, in ...string) string {
	var message string
	if !slices.Contains(in, val) {
		message = "The field is invalid"
	}
	return message
}
