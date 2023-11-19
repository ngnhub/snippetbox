package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// given
	data := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			"UTC",
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			"01 Jan 2024 at 00:00",
		},
		{
			"Empty",
			time.Time{},
			"",
		},
	}

	for _, elem := range data {
		t.Run(elem.name, func(t *testing.T) {
			// when
			actual := humanDate(elem.time)

			// then
			if actual != elem.expected {
				t.Errorf("Test name: %s. Expected %s, but was %s", elem.name, elem.expected, actual)
			}
		})
	}
}
