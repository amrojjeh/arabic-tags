package main

import (
	"bytes"
	"html/template"
	"net/http"
)

func (app *application) addArabic(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/add.tmpl",
	}
	tmpl, err := template.ParseFiles(files...)
	// TODO(Amr Ojjeh): Move to a helper function to log
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	buffer := bytes.Buffer{}
	err = tmpl.ExecuteTemplate(&buffer, "base", nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	// Ignoring error as it's unlikely to occur
	_, err = buffer.WriteTo(w)
}
