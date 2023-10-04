package main

import (
	"log/slog"
	"net/http"
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

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}
