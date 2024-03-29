package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/amrojjeh/arabic-tags/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger     *slog.Logger
	u          url
	page       map[string]*template.Template
	user       models.UserModel
	excerpt    models.ExcerptModel
	word       models.WordModel
	manuscript models.ManuscriptModel
	session    *scs.SessionManager
	mutex      sync.Mutex
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP Address")
	dsn := flag.String("dsn", "web:pass@/arabic_tags?parseTime=true",
		"Data source name")
	cert := flag.String("cert", "./tls/cert.pem", "location of tls certificate")
	key := flag.String("key", "./tls/key.pem", "location of tls private key")
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

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Store = mysqlstore.New(db)

	app := application{
		logger:     logger,
		user:       models.UserModel{Db: db},
		excerpt:    models.ExcerptModel{Db: db},
		manuscript: models.ManuscriptModel{Db: db},
		word:       models.WordModel{Db: db},
		session:    session,
	}

	if err != nil {
		logger.Error("cannot cache templates", slog.String("error", err.Error()))
		os.Exit(1)
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
	}
	server := &http.Server{
		Handler:      app.routes(),
		Addr:         *addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		TLSConfig:    tlsConfig,
	}

	logger.Info("starting server", slog.String("addr", *addr))
	err = server.ListenAndServeTLS(*cert, *key)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
