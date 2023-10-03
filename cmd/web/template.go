package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/amrojjeh/arabic-tags/internal/models"
)

type templateData struct {
	Excerpt models.Excerpt
	Form    any
}

func (app *application) cacheTemplates() error {
	app.page = make(map[string]*template.Template)

	names, err := filepath.Glob("./ui/html/pages/*")
	if err != nil {
		return err
	}

	for _, name := range names {
		baseName := filepath.Base(name)

		base, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return err
		}

		app.page[baseName], err = base.ParseFiles(name)
		if err != nil {
			return err
		}
		app.logger.Info("page cached", slog.String("name", baseName))
	}

	return nil
}

func (app *application) renderTemplate(w http.ResponseWriter, page string,
	code int, data templateData) {
	template, ok := app.page[page]
	if !ok {
		app.serverError(w, errors.New(
			fmt.Sprintf("Page %v does not exist", page)))
		return
	}

	buffer := bytes.Buffer{}
	err := template.ExecuteTemplate(&buffer, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Ignoring error as it's unlikely to occur
	w.WriteHeader(code)
	_, err = buffer.WriteTo(w)
}
