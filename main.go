package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/amrojjeh/arabic-tags/speech"
)

// TODO(Amr Ojjeh): Consider moving into a context
var templates = template.Must(template.ParseFiles("index.html"))

type mainHandle struct {
	Sentences []speech.Sentence
}

func (m mainHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", m.Sentences)
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

	// TODO(Amr Ojjeh): Load from a json file
	// TODO(Amr Ojjeh): Add an option to connect to the db (requires auth)
	handle := mainHandle{[]speech.Sentence{{speech.Word{{
		Value: "هذا",
		Case: speech.CaseType{
			Type:      speech.CaseNA,
			Indicator: speech.IndicatorDammah,
		},
	}},
		speech.Word{
			{
				Value: "بيتُ",
				Case: speech.CaseType{
					Type:      speech.CaseNominative,
					Indicator: speech.IndicatorDammah,
				},
			},
			{
				Value: "ه",
				Case: speech.CaseType{
					Type:      speech.CaseNA,
					Indicator: speech.IndicatorNA,
				},
			},
		},
	}}}

	mux.Handle("/", handle)
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./assets"))))
	mux.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("analyze.html")
		if err != nil {
			panic(err)
		}
		w.Write(data)
	})

	log.Println("Listening on http://127.0.0.1:8080")
	log.Fatal(s.ListenAndServe())
}
