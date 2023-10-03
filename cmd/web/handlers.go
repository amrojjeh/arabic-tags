package main

import (
	"net/http"
)

// TODO(Amr Ojjeh): Return error if there's no id param
func (app *application) excerptGet(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, "add.tmpl", http.StatusOK, templateData{})
}

type excerptForm struct {
	Validator
	Title string
}

func (app *application) excerptCreate(w http.ResponseWriter, r *http.Request) {
	data := templateData{}
	data.Form = excerptForm{}
	app.renderTemplate(w, "home.tmpl", http.StatusOK, data)
}

func (app *application) excerptCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := excerptForm{}

	form.Title = r.Form.Get("title")
	form.CheckField(NotBlank(form.Title),
		"title", "Title cannot be blank")
	form.CheckField(MaxChars(form.Title, 100), "title",
		"Title cannot exceed 100 characters")

	if !form.Valid() {
		data := templateData{
			Form: form,
		}

		if r.Header.Get("HX-Boosted") == "true" {
			app.renderTemplate(w, "home.tmpl", http.StatusOK, data)
		} else {
			app.renderTemplate(w, "home.tmpl", http.StatusUnprocessableEntity, data)
		}
		return
	}

	// TODO(Amr Ojjeh): Insert and then redirect with proper id
	http.Redirect(w, r, "/excerpt", http.StatusSeeOther)
}
