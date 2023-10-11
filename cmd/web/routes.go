package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TODO(Amr Ojjeh): Setup secure headers
func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = Adapt(app.notFound(), app.logRequest)

	router.Handler(http.MethodGet, "/static/*file", Adapt(http.FileServer(http.Dir("./ui/static/")),
		stripPrefix("/static"), app.logRequest))

	router.Handler(http.MethodGet, "/", Adapt(app.excerptCreateGet(),
		app.logRequest))
	router.Handler(http.MethodPost, "/", Adapt(app.excerptCreatePost(),
		app.logRequest))

	router.Handler(http.MethodGet, "/excerpt/edit", Adapt(app.excerptEditGet(),
		app.logRequest))
	router.Handler(http.MethodPut, "/excerpt/edit", Adapt(app.excerptEditPut(),
		app.logRequest))
	router.Handler(http.MethodPut, "/excerpt/edit/lock", Adapt(app.excerptEditLock(),
		app.logRequest))

	router.Handler(http.MethodGet, "/excerpt/grammar", Adapt(app.excerptGrammarGet(),
		app.logRequest))
	router.Handler(http.MethodPut, "/excerpt/grammar", Adapt(app.excerptGrammarPut(),
		app.logRequest))
	return router
}
