package main

import (
	"log"
	"net/http"
	"time"

	"github.com/amrojjeh/arabic-tags/routes"
	"github.com/gorilla/mux"
)

func setupRoutes(r *mux.Router) {
	routes.Index(r)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("assets"))))
}

func main() {
	r := mux.NewRouter()
	setupRoutes(r)

	addr := "localhost:8080"

	server := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Print("Hosting at ", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
