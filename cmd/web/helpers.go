package main

import (
	"log/slog"
	"net/http"
)

func (app *application) clientError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	app.logger.Error("server error", slog.String("error", err.Error()))
	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func (app *application) noBody(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func (app *application) excerptNotFound(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?error=Excerpt not found", http.StatusSeeOther)
}
