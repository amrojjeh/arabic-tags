package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/amrojjeh/arabic-tags/internal/models"
	"github.com/google/uuid"
)

func (app *application) notFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}

func (app *application) excerptEditGet() http.Handler {
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

func (app *application) excerptEditLock() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(Amr Ojjeh): Prevent next if there are errors
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
		err = app.excerpts.SetContentLock(id, true)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// TODO(Amr Ojjeh): Create Grammar
		err = app.excerpts.ResetGrammar(id)
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/excerpt/grammar?id=%v", idStr), http.StatusSeeOther)
	})
}

func (app *application) excerptEditPut() http.Handler {
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
		err = app.excerpts.UpdateContent(id, content)
		if err != nil {
			app.serverError(w, err)
			return
		}
		app.noBody(w)
	})
}

func (app *application) excerptGrammarGet() http.Handler {
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

		excerpt, err := app.excerpts.Get(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		data := templateData{
			Excerpt: excerpt,
			Type:    "grammar",
		}
		app.renderTemplate(w, "grammar.tmpl", http.StatusOK, data)
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
		http.Redirect(w, r, fmt.Sprintf("/excerpt/edit?id=%v", idStr), http.StatusSeeOther)
	})
}
