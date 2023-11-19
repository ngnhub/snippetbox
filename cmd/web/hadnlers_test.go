package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	// given
	record := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	// when
	ping(record, request)

	// then
	result := record.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, result.StatusCode)
	}

	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
