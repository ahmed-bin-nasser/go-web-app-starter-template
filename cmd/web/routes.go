package main

import (
	"net/http"

	"example.com/assets"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.NotFound(app.notFound)

	mux.Use(app.logAccess)
	mux.Use(app.recoverPanic)
	mux.Use(app.securityHeaders)
	mux.Use(app.getCsrfMiddleware())

	fileServer := http.FileServer(http.FS(assets.StaticFiles))
	mux.Handle("/static/*", fileServer)

	mux.Get("/", app.home)

	mux.Group(func(r chi.Router) {
		r.Use(app.authenticate)

		r.Get("/registration", app.registration)
		r.Post("/registration", app.postRegistration)

		r.Get("/login", app.login)
		r.Post("/login", app.postLogin)

		r.With(app.requireAuthentication).Get("/profile", app.profile)
		r.With(app.requireAuthentication).Post("/logout", app.postLogout)
	})

	return mux
}

func (app *application) getCsrfMiddleware() func(http.Handler) http.Handler {
	if app.config.AppEnv == "prod" {
		return csrf.Protect(
			[]byte(app.config.CsrfSecretKey),
			csrf.SameSite(csrf.SameSiteStrictMode),
		)
	}

	return csrf.Protect(
		[]byte(app.config.CsrfSecretKey),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.Secure(false),
	)
}
