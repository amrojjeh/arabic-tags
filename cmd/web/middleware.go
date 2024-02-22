package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/amrojjeh/arabic-tags/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *application) recoverPanic(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		h.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("request made",
			slog.String("url", r.URL.String()),
			slog.String("method", r.Method),
			slog.String("remote-addr", r.RemoteAddr))
		h.ServeHTTP(w, r)
	})
}

func (app *application) authRequired(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := app.session.GetString(r.Context(), authorizedEmailSessionKey)
		if email == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (app *application) getUser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.session.GetString(r.Context(), authorizedEmailSessionKey) != "" {
			user, err := app.user.Get(app.session.GetString(
				r.Context(),
				authorizedEmailSessionKey))
			if err != nil {
				app.serverError(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			r = r.WithContext(ctx)
		}

		h.ServeHTTP(w, r)
	})
}

func (app *application) excerptRequired(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := httprouter.ParamsFromContext(r.Context()).ByName("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.notFound().ServeHTTP(w, r)
			return
		}

		excerpt, err := app.excerpt.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound().ServeHTTP(w, r)
				return
			}
			app.serverError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), excerptContextKey, excerpt)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

func getExcerptFromContext(c context.Context) models.Excerpt {
	return c.Value(excerptContextKey).(models.Excerpt)
}

func getUserFromContext(c context.Context) models.User {
	return c.Value(userContextKey).(models.User)
}
