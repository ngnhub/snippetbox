package forms

import (
	"net/url"
	"regexp"

	"github.com/ngnhub/snippetbox/pkg/models/validation"
)

type Form struct {
	url.Values
	Errors validation.Errors
}

func New(val url.Values) *Form {
	return &Form{
		Values: val,
		Errors: make(validation.Errors),
	}
}

func (f *Form) Requried(fields ...string) {
	for _, field := range fields {
		val := f.Get(field)
		if message := validation.NotEmpty(val); message != "" {
			f.Errors.Add(field, message)
		}
	}
}

func (f *Form) MaxLength(length int, fields ...string) {
	for _, field := range fields {
		val := f.Get(field)
		f.maxLength(length, field, val)
	}
}

func (f *Form) maxLength(length int, field, val string) {
	if val == "" {
		return
	}
	if message := validation.MaxLength(val, length); message != "" {
		f.Errors.Add(field, message)
	}
}

func (f *Form) PermittedValues(values []string, fields ...string) {
	for _, field := range fields {
		val := f.Get(field)
		f.permittedValue(values, field, val)
	}
}

func (f *Form) permittedValue(values []string, field, val string) {
	if val == "" {
		return
	}
	if message := validation.NotContainsIn(val, values...); message != "" {
		f.Errors.Add(field, message)
	}
}

func (f *Form) MinLength(length int, fields ...string) {
	for _, field := range fields {
		val := f.Get(field)
		f.minLength(length, field, val)
	}
}

func (f *Form) minLength(length int, field, val string) {
	if val == "" {
		return
	}
	if message := validation.MinLength(val, length); message != "" {
		f.Errors.Add(field, message)
	}
}

func (f *Form) MatchPattern(field string, regex *regexp.Regexp) {
	val := f.Get(field)
	if val == "" {
		return
	}
	if message := validation.MatchPattern(val, regex); message != "" {
		f.Errors.Add(field, message)
	}
}

func (f *Form) IsValid() bool {
	return len(f.Errors) == 0
}
