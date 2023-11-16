package validation

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)

func NotEmpty(val string) string {
	var message string
	if strings.TrimSpace(val) == "" {
		message = "This field cannot be blank"
	}
	return message
}

func MaxLength(val string, length int) string {
	var message string
	if utf8.RuneCountInString(val) > 100 {
		message = fmt.Sprintf("This field cannot be longer than %d", length)
	}
	return message
}

func MinLength(val string, length int) string {
	var message string
	if utf8.RuneCountInString(val) < length {
		message = fmt.Sprintf("This field cannot be less than %d", length)
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

func MatchPattern(val string, exp *regexp.Regexp) string {
	var message string
	if !exp.MatchString(val) {
		message = "The field is invalid"
	}
	return message
}
