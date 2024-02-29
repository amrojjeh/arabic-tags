package main

import (
	"net/http"

	"github.com/amrojjeh/arabic-tags/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// TODO(Amr Ojjeh): Setup secure headers
func (app *application) routes() http.Handler {
	app.u = url{}
	router := httprouter.New()
	router.NotFound = app.notFound()

	router.Handler(http.MethodGet, "/static/*file",
		http.FileServer(http.FS(ui.Files)))

	dynamic := alice.New(
		app.session.LoadAndSave,
		app.getUser,
	)

	router.Handler(http.MethodGet, app.u.register(), dynamic.Then(app.registerGet()))
	router.Handler(http.MethodPost, app.u.register(), dynamic.Then(app.registerPost()))
	router.Handler(http.MethodGet, app.u.login(), dynamic.Then(app.loginGet()))
	router.Handler(http.MethodPost, app.u.login(), dynamic.Then(app.loginPost()))

	authRequired := dynamic.Append(app.authRequired)
	// TODO(Amr Ojjeh): Write an index page
	router.Handler(http.MethodGet, app.u.index(), authRequired.Then(app.homeGet()))
	router.Handler(http.MethodPost, app.u.logout(), authRequired.Then(app.logoutPost()))
	router.Handler(http.MethodGet, app.u.home(), authRequired.Then(app.homeGet()))
	router.Handler(http.MethodGet, app.u.excerptCreate(), authRequired.Then(app.excerptCreateGet()))
	router.Handler(http.MethodPost, app.u.excerptCreate(), authRequired.Then(app.excerptCreatePost()))

	excerptRequired := dynamic.Append(app.excerptRequired)
	router.Handler(http.MethodGet, app.u.excerpt(":id"), excerptRequired.Then(app.excerptGet()))
	router.Handler(http.MethodGet, app.u.excerptExport(":id"), excerptRequired.Then(app.excerptExportGet()))

	ownerRequired := excerptRequired.Extend(authRequired).Append(app.ownerOfExcerpt)
	router.Handler(http.MethodPost, app.u.excerpt(":id"), ownerRequired.Then(app.excerptPost()))
	router.Handler(http.MethodGet, app.u.excerptTitle(":id"), ownerRequired.Then(app.excerptTitleGet()))
	router.Handler(http.MethodPost, app.u.excerptTitle(":id"), ownerRequired.Then(app.excerptTitlePost()))
	router.Handler(http.MethodPost, app.u.excerptLock(":id"), ownerRequired.Then(app.excerptNextPost()))
	router.Handler(http.MethodGet, app.u.wordEdit(":id"), ownerRequired.Then(app.wordEditGet()))
	router.Handler(http.MethodPost, app.u.wordEdit(":id"), ownerRequired.Then(app.wordEditPost()))

	ownerRequired = ownerRequired.Append(app.wordIdRequired)
	router.Handler(http.MethodPost, app.u.wordRight(":id", ":wid"), ownerRequired.Then(app.wordRightPost()))
	router.Handler(http.MethodPost, app.u.wordLeft(":id", ":wid"), ownerRequired.Then(app.wordLeftPost()))
	router.Handler(http.MethodPost, app.u.wordAdd(":id", ":wid"), ownerRequired.Then(app.wordAddPost()))
	router.Handler(http.MethodPost, app.u.wordRemove(":id", ":wid"), ownerRequired.Then(app.wordRemovePost()))
	router.Handler(http.MethodPost, app.u.wordConnect(":id", ":wid"), ownerRequired.Then(app.wordConnectPost()))
	router.Handler(http.MethodPost, app.u.wordSentenceStart(":id", ":wid"), ownerRequired.Then(app.wordSentenceStartPost()))
	router.Handler(http.MethodPost, app.u.wordIgnore(":id", ":wid"), ownerRequired.Then(app.wordIgnorePost()))
	router.Handler(http.MethodPost, app.u.wordCase(":id", ":wid"), ownerRequired.Then(app.wordCasePost()))
	router.Handler(http.MethodPost, app.u.wordState(":id", ":wid"), ownerRequired.Then(app.wordStatePost()))

	ownerRequired = ownerRequired.Append(app.letterPosRequired)
	router.Handler(http.MethodPost, app.u.letterEdit(":id", ":wid", ":lid"), ownerRequired.Then(app.letterEditPost()))

	base := alice.New(
		app.recoverPanic,
		app.logRequest,
		app.secureHeaders,
	)
	return base.Then(router)
}
