package main

import (
	"fmt"
	"net/http"

	"github.com/amrojjeh/arabic-tags/internal/validator"
	"github.com/amrojjeh/arabic-tags/ui/pages"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) notFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}

func (app *application) homeGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := pages.HomePage(pages.HomeProps{
			TitleField: "",
			TitleError: "",
		}).Render(w)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
		}
	})
}

func (app *application) homePost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		title := r.Form.Get("title")
		password := r.Form.Get("password")

		titleError :=
			validator.NewValidator("title", title).NotBlank().MaxLength(100).
				Validate()
		passwordError :=
			validator.NewValidator("password", password).NotBlank().MaxBytes(72).
				Validate()

		if titleError != "" || passwordError != "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err = pages.HomePage(pages.HomeProps{
				TitleField:    title,
				TitleError:    titleError,
				PasswordError: passwordError,
			}).Render(w)

			if err != nil {
				app.serverError(w, err)
			}

			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			app.serverError(w, err)
			return
		}

		id, err := app.excerpts.Insert(title, hashed)
		if err != nil {
			app.serverError(w, err)
			return
		}

		idStr := idToString(id)
		http.Redirect(w, r, fmt.Sprintf("/excerpt/manuscript?id=%v", idStr),
			http.StatusSeeOther)
	})
}

func (app *application) manuscriptGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(Amr Ojjeh): Check if logged into excerpt
		// If not, show readonly mode
		// Otherwise, show editable mode
	})
}

// func (app *application) excerptEditUnlock() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		id := r.Context().Value(idContextKey).(uuid.UUID)
// 		err := app.excerpts.SetContentLock(id, false)
// 		if err != nil {
// 			app.serverError(w, err)
// 			return
// 		}

// 		http.Redirect(w, r, fmt.Sprintf("/excerpt/edit?id=%v", idToString(id)), http.StatusSeeOther)
// 	})
// }

// func (app *application) excerptEditLock() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		id := r.Context().Value(idContextKey).(uuid.UUID)
// 		excerpt := r.Context().Value(excerptContextKey).(models.Excerpt)
// 		if !kalam.IsContentClean(excerpt.Content) {
// 			http.Redirect(w, r, fmt.Sprintf(
// 				"/excerpt/edit?id=%v&error=Could not proceed. Found errors.", id),
// 				http.StatusSeeOther)
// 			return
// 		}
// 		content := kalam.RemoveExtraWhitespace(excerpt.Content)
// 		if content == "" {
// 			http.Redirect(w, r, fmt.Sprintf(
// 				"/excerpt/edit?id=%v&error=Could not proceed. There's no text.", id),
// 				http.StatusSeeOther)
// 			return
// 		}
// 		err := app.excerpts.UpdateContent(id, content)
// 		if err != nil {
// 			app.serverError(w, err)
// 			return
// 		}

// 		err = app.excerpts.SetContentLock(id, true)
// 		if err != nil {
// 			app.serverError(w, err)
// 			return
// 		}

// 		err = app.excerpts.ResetGrammar(id)
// 		if err != nil {
// 			app.serverError(w, err)
// 			return
// 		}
// 		idStr := idToString(id)
// 		http.Redirect(w, r, fmt.Sprintf("/excerpt/grammar?id=%v", idStr),
// 			http.StatusSeeOther)
// 	})
// }

// func (app *application) excerptEditPut() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		id := r.Context().Value(idContextKey).(uuid.UUID)
// 		content := r.Form.Get("content")
// 		var err error
// 		err = app.excerpts.UpdateContent(id, content)
// 		if err != nil {
// 			app.serverError(w, err)
// 			return
// 		}
// 		app.noBody(w)
// 	})
// }

// func (app *application) excerptGrammarGet() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		data, err := newTemplateData(r)
// 		if err != nil {
// 			app.clientError(w, http.StatusBadRequest)
// 			return
// 		}
// 		app.renderTemplate(w, "grammar.tmpl", http.StatusOK, data)
// 	})
// }

// TODO(Amr Ojjeh): Verify tags
// func (app *application) excerptGrammarPut() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		content := r.Form.Get("content")
// 		var grammar models.Grammar
// 		err := json.Unmarshal([]byte(content), &grammar)
// 		if err != nil {
// 			app.clientError(w, http.StatusBadRequest)
// 			return
// 		}

// 		id := r.Context().Value(idContextKey).(uuid.UUID)
// 		err = app.excerpts.UpdateGrammar(id, grammar)
// 		if err != nil {
// 			app.serverError(w, err)
// 			return
// 		}
// 		app.noBody(w)
// 	})
// }

// func (app *application) excerptGrammarLock() http.Handler {
// 	// TODO(Amr Ojjeh): Verify tags
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		id := r.Context().Value(idContextKey).(uuid.UUID)
// 		err := app.excerpts.SetGrammarLock(id, true)
// 		if err != nil {
// 			app.serverError(w, err)
// 			return
// 		}

// 		err = app.excerpts.ResetTechnical(id)
// 		if err != nil {
// 			app.serverError(w, err)
// 		}

// 		idStr := idToString(id)
// 		http.Redirect(w, r, fmt.Sprintf("/excerpt/technical?id=%v", idStr),
// 			http.StatusSeeOther)
// 	})
// }

// func (app *application) excerptGrammarUnlock() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		id := r.Context().Value(idContextKey).(uuid.UUID)
// 		err := app.excerpts.SetGrammarLock(id, false)
// 		if err != nil {
// 			app.serverError(w, err)
// 			return
// 		}
// 		idStr := idToString(id)
// 		http.Redirect(w, r, fmt.Sprintf("/excerpt/grammar?id=%v", idStr),
// 			http.StatusSeeOther)
// 	})
// }

// func (app *application) excerptTechnicalGet() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		data, err := newTemplateData(r)
// 		if err != nil {
// 			app.clientError(w, http.StatusBadRequest)
// 			return
// 		}
// 		app.renderTemplate(w, "technical.tmpl", http.StatusOK, data)
// 	})
// }

// func (app *application) excerptTechnicalVowelPut() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		excerpt := r.Context().Value(excerptContextKey).(models.Excerpt)
// 		wordIndex := r.Context().Value(wordIndexContextKey).(int)
// 		letterIndex, err := strconv.Atoi(r.Form.Get("letter"))
// 		if err != nil {
// 			app.clientError(w, http.StatusBadRequest)
// 			return
// 		}

// 		vowel := r.Form.Get(strconv.Itoa(letterIndex))

// 		if !Radio(vowel, []string{
// 			string(kalam.Damma),
// 			string(kalam.Dammatan),
// 			string(kalam.Kasra),
// 			string(kalam.Kasratan),
// 			string(kalam.Fatha),
// 			string(kalam.Fathatan),
// 			string(kalam.Sukoon),
// 		}) {
// 			app.clientError(w, http.StatusBadRequest)
// 			return
// 		}

// 		excerpt.Technical.Words[wordIndex].Letters[letterIndex].Vowel = vowel
// 		err = app.excerpts.UpdateTechnical(excerpt.ID, excerpt.Technical)
// 		if err != nil {
// 			app.serverError(w, err)
// 			return
// 		}
// 		data, err := newTemplateData(r)
// 		if err != nil {
// 			app.clientError(w, http.StatusBadRequest)
// 			return
// 		}
// 		data.TSelectedWord = wordIndex
// 		app.renderTemplate(w, "htmx-technical-inspector-update.tmpl", http.StatusOK, data)
// 	})
// }

// func (app *application) excerptTechnicalShadda() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		e := r.Context().Value(excerptContextKey).(models.Excerpt)
// 		wi := r.Context().Value(wordIndexContextKey).(int)
// 		letterIndex, err := strconv.Atoi(r.Form.Get("letter"))
// 		if err != nil {
// 			app.clientError(w, http.StatusBadRequest)
// 			return
// 		}
// 		e.Technical.Words[wi].Letters[letterIndex].Shadda =
// 			r.Form.Get("shadda") == "true"
// 		err = app.excerpts.UpdateTechnical(e.ID, e.Technical)
// 		if err != nil {
// 			app.serverError(w, err)
// 			return
// 		}
// 		data, err := newTemplateData(r)
// 		if err != nil {
// 			app.clientError(w, http.StatusBadRequest)
// 			return
// 		}
// 		data.TSelectedWord = wi
// 		app.renderTemplate(w, "htmx-technical-inspector-update.tmpl", http.StatusOK, data)
// 	})
// }

// func (app *application) excerptTechnicalWordGet() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		index := r.Context().Value(wordIndexContextKey).(int)
// 		data, err := newTemplateData(r)
// 		if err != nil {
// 			app.clientError(w, http.StatusBadRequest)
// 			return
// 		}

// 		data.TSelectedWord = index
// 		app.renderTemplate(w, "htmx-technical-change-word.tmpl", http.StatusOK, data)
// 	})
// }

// func (app *application) excerptTechnicalSentenceStart() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		e := r.Context().Value(excerptContextKey).(models.Excerpt)
// 		wi := r.Context().Value(wordIndexContextKey).(int)
// 		e.Technical.Words[wi].SentenceStart =
// 			r.Form.Get("sentenceStart") == "true"
// 		err := app.excerpts.UpdateTechnical(e.ID, e.Technical)
// 		if err != nil {
// 			app.serverError(w, err)
// 			return
// 		}
// 		app.noBody(w)
// 	})
// }

// func (app *application) excerptTechnicalIgnore() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		e := r.Context().Value(excerptContextKey).(models.Excerpt)
// 		wi := r.Context().Value(wordIndexContextKey).(int)
// 		e.Technical.Words[wi].Ignore =
// 			r.Form.Get("ignore") == "true"
// 		err := app.excerpts.UpdateTechnical(e.ID, e.Technical)
// 		if err != nil {
// 			app.serverError(w, err)
// 			return
// 		}
// 		app.noBody(w)
// 	})
// }

// func (app *application) excerptTechnicalExport() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		e := r.Context().Value(excerptContextKey).(models.Excerpt)
// 		buff, err := e.Export()
// 		if err != nil {
// 			app.serverError(w, err)
// 			return
// 		}
// 		w.Write(buff)
// 	})
// }
