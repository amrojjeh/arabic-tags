package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/amrojjeh/arabic-tags/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger   *slog.Logger
	page     map[string]*template.Template
	excerpts models.ExcerptModel
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP Address")
	dsn := flag.String("dsn", "web:pass@/arabic_tags?parseTime=true",
		"Data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	db, err := sql.Open("mysql", *dsn)
	if err != nil {
		logger.Error("cannot open db", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Error("cannot open connection with db", slog.String("error",
			err.Error()))
		os.Exit(1)
	}

	app := application{
		logger:   logger,
		excerpts: models.ExcerptModel{DB: db},
	}

	err = app.cacheTemplates()
	if err != nil {
		logger.Error("cannot cache templates", slog.String("error", err.Error()))
		os.Exit(1)
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
