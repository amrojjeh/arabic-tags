package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/amrojjeh/arabic-tags/internal/models"
	"github.com/amrojjeh/arabic-tags/internal/validator"
	"github.com/amrojjeh/arabic-tags/ui/layers"
	"github.com/amrojjeh/arabic-tags/ui/pages"
	"github.com/julienschmidt/httprouter"
)

func (app *application) notFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}

func (app *application) registerGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := pages.RegisterPage(pages.RegisterProps{}).Render(w)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
		}
	})
}

func (app *application) registerPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := pages.NewRegisterResponse(r)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		props := res.Props()
		valid := true
		props.EmailError = validator.NewValidator("email", res.Email).
			Required().
			IsEmail().
			MaxLength(255).
			Validate(&valid)

		props.UsernameError = validator.NewValidator("username", res.Username).
			Required().
			MaxLength(255).
			Validate(&valid)

		props.PasswordError =
			validator.NewValidator("password", res.Password).
				Required().
				SameAs(res.RePassword).
				MaxBytes(72).
				Validate(&valid)

		if !valid {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err = pages.RegisterPage(props).Render(w)
			if err != nil {
				app.serverError(w, err)
			}
			return
		}

		err = app.user.Register(res.Username, res.Email, res.Password)
		if err != nil {
			if errors.Is(err, models.ErrDuplicateEmail) {
				props.EmailError = "Email is already taken"
				err = pages.RegisterPage(props).Render(w)
				if err != nil {
					app.serverError(w, err)
				}
				return
			} else if errors.Is(err, models.ErrDuplicateUsername) {
				props.UsernameError = "Username is already taken"
				err = pages.RegisterPage(props).Render(w)
				if err != nil {
					app.serverError(w, err)
				}
				return
			}
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}

func (app *application) loginGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := pages.LoginPage(pages.LoginProps{}).Render(w)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
		}
	})
}

func (app *application) loginPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := pages.NewLoginResponse(r)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		props := res.Props()
		valid := true
		props.EmailError = validator.NewValidator("email", res.Email).
			Required().
			IsEmail().
			Validate(&valid)
		props.PasswordError = validator.NewValidator("password", res.Password).
			Required().
			Validate(&valid)

		if !valid {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err = pages.LoginPage(props).Render(w)
			if err != nil {
				app.serverError(w, err)
			}
			return
		}

		auth, err := app.user.Authenticate(res.Email, res.Password)
		if err != nil {
			app.serverError(w, err)
			return
		}

		if !auth {
			props.LoginError = "The email or password is incorrect"
			err = pages.LoginPage(props).Render(w)
			if err != nil {
				app.serverError(w, err)
			}
			return
		}

		app.session.Put(r.Context(), authorizedEmailSessionKey, res.Email)
		err = app.session.RenewToken(r.Context())
		if err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	})
}

func (app *application) logoutPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.session.Remove(r.Context(), authorizedEmailSessionKey)
		err := app.session.RenewToken(r.Context())
		if err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}

func (app *application) homeGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := app.session.GetString(r.Context(), authorizedEmailSessionKey)
		user, err := app.user.Get(email)
		if err != nil {
			app.serverError(w, err)
			return
		}
		excerpts, err := app.excerpt.GetByEmail(email)
		if err != nil {
			app.serverError(w, err)
			return
		}

		homeExcerpts := make([]pages.HomeExcerpt, len(excerpts))
		for x := 0; x < len(homeExcerpts); x++ {
			homeExcerpts[x] = pages.HomeExcerpt{
				Name: excerpts[x].Title,
				Url:  fmt.Sprintf("/excerpt/%d", excerpts[x].Id),
			}
		}
		err = pages.HomePage(pages.HomeProps{
			Username: user.Username,
			Excerpts: homeExcerpts,
			AddUrl:   "/excerpt",
		}).Render(w)
		if err != nil {
			app.serverError(w, err)
		}
	})
}

func (app *application) createExcerptGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := layers.ExcerptLayer("/excerpt").Render(w)
		if err != nil {
			app.serverError(w, err)
		}
	})
}

func (app *application) createExcerptPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := layers.NewExcerptResponse(r)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		email := app.session.GetString(r.Context(), authorizedEmailSessionKey)
		id, err := app.excerpt.Insert(res.Title, email)
		if err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/excerpt/%v", id), http.StatusSeeOther)
	})
}

func (app *application) excerptGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(httprouter.ParamsFromContext(r.Context()).ByName("id"))
		if err != nil {
			app.clientError(w, http.StatusUnprocessableEntity)
			return
		}
		w.Write([]byte(strconv.Itoa(id)))
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
