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

	idRequired := []Adapter{
		app.idRequired,
		app.logRequest,
	}

	excerptRequired := []Adapter{
		app.excerptRequired,
		app.idRequired,
		app.logRequest,
	}

	contentLockRequired := []Adapter{
		app.contentLockRequired,
		app.excerptRequired,
		app.idRequired,
		app.logRequest,
	}

	router.Handler(http.MethodGet, "/excerpt/edit", Adapt(app.excerptEditGet(),
		excerptRequired...))
	router.Handler(http.MethodPut, "/excerpt/edit", Adapt(app.excerptEditPut(),
		idRequired...))
	router.Handler(http.MethodPut, "/excerpt/edit/lock", Adapt(app.excerptEditLock(),
		excerptRequired...))
	router.Handler(http.MethodPut, "/excerpt/edit/unlock", Adapt(app.excerptEditUnlock(),
		idRequired...))

	router.Handler(http.MethodGet, "/excerpt/grammar", Adapt(app.excerptGrammarGet(),
		contentLockRequired...))
	router.Handler(http.MethodPut, "/excerpt/grammar", Adapt(app.excerptGrammarPut(),
		contentLockRequired...))
	return router
}
