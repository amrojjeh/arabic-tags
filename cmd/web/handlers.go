package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/amrojjeh/arabic-tags/internal/models"
	"github.com/google/uuid"
)

func (app *application) excerptGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		idStr := r.Form.Get("id")
		if idStr == "" {
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		excerpt, err := app.excerpts.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				// TODO(Amr Ojjeh): Write an excerpt not found page
				w.Write([]byte("Excerpt not found..."))
				return
			}
			app.serverError(w, err)
			return
		}
		data := templateData{
			Excerpt: excerpt,
		}
		app.renderTemplate(w, "add.tmpl", http.StatusOK, data)
	})
}

func (app *application) excerptPut() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		idStr := r.Form.Get("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		content := r.Form.Get("content")
		app.excerpts.UpdateContent(id, content)
		app.noBody(w)

	})
}

type excerptForm struct {
	Validator
	Title string
}

func (app *application) excerptCreateGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := templateData{}
		data.Form = excerptForm{}
		app.renderTemplate(w, "home.tmpl", http.StatusOK, data)
	})
}

func (app *application) excerptCreatePost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		id, err := app.excerpts.Insert(form.Title)
		if err != nil {
			app.serverError(w, err)
			return
		}

		idStr := strings.ReplaceAll(id.String(), "-", "")
		http.Redirect(w, r, fmt.Sprintf("/excerpt?id=%v", idStr), http.StatusSeeOther)
	})
}
