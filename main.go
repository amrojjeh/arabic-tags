package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/amrojjeh/arabic-tags/speech"
	"github.com/gorilla/mux"
)

type app struct {
	templates *template.Template
	data      []speech.Sentence
}

func (a app) index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(Amr Ojjeh): If there aren't any sentences, change the current message to be a button to add a new sentence
		a.templates.ExecuteTemplate(w, "index.html", a.data)
	})
}

func (a app) newSentence() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.templates.ExecuteTemplate(w, "new_sentence.html", nil)
	})
}

func main() {
	r := mux.NewRouter()
	s := http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	defer s.Close()

	f, err := os.OpenFile("data.json", os.O_RDONLY|os.O_CREATE, 0o755)
	if err != nil {
		log.Fatal("Could not read or create file:", err)
	}
	data := make([]speech.Sentence, 0, 50)
	dec := json.NewDecoder(f)
	for dec.More() {
		var sen speech.Sentence
		dec.Decode(&sen)
		data = append(data, sen)
	}
	f.Close()

	// TODO(Amr Ojjeh): Add an option to connect to the db (requires auth)
	a := app{
		templates: template.Must(template.ParseGlob("templates/*.html")),
		data:      data}

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./assets")))).
		Methods("GET")

	r.Handle("/sentence/new", a.newSentence()).
		Methods("GET")

	r.Handle("/", a.index()).
		Methods("GET")

	// TODO(Amr Ojjeh): Add analyze option
	r.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("analyze.html")
		if err != nil {
			panic(err)
		}
		w.Write(data)
	})

	log.Println("Listening on http://127.0.0.1:8080")
	log.Fatal(s.ListenAndServe())
}
