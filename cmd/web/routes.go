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

	router.Handler(http.MethodGet, "/", app.homeGet())
	router.Handler(http.MethodPost, "/", app.homePost())

	// TODO(Amr Ojjeh): Improve Adapter interface
	// idRequired := alice.New(app.idRequired)
	// contentExcerptRequired := idRequired.Append(app.excerptRequired)
	// grammarExcerptRequired := contentExcerptRequired.Append(app.contentLockRequired)
	// technicalExcerptRequired := grammarExcerptRequired.Append(app.grammarLockRequired)
	// technicalWordRequired := technicalExcerptRequired.Append(app.technicalWordRequired)

	// router.Handler(http.MethodGet, "/excerpt/edit",
	// 	contentExcerptRequired.Then(app.excerptEditGet()))
	// router.Handler(http.MethodPut, "/excerpt/edit",
	// 	idRequired.Then(app.excerptEditPut()))
	// router.Handler(http.MethodPut, "/excerpt/edit/lock",
	// 	contentExcerptRequired.Then(app.excerptEditLock()))
	// router.Handler(http.MethodPut, "/excerpt/edit/unlock",
	// 	idRequired.Then(app.excerptEditUnlock()))
	// router.Handler(http.MethodGet, "/excerpt/grammar",
	// 	grammarExcerptRequired.Then(app.excerptGrammarGet()))
	// router.Handler(http.MethodPut, "/excerpt/grammar",
	// 	grammarExcerptRequired.Then(app.excerptGrammarPut()))
	// router.Handler(http.MethodPut, "/excerpt/grammar/lock",
	// 	grammarExcerptRequired.Then(app.excerptGrammarLock()))
	// router.Handler(http.MethodPut, "/excerpt/grammar/unlock",
	// 	grammarExcerptRequired.Then(app.excerptGrammarUnlock()))

	// router.Handler(http.MethodGet, "/excerpt/technical",
	// 	technicalExcerptRequired.Then(app.excerptTechnicalGet()))
	// router.Handler(http.MethodPut, "/excerpt/technical/tashkeel",
	// 	technicalWordRequired.Then(app.excerptTechnicalVowelPut()))
	// router.Handler(http.MethodPut, "/excerpt/technical/shadda",
	// 	technicalWordRequired.Then(app.excerptTechnicalShadda()))
	// router.Handler(http.MethodGet, "/excerpt/technical/word",
	// 	technicalWordRequired.Then(app.excerptTechnicalWordGet()))
	// router.Handler(http.MethodPut, "/excerpt/technical/sentenceStart",
	// 	technicalWordRequired.Then(app.excerptTechnicalSentenceStart()))
	// router.Handler(http.MethodPut, "/excerpt/technical/ignore",
	// 	technicalWordRequired.Then(app.excerptTechnicalIgnore()))
	// router.Handler(http.MethodGet, "/excerpt/technical/export.json",
	// 	technicalExcerptRequired.Then(app.excerptTechnicalExport()))
	base := alice.New(app.logRequest)
	return base.Then(router)
}
