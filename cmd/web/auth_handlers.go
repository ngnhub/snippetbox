package main

import (
	"errors"
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
	app.renderTemplate(writer, request, "login.page.tmpl", &templateData{Form: forms.New(nil)})
}

func (app *application) logIn(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}
	form := forms.New(request.PostForm)

	id, err := app.user.Authencticate(form.Get("email"), form.Get("password"))

	if err != nil {
		if errors.Is(models.ErrorInvalidCredentials, err) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.renderTemplate(writer, request, "login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(writer, err)
		}
		return
	}

	// saving user id in the session
	app.session.Put(request, AuthIdKey, id)
	http.Redirect(writer, request, "/snippet/create", http.StatusSeeOther)
}

func (app *application) logOut(writer http.ResponseWriter, request *http.Request) {
	app.session.Remove(request, AuthIdKey)
	app.session.Put(request, "flash", "You've been logged out successfully!")
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}
