package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ngnhub/snippetbox/pkg/models"
	"github.com/ngnhub/snippetbox/pkg/models/forms"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.GetLatest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := templateData{Snippets: s}
	app.renderTemplate(w, r, "home.page.tmpl", &data)
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

	data := templateData{
		Snippet: snippet,
	}
	app.renderTemplate(w, r, "show.page.tmpl", &data)
}

// Add a new createSnippetForm handler, which for now returns a placeholder response.
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil)})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Requried("title", "content", "expires")
	form.RequiredLengths(100, "title")
	form.PermittedValues([]string{"365", "7", "1"}, "expires")

	if !form.IsValid() {
		app.renderTemplate(w, r, "create.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
