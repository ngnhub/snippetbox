package main

import "net/http"

const AuthIdKey = "authenticationUserId"

func (app *application) IsAuthenticated(request *http.Request) bool {
	return app.session.Exists(request, AuthIdKey)
}
