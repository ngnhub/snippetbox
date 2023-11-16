package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Request middleware chain
	standardMidleware := alice.New(app.handlePanic, app.logRequest, addSecureHeaders)

	sessionMidleware := alice.New(app.session.Enable, noSurf)

	// Create handlers aka Controllers
	mux := pat.New()
	mux.Get("/", sessionMidleware.ThenFunc(app.home))
	mux.Get("/snippet/create", sessionMidleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", sessionMidleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", sessionMidleware.ThenFunc(app.showSnippet))

	mux.Get("/user/signup", sessionMidleware.ThenFunc(app.getSignupForm))
	mux.Post("/user/signup", sessionMidleware.ThenFunc(app.signUp))
	mux.Get("/user/login", sessionMidleware.ThenFunc(app.getLoginForm))
	mux.Post("/user/login", sessionMidleware.ThenFunc(app.logIn))
	mux.Post("/user/logout", sessionMidleware.Append(app.requireAuthentication).ThenFunc(app.logOut))

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMidleware.Then(mux)
}
