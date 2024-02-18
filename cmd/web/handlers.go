package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/amrojjeh/arabic-tags/internal/models"
	"github.com/amrojjeh/arabic-tags/internal/validator"
	"github.com/amrojjeh/arabic-tags/ui/layers"
	"github.com/amrojjeh/arabic-tags/ui/pages"
	"github.com/amrojjeh/kalam"
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

		_, err = app.manuscript.Insert(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/excerpt/%v", id), http.StatusSeeOther)
	})
}

func (app *application) excerptGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		excerpt := getExcerptFromContext(r.Context())

		manuscript, err := app.manuscript.GetByExcerptId(excerpt.Id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound().ServeHTTP(w, r)
				return
			}
			app.serverError(w, err)
			return
		}

		err = pages.ManuscriptPage(pages.ManuscriptProps{
			ExcerptTitle:        excerpt.Title,
			ReadOnly:            false,
			AcceptedPunctuation: kalam.PunctuationRegex().String(),
			Content:             manuscript.Content,
			SubmitUrl:           r.URL.String(),
		}).Render(w)

		if err != nil {
			app.serverError(w, err)
		}
	})
}

func (app *application) excerptPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusUnprocessableEntity)
			return
		}

		e := getExcerptFromContext(r.Context())
		if err != nil {
			app.serverError(w, err)
			return
		}

		c := r.Form.Get("content")
		err = app.manuscript.UpdateByExcerptId(e.Id, c)
		if err != nil {
			app.serverError(w, err)
			return
		}
	})
}
