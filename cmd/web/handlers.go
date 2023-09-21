package main

import (
	"net/http"

	"github.com/amrojjeh/arabic-tags/cmd/web/views"
)

func (app *application) addArabic(w http.ResponseWriter, r *http.Request) {
	views.Page(views.Prop{
		Title:   "Arabic",
		Content: views.AddArabic(),
	}).Render(w)
}
