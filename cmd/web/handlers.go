package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ngnhub/snippetbox/pkg/models"
	"github.com/ngnhub/snippetbox/pkg/models/validation"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.GetLatest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := templateData{Snippets: s}
	app.renderTemplate(w, "home.page.tmpl", &data)
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErorNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := templateData{Snippet: snippet}
	app.renderTemplate(w, "show.page.tmpl", &data)
}

// Add a new createSnippetForm handler, which for now returns a placeholder response.
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, "create.page.tmpl", nil)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	errors := make(map[string]string)
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")
	if errorMessage := validateTitle(title); errorMessage != "" {
		errors["title"] = errorMessage
	}
	if errorMessage := validateContent(content); errorMessage != "" {
		errors["content"] = errorMessage
	}
	if errorMessage := validateExpires(expires); errorMessage != "" {
		errors["expires"] = errorMessage
	}

	if len(errors) > 0 {
		app.renderTemplate(w, "create.page.tmpl", &templateData{
			FormData:   r.PostForm,
			FormErrors: errors,
		})
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func validateTitle(title string) string {
	var message string
	if message = validation.NotEmpty(title); message != "" {
		return message
	}
	message = validation.NoLongerThan(title, 100)
	return message
}

func validateContent(content string) string {
	return validation.NotEmpty(content)
}

func validateExpires(expires string) string {
	var message string
	if message = validation.NotEmpty(expires); message != "" {
		return message
	}
	allowed := []string{"365", "7", "1"}
	return validation.NotContainsIn(expires, allowed...)
}
