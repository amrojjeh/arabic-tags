package main

import (
	"net/http"

	"github.com/amrojjeh/arabic-tags/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// TODO(Amr Ojjeh): Setup secure headers
func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = app.notFound()

	router.Handler(http.MethodGet, "/static/*file",
		http.FileServer(http.FS(ui.Files)))

	router.Handler(http.MethodGet, "/register", app.registerGet())
	router.Handler(http.MethodPost, "/register", app.registerPost())
	router.Handler(http.MethodGet, "/login", app.loginGet())
	router.Handler(http.MethodPost, "/login", app.loginPost())

	authRequired := alice.New(app.authRequired)
	// TODO(Amr Ojjeh): Write an index page
	router.Handler(http.MethodGet, "/", authRequired.Then(app.homeGet()))
	router.Handler(http.MethodPost, "/logout", authRequired.Then(app.logoutPost()))
	router.Handler(http.MethodGet, "/home", authRequired.Then(app.homeGet()))
	router.Handler(http.MethodGet, "/excerpt", authRequired.Then(app.createExcerptGet()))
	router.Handler(http.MethodPost, "/excerpt", authRequired.Then(app.createExcerptPost()))

	excerptRequired := alice.New(app.excerptRequired)
	router.Handler(http.MethodGet, "/excerpt/:id", excerptRequired.Then(app.excerptGet()))
	router.Handler(http.MethodPost, "/excerpt/:id", excerptRequired.Then(app.excerptPost()))
	router.Handler(http.MethodGet, "/excerpt/:id/title", excerptRequired.Then(app.excerptTitleGet()))
	router.Handler(http.MethodPost, "/excerpt/:id/title", excerptRequired.Then(app.excerptTitlePost()))

	base := alice.New(app.session.LoadAndSave, app.recoverPanic, app.logRequest)
	return base.Then(router)
}
