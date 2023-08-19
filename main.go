package main

import (
	"net/http"
	"os"
	"log"
	"time"
	"html/template"
)

var templates = template.Must(template.ParseFiles("index.html"))


type SpeechType int

type Function int

type Case struct {
	class     CaseClass
	explicit  bool
	indicator string
	index int // the index in which the case is explicit within the word
}

type CaseClass int

const (
	Noun     SpeechType = iota // اسم
	Verb                       // فعل
	Particle                   // حرف
)

const (
	Subject Function = iota
	Object
)

const (
	NA CaseClass = iota // not applicable
	Nominative
	Accusative
	Genetive
	// Add case for jazm
)

type Word struct {
	value    string
	speech   SpeechType
	function Function
	cas      Case // so that it doesn't conflict with case; also how it's named in CAMeL tools
	tokens   []int
}

type Sentence []Word

type mainHandle struct {
	Sentences []Sentence
}

func (m mainHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func main() {
	mux := http.NewServeMux()
	s := http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	defer s.Close()

	handle := mainHandle{nil}

	mux.Handle("/", handle)
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./assets"))))
	mux.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("analyze.html")
		if err != nil {
			panic(err)
		}
		w.Write(data)
	})

	log.Println("Listening on localhost:8080")
	log.Fatal(s.ListenAndServe())
}
