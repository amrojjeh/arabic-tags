package main

import (
	"net/http"

	"github.com/amrojjeh/arabic-tags/ui"
	"github.com/julienschmidt/httprouter"
)

// TODO(Amr Ojjeh): Setup secure headers
func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = Adapt(app.notFound(), app.logRequest)

	router.Handler(http.MethodGet, "/static/*file", Adapt(http.FileServer(http.FS(ui.Files)),
		app.logRequest))

	router.Handler(http.MethodGet, "/", Adapt(app.excerptCreateGet(),
		app.logRequest))
	router.Handler(http.MethodPost, "/", Adapt(app.excerptCreatePost(),
		app.logRequest))

	// TODO(Amr Ojjeh): Improve Adapter interface
	idRequired := []Adapter{
		app.idRequired,
		app.logRequest,
	}

	contentExcerptRequired := []Adapter{
		app.excerptRequired,
		app.idRequired,
		app.logRequest,
	}

	grammarExcerptRequired := []Adapter{
		app.contentLockRequired,
		app.excerptRequired,
		app.idRequired,
		app.logRequest,
	}

	technicalExcerptRequired := []Adapter{
		app.grammarLockRequired,
		app.contentLockRequired,
		app.excerptRequired,
		app.idRequired,
		app.logRequest,
	}

	technicalWordRequired := []Adapter{
		app.technicalWordRequired,
		app.grammarLockRequired,
		app.contentLockRequired,
		app.excerptRequired,
		app.idRequired,
		app.logRequest,
	}

	router.Handler(http.MethodGet, "/excerpt/edit", Adapt(app.excerptEditGet(),
		contentExcerptRequired...))
	router.Handler(http.MethodPut, "/excerpt/edit", Adapt(app.excerptEditPut(),
		idRequired...))
	router.Handler(http.MethodPut, "/excerpt/edit/lock", Adapt(app.excerptEditLock(),
		contentExcerptRequired...))
	router.Handler(http.MethodPut, "/excerpt/edit/unlock", Adapt(app.excerptEditUnlock(),
		idRequired...))

	router.Handler(http.MethodGet, "/excerpt/grammar", Adapt(app.excerptGrammarGet(),
		grammarExcerptRequired...))
	router.Handler(http.MethodPut, "/excerpt/grammar", Adapt(app.excerptGrammarPut(),
		grammarExcerptRequired...))
	router.Handler(http.MethodPut, "/excerpt/grammar/lock", Adapt(app.excerptGrammarLock(),
		grammarExcerptRequired...))
	router.Handler(http.MethodPut, "/excerpt/grammar/unlock", Adapt(app.excerptGrammarUnlock(),
		grammarExcerptRequired...))

	router.Handler(http.MethodGet, "/excerpt/technical", Adapt(app.excerptTechnicalGet(),
		technicalExcerptRequired...))
	router.Handler(http.MethodPut, "/excerpt/technical/tashkeel", Adapt(app.excerptTechnicalVowelPut(),
		technicalWordRequired...))
	router.Handler(http.MethodPut, "/excerpt/technical/shadda", Adapt(app.excerptTechnicalShadda(),
		technicalWordRequired...))
	router.Handler(http.MethodGet, "/excerpt/technical/word", Adapt(app.excerptTechnicalWordGet(),
		technicalWordRequired...))
	router.Handler(http.MethodPut, "/excerpt/technical/sentenceStart",
		Adapt(app.excerptTechnicalSentenceStart(), technicalWordRequired...))
	router.Handler(http.MethodPut, "/excerpt/technical/ignore",
		Adapt(app.excerptTechnicalIgnore(), technicalWordRequired...))
	router.Handler(http.MethodGet, "/excerpt/technical/export.json",
		Adapt(app.excerptTechnicalExport(), technicalExcerptRequired...))
	return router
}
