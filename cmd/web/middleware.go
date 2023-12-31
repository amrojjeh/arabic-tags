package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/amrojjeh/arabic-tags/internal/models"
	"github.com/google/uuid"
)

type Adapter func(http.Handler) http.Handler

func (app *application) logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("request made", slog.String("url", r.URL.String()),
			slog.String("proto", r.Proto),
			slog.String("method", r.Method),
			slog.String("remote-addr", r.RemoteAddr))
		h.ServeHTTP(w, r)
	})
}

func stripPrefix(prefix string) Adapter {
	return func(h http.Handler) http.Handler {
		return http.StripPrefix(prefix, h)
	}
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

func (app *application) excerptRequired(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(idContextKey).(uuid.UUID)
		shared := r.Form.Get("share") == "true"
		var excerpt models.Excerpt
		var err error
		if shared {
			if strings.Contains(r.URL.Path, "edit") {
				excerpt, err = app.excerpts.GetSharedContent(id)
			} else if strings.Contains(r.URL.Path, "grammar") {
				excerpt, err = app.excerpts.GetSharedGrammar(id)
			} else if strings.Contains(r.URL.Path, "technical") {
				excerpt, err = app.excerpts.GetSharedTechnical(id)
			} else {
				app.clientError(w, http.StatusNotFound)
				return
			}
		} else {
			excerpt, err = app.excerpts.Get(id)
		}
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.excerptNotFound(w, r)
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

func (app *application) contentLockRequired(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		excerpt := r.Context().Value(excerptContextKey).(models.Excerpt)
		if !excerpt.CLocked {
			id := r.Context().Value(idContextKey).(uuid.UUID)
			http.Redirect(w, r,
				fmt.Sprintf("/excerpt/edit?id=%v", idToString(id)),
				http.StatusSeeOther)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (app *application) grammarLockRequired(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		excerpt := r.Context().Value(excerptContextKey).(models.Excerpt)
		if !excerpt.GLocked {
			id := r.Context().Value(idContextKey).(uuid.UUID)
			http.Redirect(w, r,
				fmt.Sprintf("/excerpt/grammar?id=%v", idToString(id)),
				http.StatusSeeOther)
			return
		}
		h.ServeHTTP(w, r)
	})
}

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

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}
