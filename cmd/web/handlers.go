package main

import (
	"errors"
	"net/http"

	"github.com/amrojjeh/arabic-tags/internal/models"
	"github.com/amrojjeh/arabic-tags/internal/validator"
	"github.com/amrojjeh/arabic-tags/ui/layers"
	"github.com/amrojjeh/arabic-tags/ui/pages"
	"github.com/amrojjeh/arabic-tags/ui/partials"
	"github.com/amrojjeh/kalam"
)

func (app *application) notFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}

func (app *application) registerGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := pages.RegisterPage(pages.RegisterProps{
			LoginUrl:    app.u.login(),
			RegisterUrl: app.u.register(),
			LogoutUrl:   app.u.logout(),
		}).Render(w)
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

		props := res.Props(app.u.login(), app.u.register(), app.u.logout())
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

		http.Redirect(w, r, app.u.login(), http.StatusSeeOther)
	})
}

func (app *application) loginGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := pages.LoginPage(pages.LoginProps{
			LoginUrl:    app.u.login(),
			RegisterUrl: app.u.register(),
			LogoutUrl:   app.u.logout(),
		}).Render(w)
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

		props := res.Props(app.u.login(), app.u.register(), app.u.logout())
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

		http.Redirect(w, r, app.u.home(), http.StatusSeeOther)
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

		http.Redirect(w, r, app.u.login(), http.StatusSeeOther)
	})
}

func (app *application) homeGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := app.getAuthenticatedEmail(r)
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
				Url:  app.u.excerpt(excerpts[x].Id),
			}
		}
		err = pages.HomePage(pages.HomeProps{
			Username:    user.Username,
			Excerpts:    homeExcerpts,
			Error:       app.session.PopString(r.Context(), errorSessionKey),
			AddUrl:      app.u.createExcerpt(),
			LoginUrl:    app.u.login(),
			RegisterUrl: app.u.register(),
			LogoutUrl:   app.u.logout(),
		}).Render(w)
		if err != nil {
			app.serverError(w, err)
		}
	})
}

func (app *application) createExcerptGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := layers.ExcerptLayer(app.u.createExcerpt()).Render(w)
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
		valid := true
		msg := validator.NewValidator("title", res.Title).
			Required().
			MaxLength(100).
			Validate(&valid)
		if !valid {
			app.session.Put(r.Context(), errorSessionKey, msg)
			http.Redirect(w, r, app.u.home(), http.StatusSeeOther)
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

		http.Redirect(w, r, app.u.excerpt(id), http.StatusSeeOther)
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

		props := pages.ManuscriptProps{
			ExcerptTitle:        excerpt.Title,
			ReadOnly:            false,
			AcceptedPunctuation: kalam.PunctuationRegex().String(),
			Content:             manuscript.Content,
			Error:               app.session.PopString(r.Context(), errorSessionKey),
			SubmitUrl:           r.URL.String(),
			TitleUrl:            app.u.excerptTitle(excerpt.Id),
			LoginUrl:            app.u.login(),
			RegisterUrl:         app.u.register(),
			LogoutUrl:           app.u.logout(),
		}

		email := app.getAuthenticatedEmail(r)
		if loggedIn := email != ""; !loggedIn {
			props.ReadOnly = true
			props.Warning = "Log in as the owner if you wish to edit the excerpt"
		} else {
			user, err := app.user.Get(email)
			if err != nil {
				app.serverError(w, err)
				return
			}
			props.Username = user.Username

			if owner := email == excerpt.AuthorEmail; !owner {
				props.ReadOnly = true
				props.Warning = "You cannot make changes as you're not the owner of the excerpt"
			}
		}

		err = pages.ManuscriptPage(props).Render(w)

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

func (app *application) excerptTitleGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e := getExcerptFromContext(r.Context())
		err := partials.TitleForm(app.u.excerptTitle(e.Id), e.Title).Render(w)
		if err != nil {
			app.serverError(w, err)
		}
	})
}

func (app *application) excerptTitlePost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.serverError(w, err)
			return
		}

		excerpt := getExcerptFromContext(r.Context())
		title := r.Form.Get("title")
		valid := true
		msg := validator.NewValidator("title", title).
			Required().
			MaxLength(100).
			Validate(&valid)
		if !valid {
			err = partials.WithError(msg, partials.TitleRegular(
				app.u.excerptTitle(excerpt.Id), excerpt.Title)).Render(w)
			if err != nil {
				app.serverError(w, err)
			}
			return
		}
		err = app.excerpt.UpdateTitle(excerpt.Id, title)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = partials.TitleRegular(app.u.excerptTitle(excerpt.Id),
			title).Render(w)
		if err != nil {
			app.serverError(w, err)
		}
	})
}
