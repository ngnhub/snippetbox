package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data *templateData) {
	template, ok := app.templateCache[templateName]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", templateName))
		return
	}
	buffer := new(bytes.Buffer)
	data = app.addData(data, r)
	err := template.Execute(buffer, data)
	if err != nil {
		app.serverError(w, err)
	}

	buffer.WriteTo(w)
}

func (app *application) addData(data *templateData, r *http.Request) *templateData {
	if data == nil {
		data = &templateData{}
	}
	flash := app.session.PopString(r, "flash")
	data.Flash = flash
	data.CurrentYear = time.Now().Year()
	data.IsAuthenticated = app.IsAuthenticated(r)
	data.CSRFToken = nosurf.Token(r)
	return data
}
