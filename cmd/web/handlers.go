package main

import (
	"errors"
	"net/http"
	"time"

	"example.com/internal/database"
	"example.com/pkg/request"
	"github.com/gorilla/sessions"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w, r)
		return
	}

	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "home", data)
}

func (app *application) profile(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "profile", data)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "login", data)
}

func (app *application) registration(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "registration", data)
}

func (app *application) postRegistration(w http.ResponseWriter, r *http.Request) {
	var form RegistrationForm

	err := request.DecodePostForm(r, &form)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if !form.isValid() {
		data := app.newTemplateData(r)
		data["Form"] = form
		app.render(w, r, http.StatusUnprocessableEntity, "registration", data)
		return
	}

	_, err = app.db.RegisterUser(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, database.ErrDuplicateEmail) {
			form.AddError("sorry, the email has been taken")
			data := app.newTemplateData(r)
			data["Form"] = form
			app.render(w, r, http.StatusUnauthorized, "registration", data)
			return
		}

		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	app.render(w, r, http.StatusSeeOther, "login", data)
}

func (app *application) postLogin(w http.ResponseWriter, r *http.Request) {
	var form LoginForm

	err := request.DecodePostForm(r, &form)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if !form.isValid() {
		data := app.newTemplateData(r)
		data["Form"] = form
		app.render(w, r, http.StatusUnprocessableEntity, "login", data)
		return
	}

	user, err := app.db.AuthenticateUser(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, database.ErrInvalidCredentials) {
			form.AddError("invalid credential")
			data := app.newTemplateData(r)
			data["Form"] = form
			app.render(w, r, http.StatusUnauthorized, "login", data)
			return
		}

		app.serverError(w, r, err)
		return
	}

	session, err := app.sessionStore.Get(r, sessionsKey)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	session.Values[authenticatedUserIDKey] = user.ID
	session.Values[timestampKey] = time.Now()
	session.Options = &sessions.Options{
		HttpOnly: true,
		MaxAge:   86400 * app.config.CookieMaxAge,
	}

	err = session.Save(r, w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (app *application) postLogout(w http.ResponseWriter, r *http.Request) {
	session, err := app.sessionStore.Get(r, sessionsKey)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	session.Values[authenticatedUserIDKey] = uint64(0)
	session.Values[timestampKey] = time.Time{}
	session.Options = &sessions.Options{
		MaxAge:   -1, // Set MaxAge to -1 to delete the cookie
		HttpOnly: true,
	}

	err = session.Save(r, w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	app.render(w, r, http.StatusSeeOther, "login", data)
}
