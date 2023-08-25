package main

// TODO(Amr Ojjeh): Add documentation
// TODO(Amr Ojjeh): Improve logging

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/amrojjeh/arabic-tags/speech"
	"github.com/amrojjeh/goarabic"
	"github.com/gorilla/mux"
)

func getVar(w http.ResponseWriter, r *http.Request, v string) string {
	vars := mux.Vars(r)
	str, ok := vars[v]
	if !ok {
		http.Error(w, fmt.Sprintf("%v was not passed into url", v), http.StatusInternalServerError)
		log.Fatalf("getSentence - %v was not passed into url: %v", v, r.URL)
	}
	return str
}

func getId(w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.ParseInt(getVar(w, r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "id has to be a positive integer", http.StatusBadRequest)
		return 0, err
	}
	return int(id), nil
}

type app struct {
	templates *template.Template
	paragraph speech.Paragraph
	fileName  string
	mux       *mux.Router
}

func (a *app) load() error {
	a.paragraph = speech.Paragraph{}
	b, err := os.ReadFile(a.fileName)
	if errors.Is(err, os.ErrNotExist) {
		log.Println("file", a.fileName, "does not exist. Starting from zero...")
		return nil
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &a.paragraph)
	if err != nil {
		return err
	}
	log.Println("loaded paragraph from:", a.fileName)
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

	log.Println("paragraph was saved to:", a.fileName)
	return nil
}

func (a *app) index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.templates.ExecuteTemplate(w, "index.html", a.paragraph)
	})
}

func (a *app) sentences() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.templates.ExecuteTemplate(w, "sentences.tmpl", a.paragraph)
	})
}

func (a *app) newSentence() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ar, _ := goarabic.SafeBWToAr("hVA mvAl")
		ar += " " + strconv.Itoa(len(a.paragraph.Sentences))
		sen := a.paragraph.AddSentence(ar)
		a.save()
		log.Println("new Sentence:", sen)
		log.Println("new sentence id:", sen.Id)
		a.templates.ExecuteTemplate(w, "sentence-outer.tmpl", sen)
	})
}

func (a *app) getSentence() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := getId(w, r)
		if err != nil {
			http.Error(w, "id has to be a positive integer", http.StatusBadRequest)
			return
		}
		sen, err := a.paragraph.GetSentenceId(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("sentence with id %v does not exist", id), http.StatusBadRequest)
			return
		}
		a.templates.ExecuteTemplate(w, "sentence-outer.tmpl", sen)
	})
}

func (a *app) loadInspector() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := getId(w, r)
		if err != nil {
			return
		}
		log.Println("load inspector for sentence id:", id)
		sen, err := a.paragraph.GetSentenceId(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("sentence with id %v does not exist", id), http.StatusBadRequest)
			return
		}
		a.templates.ExecuteTemplate(w, "inspector.tmpl", sen)
	})
}

func (a *app) deleteSentence() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := getId(w, r)
		if err != nil {
			return
		}
		a.paragraph.DeleteSentence(int(id))
		a.save()
		log.Println("deleted sentence id:", id)
		a.templates.ExecuteTemplate(w, "inspector.tmpl", nil)
	})
}

func (a *app) moveSentenceUp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := getId(w, r)
		if err != nil {
			return
		}
		// TODO(Amr Ojjeh): Return BadRequest if sentence does not exist
		a.paragraph.MoveSentenceUp(id)
		a.save()
		a.templates.ExecuteTemplate(w, "sentences.tmpl", a.paragraph)
	})
}

func (a *app) moveSentenceDown() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := getId(w, r)
		if err != nil {
			return
		}
		// TODO(Amr Ojjeh): Return BadRequest if sentence does not exist
		a.paragraph.MoveSentenceDown(id)
		a.save()
		a.templates.ExecuteTemplate(w, "sentences.tmpl", a.paragraph)
	})
}

func (a *app) sentenceEditView() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := getId(w, r)
		if err != nil {
			return
		}
		sen, err := a.paragraph.GetSentenceId(id)
		// TODO(Amr Ojjeh): Improve helper functions
		if err != nil {
			http.Error(w, fmt.Sprintf("sentence with id %v does not exist", id), http.StatusBadRequest)
			return
		}
		a.templates.ExecuteTemplate(w, "sentence-edit.tmpl", sen)
	})
}

func (a *app) sentenceEdit() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := getId(w, r)
		if err != nil {
			return
		}
		sen, err := a.paragraph.GetSentenceId(id)
		// TODO(Amr Ojjeh): Improve helper functions
		if err != nil {
			http.Error(w, fmt.Sprintf("sentence with id %v does not exist", id), http.StatusBadRequest)
			return
		}
		r.ParseForm()
		if !r.Form.Has("v") {
			http.Error(w, "v parameter is required", http.StatusBadRequest)
			return
		}
		value := r.Form.Get("v")
		a.paragraph.EditSentence(id, value)
		a.save()
		a.templates.ExecuteTemplate(w, "sentence-outer.tmpl", sen)
		a.templates.ExecuteTemplate(w, "inspector.tmpl", sen)
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
		mux:       r,
		fileName:  "data.json",
		templates: template.Must(template.ParseGlob("templates/*.html"))}
	err := a.load()

	if err != nil {
		log.Fatal("Could not load: ", a.fileName, ": ", err)
	}

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./assets")))).
		Methods(http.MethodGet)

	r.Handle("/sentences", a.sentences()).
		Methods(http.MethodGet)
	r.Handle("/sentences", a.newSentence()).
		Methods(http.MethodPost)

	r.Handle("/sentences/{id}/move-up", a.moveSentenceUp()).
		Methods(http.MethodPatch)
	r.Handle("/sentences/{id}/move-down", a.moveSentenceDown()).
		Methods(http.MethodPatch)

	r.Handle("/sentences/{id}/edit", a.sentenceEditView()).
		Methods(http.MethodGet)
	r.Handle("/sentences/{id}", a.sentenceEdit()).
		Methods(http.MethodPut)

	r.Handle("/", a.index()).
		Methods(http.MethodGet)

	r.Handle("/sentences/{id}", a.getSentence()).
		Methods(http.MethodGet)

	r.Handle("/sentences/{id}", a.deleteSentence()).
		Methods(http.MethodDelete)

	r.Handle("/sentences/{id}/inspector", a.loadInspector()).
		Methods(http.MethodGet)

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
