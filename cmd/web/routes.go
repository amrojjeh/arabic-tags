package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TODO(Amr Ojjeh): Setup secure headers
func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.Handler(http.MethodGet, "/static/*file", Adapt(http.FileServer(http.Dir("./ui/static/")),
		stripPrefix("/static"), app.logRequest))

	router.Handler(http.MethodGet, "/", Adapt(app.excerptCreateGet(),
		app.logRequest))
	router.Handler(http.MethodPost, "/", Adapt(app.excerptCreatePost(),
		app.logRequest))

	router.Handler(http.MethodGet, "/excerpt/edit", Adapt(app.excerptGet(),
		app.logRequest))
	router.Handler(http.MethodPut, "/excerpt/edit", Adapt(app.excerptPut(),
		app.logRequest))
	return router
}
