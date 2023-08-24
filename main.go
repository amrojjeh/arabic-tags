package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/amrojjeh/arabic-tags/speech"
	"github.com/amrojjeh/goarabic"
	"github.com/gorilla/mux"
)

type app struct {
	templates *template.Template
	paragraph speech.Paragraph
	fileName  string
}

func (a *app) load() error {
	a.paragraph = speech.Paragraph{}
	b, err := os.ReadFile(a.fileName)
	if errors.Is(err, os.ErrNotExist) {
		log.Println("File", a.fileName, "does not exist. Starting from zero...")
		return nil
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &a.paragraph)
	if err != nil {
		return err
	}
	log.Println("Loaded paragraph from:", a.fileName)
	return nil
}

func (a *app) save() error {
	b, err := json.Marshal(a.paragraph)
	if err != nil {
		return err
	}
	err = os.WriteFile(a.fileName, b, 0o666)
	if err != nil {
		return err
	}

	log.Println("Paragraph was saved to:", a.fileName)
	return nil
}

func (a *app) index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(Amr Ojjeh): If there aren't any sentences, change the current message to be a button to add a new sentence
		a.templates.ExecuteTemplate(w, "index.html", a.paragraph)
	})
}

func (a *app) newSentence() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ar, _ := goarabic.SafeBWToAr("hVA mvAl")
		sen := speech.NewSentence(ar)
		a.paragraph.AddSentence(sen)
		a.save()
		log.Println("New Sentence:", sen)
		a.templates.ExecuteTemplate(w, "sentence-outer", sen)
	})
}

func main() {
	log.SetPrefix("main:")
	r := mux.NewRouter()
	s := http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	defer s.Close()

	// TODO(Amr Ojjeh): Add an option to connect to the db (requires auth)
	a := app{
		fileName:  "data.json",
		templates: template.Must(template.ParseGlob("templates/*.html"))}
	err := a.load()

	if err != nil {
		log.Fatal("Could not load: ", a.fileName, ": ", err)
	}

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./assets")))).
		Methods("GET")

	r.Handle("/sentence", a.newSentence()).
		Methods("POST")

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
