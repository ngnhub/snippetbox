package main

import (
	"net/http"
)

const AuthIdKey = "authenticationUserId"

func (app *application) IsAuthenticated(request *http.Request) bool {
	exist, ok := request.Context().Value(AuthKey).(bool)
	if !ok {
		return false
	}
	return exist
}
