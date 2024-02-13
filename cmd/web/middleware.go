package main

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (app *application) logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("request made", slog.String("url", r.URL.String()),
			slog.String("proto", r.Proto),
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

func (app *application) idRequired(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		idStr := r.Form.Get("id")
		if idStr == "" {
			http.Redirect(w, r, "/?error=Excerpt id was not provided",
				http.StatusSeeOther)
			return
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Redirect(w, r, "/?error=id is invalid",
				http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), idContextKey, id)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

// func (app *application) excerptRequired(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		id := r.Context().Value(idContextKey).(uuid.UUID)
// 		var excerpt models.Excerpt
// 		var err error
// 		excerpt, err = app.excerpt.Get(id)
// 		if err != nil {
// 			if errors.Is(err, models.ErrNoRecord) {
// 				app.excerptNotFound(w, r)
// 				return
// 			}
// 			app.serverError(w, err)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), excerptContextKey, excerpt)
// 		r = r.WithContext(ctx)

// 		h.ServeHTTP(w, r)
// 	})
// }

func (app *application) technicalWordRequired(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		index_str := r.Form.Get("word")
		index, err := strconv.Atoi(index_str)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), wordIndexContextKey, index)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
