package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TODO(Amr Ojjeh): Setup request logging
// TODO(Amr Ojjeh): Setup secure headers
func (app *application) routes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*file",
		http.StripPrefix("/static", fileServer))
	router.HandlerFunc(http.MethodGet, "/", app.excerptCreate)
	router.HandlerFunc(http.MethodPost, "/", app.excerptCreatePost)
	router.HandlerFunc(http.MethodGet, "/excerpt", app.excerptGet)
	return router
}
