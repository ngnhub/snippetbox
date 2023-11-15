package forms

import (
	"net/url"

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

func (f *Form) RequiredLengths(length int, fields ...string) {
	for _, field := range fields {
		val := f.Get(field)
		f.requiredLength(length, field, val)
	}
}

func (f *Form) requiredLength(length int, field, val string) {
	if val == "" {
		return
	}
	if message := validation.NoLongerThan(val, length); message != "" {
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

func (f *Form) IsValid() bool {
	return len(f.Errors) == 0
}
