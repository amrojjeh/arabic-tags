package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/amrojjeh/arabic-tags/internal/models"
	"github.com/amrojjeh/arabic-tags/internal/speech"
	"github.com/amrojjeh/arabic-tags/ui"
	"github.com/google/uuid"
)

type symbol struct {
	PDamma    string
	PDammatan string
	PFatha    string
	PFathatan string
	PKasra    string
	PKasratan string
	PSukoon   string

	Damma    string
	Dammatan string
	Fatha    string
	Fathatan string
	Kasra    string
	Kasratan string
	Sukoon   string

	PShadda string
	Shadda  string
}

type templateData struct {
	Excerpt             models.Excerpt
	Form                any
	Error               string
	GrammaticalTags     []string
	Host                string
	ExcerptShared       bool
	TSelectedWord       int
	AcceptedPunctuation string
	Sym                 symbol
}

func newTemplateData(r *http.Request) (templateData, error) {
	err := r.ParseForm()
	if err != nil {
		return templateData{}, err
	}
	data := templateData{
		Error:               r.Form.Get("error"),
		GrammaticalTags:     speech.GrammaticalTags,
		Host:                r.Host,
		ExcerptShared:       r.Form.Get("share") == "true",
		TSelectedWord:       0,
		AcceptedPunctuation: speech.Punctuation.String(),
		Sym: symbol{
			PDamma:    speech.Placeholder + speech.Damma,
			PDammatan: speech.Placeholder + speech.Dammatan,
			PFatha:    speech.Placeholder + speech.Fatha,
			PFathatan: speech.Placeholder + speech.Fathatan,
			PKasra:    speech.Placeholder + speech.Kasra,
			PKasratan: speech.Placeholder + speech.Kasratan,
			PSukoon:   speech.Placeholder + speech.Sukoon,
			PShadda:   speech.Placeholder + speech.Shadda,
			Damma:     speech.Damma,
			Dammatan:  speech.Dammatan,
			Fatha:     speech.Fatha,
			Fathatan:  speech.Fathatan,
			Kasra:     speech.Kasra,
			Kasratan:  speech.Kasratan,
			Sukoon:    speech.Sukoon,
			Shadda:    speech.Shadda,
		},
	}

	if r.Context().Value(excerptContextKey) != nil {
		data.Excerpt = r.Context().Value(excerptContextKey).(models.Excerpt)
	}

	return data, nil
}

func JSONFunc(s any) (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	return string(b[:]), nil
}

func IdFunc(s uuid.UUID) string {
	return strings.ReplaceAll(s.String(), "-", "")
}

func EvenFunc(s int) bool {
	return s%2 == 0
}

func OddFunc(s int) bool {
	return s%2 != 0
}

func (app *application) cacheTemplates() error {
	app.page = make(map[string]*template.Template)
	funcs := template.FuncMap{
		"json": JSONFunc,
		"id":   IdFunc,
		"even": EvenFunc,
		"odd":  OddFunc,
	}

	names, err := fs.Glob(ui.Files, "html/pages/*")
	if err != nil {
		return err
	}

	for _, name := range names {
		baseName := filepath.Base(name)

		base := template.New(name).Funcs(funcs)
		var err error

		if !strings.HasPrefix(baseName, "htmx-") {
			base, err = base.ParseFS(ui.Files, "html/base.tmpl")
			if err != nil {
				return err
			}
		}

		partials, err := fs.Glob(ui.Files, "html/partials/*")
		if err != nil {
			return err
		}

		for _, name := range partials {
			base, err = base.ParseFS(ui.Files, name)
			if err != nil {
				return err
			}
		}

		app.page[baseName], err = base.ParseFS(ui.Files, name)
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
	var err error
	if strings.HasPrefix(page, "htmx-") {
		err = template.ExecuteTemplate(&buffer, "htmx", data)
	} else {
		err = template.ExecuteTemplate(&buffer, "base", data)
	}
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Ignoring error as it's unlikely to occur
	w.WriteHeader(code)
	_, err = buffer.WriteTo(w)
}
