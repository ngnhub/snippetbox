package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ngnhub/snippetbox/pkg/models"
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
	w.Write([]byte("Create a new snippet..."))
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// temporal dummy data
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
