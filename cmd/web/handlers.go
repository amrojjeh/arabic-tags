package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/amrojjeh/arabic-tags/internal/models"
	"github.com/amrojjeh/arabic-tags/internal/speech"
	"github.com/google/uuid"
)

func (app *application) notFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}

func (app *application) excerptEditGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := newTemplateData(r)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		app.renderTemplate(w, "add.tmpl", http.StatusOK, data)
	})
}

func (app *application) excerptEditUnlock() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(idContextKey).(uuid.UUID)
		err := app.excerpts.SetContentLock(id, false)
		if err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/excerpt/edit?id=%v", idToString(id)), http.StatusSeeOther)
	})
}

func (app *application) excerptEditLock() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(idContextKey).(uuid.UUID)
		excerpt := r.Context().Value(excerptContextKey).(models.Excerpt)
		content, err := speech.CleanContent(excerpt.Content)
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf(
				"/excerpt/edit?id=%v&error=Could not proceed. Found errors.", id),
				http.StatusSeeOther)
			return
		}
		if content == "" {
			http.Redirect(w, r, fmt.Sprintf(
				"/excerpt/edit?id=%v&error=Could not proceed. There's no text.", id),
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
		id := r.Context().Value(idContextKey).(uuid.UUID)
		content := r.Form.Get("content")
		shared := r.Form.Get("share") == "true"
		var err error
		if shared {
			err = app.excerpts.UpdateSharedContent(id, content)
		} else {
			err = app.excerpts.UpdateContent(id, content)
		}
		if err != nil {
			app.serverError(w, err)
			return
		}
		app.noBody(w)
	})
}

func (app *application) excerptGrammarGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := newTemplateData(r)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		app.renderTemplate(w, "grammar.tmpl", http.StatusOK, data)
	})
}

// TODO(Amr Ojjeh): Verify tags
func (app *application) excerptGrammarPut() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content := r.Form.Get("content")
		shared := r.Form.Get("share") == "true"
		var grammar models.Grammar
		err := json.Unmarshal([]byte(content), &grammar)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		id := r.Context().Value(idContextKey).(uuid.UUID)
		if shared {
			err = app.excerpts.UpdateSharedGrammar(id, grammar)
		} else {
			err = app.excerpts.UpdateGrammar(id, grammar)
		}
		if err != nil {
			app.serverError(w, err)
			return
		}
		app.noBody(w)
	})
}

func (app *application) excerptGrammarLock() http.Handler {
	// TODO(Amr Ojjeh): Verify tags
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(idContextKey).(uuid.UUID)
		err := app.excerpts.SetGrammarLock(id, true)
		if err != nil {
			app.serverError(w, err)
			return
		}

		err = app.excerpts.ResetTechnical(id)
		if err != nil {
			app.serverError(w, err)
		}

		idStr := idToString(id)
		http.Redirect(w, r, fmt.Sprintf("/excerpt/technical?id=%v", idStr),
			http.StatusSeeOther)
	})
}

func (app *application) excerptGrammarUnlock() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(idContextKey).(uuid.UUID)
		err := app.excerpts.SetGrammarLock(id, false)
		if err != nil {
			app.serverError(w, err)
			return
		}
		idStr := idToString(id)
		http.Redirect(w, r, fmt.Sprintf("/excerpt/grammar?id=%v", idStr),
			http.StatusSeeOther)
	})
}

func (app *application) excerptTechnicalGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := newTemplateData(r)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		app.renderTemplate(w, "technical.tmpl", http.StatusOK, data)
	})
}

func (app *application) excerptTechnicalVowelPut() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		excerpt := r.Context().Value(excerptContextKey).(models.Excerpt)
		wordIndex := r.Context().Value(wordIndexContextKey).(int)
		letterIndex, err := strconv.Atoi(r.Form.Get("letter"))
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		vowel := r.Form.Get(strconv.Itoa(letterIndex))

		if !Radio(vowel, []string{
			speech.Damma,
			speech.Dammatan,
			speech.Kasra,
			speech.Kasratan,
			speech.Fatha,
			speech.Fathatan,
			speech.Sukoon,
		}) {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		excerpt.Technical.Words[wordIndex].Letters[letterIndex].Vowel = vowel
		err = app.excerpts.UpdateTechnical(excerpt.ID, excerpt.Technical)
		if err != nil {
			app.serverError(w, err)
			return
		}
		data, err := newTemplateData(r)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		data.TSelectedWord = wordIndex
		app.renderTemplate(w, "htmx-technical.tmpl", http.StatusOK, data)
	})
}

func (app *application) excerptTechnicalWordGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		index := r.Context().Value(wordIndexContextKey).(int)
		data, err := newTemplateData(r)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		data.TSelectedWord = index
		app.renderTemplate(w, "htmx-technical.tmpl", http.StatusOK, data)
	})
}

func (app *application) excerptTechnicalSentenceStart() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(idContextKey).(uuid.UUID)
		e := r.Context().Value(excerptContextKey).(models.Excerpt)
		wi := r.Context().Value(wordIndexContextKey).(int)
		e.Technical.Words[wi].SentenceStart =
			r.Form.Get("sentenceStart") == "true"
		data, err := newTemplateData(r)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		err = app.excerpts.UpdateTechnical(id, e.Technical)
		if err != nil {
			app.serverError(w, err)
			return
		}
		app.renderTemplate(w, "htmx-technical.tmpl", http.StatusOK, data)
	})
}

type excerptForm struct {
	Validator
	Title string
}

func (app *application) excerptCreateGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := newTemplateData(r)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
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
			data, err := newTemplateData(r)
			if err != nil {
				app.clientError(w, http.StatusBadRequest)
				return
			}
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
