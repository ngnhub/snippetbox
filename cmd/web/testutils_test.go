package main

import (
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

func getTestApp() application {
	return application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}
}

type testServer struct {
	*httptest.Server
}

func getTestServer(t *testing.T, handler http.Handler) *testServer {
	server := httptest.NewTLSServer(handler)
	// Initialize a new cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	server.Client().Jar = jar
	// disable redirect
	server.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{Server: server}
}

func (server *testServer) getRequst(t *testing.T, path string) *http.Response {
	response, err := server.Client().Get(server.URL + path)
	if err != nil {
		t.Fatal(err)
	}
	return response
}
