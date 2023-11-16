package models

import (
	"errors"
	"time"
)

var ErorNoRecord = errors.New("models: no matching errors found")
var ErrorInvalidCredentials = errors.New("models: invalid credentials")
var ErrorDuplicatedEmail = errors.New("models: duplicate email")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}
