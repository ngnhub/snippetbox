package models

import (
	"errors"
	"time"
)

var ErorNoRecord = errors.New("models: no matching errors found")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
