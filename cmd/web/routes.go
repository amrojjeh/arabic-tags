package main

import (
	"fmt"
	"net/http"

	"github.com/amrojjeh/arabic-tags/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type url struct{}

func (u url) index() string {
	return "/"
}

func (u url) login() string {
	return "/login"
}

func (u url) register() string {
	return "/register"
}

func (u url) logout() string {
	return "/logout"
}

func (u url) home() string {
	return "/home"
}

func (u url) createExcerpt() string {
	return "/excerpt"
}

func (u url) excerpt(id any) string {
	return fmt.Sprintf("/excerpt/%v", id)
}

func (u url) excerptTitle(id any) string {
	return fmt.Sprintf("/excerpt/%v/title", id)
}

// TODO(Amr Ojjeh): Setup secure headers
func (app *application) routes() http.Handler {
	app.u = url{}
	router := httprouter.New()
	router.NotFound = app.notFound()

	router.Handler(http.MethodGet, "/static/*file",
		http.FileServer(http.FS(ui.Files)))

	router.Handler(http.MethodGet, app.u.register(), app.registerGet())
	router.Handler(http.MethodPost, app.u.register(), app.registerPost())
	router.Handler(http.MethodGet, app.u.login(), app.loginGet())
	router.Handler(http.MethodPost, app.u.login(), app.loginPost())

	authRequired := alice.New(app.authRequired)
	// TODO(Amr Ojjeh): Write an index page
	router.Handler(http.MethodGet, app.u.index(), authRequired.Then(app.homeGet()))
	router.Handler(http.MethodPost, app.u.logout(), authRequired.Then(app.logoutPost()))
	router.Handler(http.MethodGet, app.u.home(), authRequired.Then(app.homeGet()))
	router.Handler(http.MethodGet, app.u.createExcerpt(), authRequired.Then(app.createExcerptGet()))
	router.Handler(http.MethodPost, app.u.createExcerpt(), authRequired.Then(app.createExcerptPost()))

	excerptRequired := alice.New(app.excerptRequired)
	router.Handler(http.MethodGet, app.u.excerpt(":id"), excerptRequired.Then(app.excerptGet()))
	router.Handler(http.MethodPost, app.u.excerpt(":id"), excerptRequired.Then(app.excerptPost()))
	router.Handler(http.MethodGet, app.u.excerptTitle(":id"), excerptRequired.Then(app.excerptTitleGet()))
	router.Handler(http.MethodPost, app.u.excerptTitle(":id"), excerptRequired.Then(app.excerptTitlePost()))

	base := alice.New(app.session.LoadAndSave, app.recoverPanic, app.logRequest)
	return base.Then(router)
}
