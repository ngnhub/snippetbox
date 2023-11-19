package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/ngnhub/snippetbox/pkg/models"
)

const XssProtectionHeader = "X-XSS-Protection"
const XssFrameOptionHeader = "X-Frame-Options"
const ConnectionHeader = "Connection"

func addSecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(XssProtectionHeader, "1; mode=block")
		w.Header().Set(XssFrameOptionHeader, "deny")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) handlePanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set(ConnectionHeader, "closed")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticated := app.IsAuthenticated(r)

		if !authenticated {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache (or
		// other intermediary cache).
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

// Create a NoSurf middleware function which uses a customized CSRF cookie with
// the Secure, Path and HttpOnly flags set.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticated := app.session.Exists(r, AuthIdKey)
		if !authenticated {
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.user.GetBy(app.session.GetInt(r, AuthIdKey))
		if err != nil {
			if errors.Is(err, models.ErorNoRecord) {
				app.session.Remove(r, AuthIdKey)
				next.ServeHTTP(w, r)
			} else {
				app.serverError(w, err)
			}
			return
		}

		if !user.Active {
			app.session.Remove(r, AuthIdKey)
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), AuthKey, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
