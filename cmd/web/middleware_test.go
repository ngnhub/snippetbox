package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddSecureHeaders(t *testing.T) {
	// given
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))

	})

	// when
	addSecureHeaders(next).ServeHTTP(recorder, request)

	// then
	response := recorder.Result()
	frameOpt := response.Header.Get(XssFrameOptionHeader)
	if frameOpt != "deny" {
		t.Errorf("want %q; got %q", "deny", frameOpt)
	}

	protheader := response.Header.Get(XssProtectionHeader)
	if protheader != "1; mode=block" {
		t.Errorf("want %q; got %q", "1; mode=block", protheader)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, response.StatusCode)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
