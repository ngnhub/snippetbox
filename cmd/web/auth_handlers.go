package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ngnhub/snippetbox/pkg/models"
	"github.com/ngnhub/snippetbox/pkg/models/forms"
	"github.com/ngnhub/snippetbox/pkg/models/validation"
)

func (app *application) getSignupForm(writer http.ResponseWriter, request *http.Request) {
	app.renderTemplate(writer, request, "signup.page.tmpl", &templateData{Form: forms.New(nil)})
}

func (app *application) signUp(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	form := forms.New(request.PostForm)
	form.Requried("name", "email", "password")
	form.MinLength(10, "password")
	form.MaxLength(255, "name", "email")
	form.MatchPattern("email", validation.EmailRX)

	if !form.IsValid() {
		app.renderTemplate(writer, request, "signup.page.tmpl", &templateData{Form: form})
	}
	err := app.user.Insert(form.Get("email"), form.Get("name"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrorDuplicatedEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.renderTemplate(writer, request, "signup.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(writer, err)
		}
		return
	}

	app.session.Put(request, "flash", "Your signup was successful. Please log in.")
	http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
}

func (app *application) getLoginForm(writer http.ResponseWriter, request *http.Request) {

	fmt.Fprintln(writer, "not supported yet")
}

func (app *application) logIn(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "not supported yet")
}

func (app *application) logOut(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "not supported yet")
}
