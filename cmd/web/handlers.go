package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode"

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
		id := r.Context().Value("id").(uuid.UUID)
		excerpt, err := app.excerpts.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.excerptNotFound(w, r)
				return
			}
			app.serverError(w, err)
			return
		}
		data := newTemplateData(r)
		data.Excerpt = excerpt
		app.renderTemplate(w, "add.tmpl", http.StatusOK, data)
	})
}

func (app *application) excerptEditUnlock() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value("id").(uuid.UUID)
		err := app.excerpts.SetContentLock(id, false)
		if err != nil {
			app.serverError(w, err)
			return
		}

	})
}

func cleanContent(content string) (string, error) {
	for _, c := range content {
		if !(isArabicLetter(c) || isWhitespace(c)) {
			return "", errors.New(fmt.Sprintf("%v is an invalid letter", c))
		}
	}

	// Remove double spaces
	r, _ := regexp.Compile(" +")
	content = r.ReplaceAllString(content, " ")

	// Trim sentence
	content = strings.TrimFunc(content, unicode.IsSpace)
	return content, nil
}

// isArabicLetter does not include tashkeel
func isArabicLetter(letter rune) bool {
	if letter >= 0x0621 && letter <= 0x063A {
		return true
	}
	if letter >= 0x0641 && letter <= 0x064A {
		return true
	}
	return false
}

func isWhitespace(letter rune) bool {
	return letter == ' '
}

func (app *application) excerptEditLock() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value("id").(uuid.UUID)
		excerpt, err := app.excerpts.Get(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		content, err := cleanContent(excerpt.Content)
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf(
				"/excerpt/edit?id=%v&error=Could not proceed. Found errors.", id),
				http.StatusSeeOther)
			return
		}

		err = app.excerpts.UpdateContent(id, content)
		if err != nil {
			app.serverError(w, err)
			return
		}

		err = app.excerpts.SetContentLock(id, true)
		if err != nil {
			app.serverError(w, err)
			return
		}

		err = app.excerpts.ResetGrammar(id)
		if err != nil {
			app.serverError(w, err)
			return
		}
		idStr := idToString(id)
		http.Redirect(w, r, fmt.Sprintf("/excerpt/grammar?id=%v", idStr),
			http.StatusSeeOther)
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
		id := r.Context().Value("id").(uuid.UUID)
		excerpt, err := app.excerpts.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.excerptNotFound(w, r)
				return
			}
			app.serverError(w, err)
			return
		}

		data := newTemplateData(r)
		data.Type = "grammar"
		data.Excerpt = excerpt
		app.renderTemplate(w, "grammar.tmpl", http.StatusOK, data)
	})
}

// TODO(Amr Ojjeh): Verify tags
func (app *application) excerptGrammarPut() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		content := r.Form.Get("content")
		var grammar models.Grammar
		err = json.Unmarshal([]byte(content), &grammar)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		id := r.Context().Value("id").(uuid.UUID)
		err = app.excerpts.UpdateGrammar(id, grammar)
		if err != nil {
			app.serverError(w, err)
			return
		}
		app.noBody(w)
	})
}

type excerptForm struct {
	Validator
	Title string
}

func (app *application) excerptCreateGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := newTemplateData(r)
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
			data := newTemplateData(r)
			data.Form = form

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

		idStr := idToString(id)
		http.Redirect(w, r, fmt.Sprintf("/excerpt/edit?id=%v", idStr),
			http.StatusSeeOther)
	})
}
