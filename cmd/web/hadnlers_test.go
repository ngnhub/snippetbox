package main

import (
	"io"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	// given
	app := getTestApp()
	server := getTestServer(t, app.routes())
	defer server.Close()

	// when
	response := server.getRequst(t, "/ping")

	// then
	if response.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, response.StatusCode)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
