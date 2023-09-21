package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type application struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP Address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := application{
		logger: logger,
	}
	server := &http.Server{
		Handler:      app.routes(),
		Addr:         *addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("starting server", slog.String("addr", *addr))
	if err := server.ListenAndServe(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
