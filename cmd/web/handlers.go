package main

import (
	"net/http"
)

func (app *application) homeGet(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, "home.tmpl", templateData{})
}

func (app *application) excerptGet(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, "add.tmpl", templateData{})
}

func (app *application) excerptPost(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/excerpt", http.StatusSeeOther)
}
