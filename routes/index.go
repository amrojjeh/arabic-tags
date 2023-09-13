package routes

import (
	"net/http"

	"github.com/amrojjeh/arabic-tags/views"
	"github.com/gorilla/mux"
)

func Index(r *mux.Router) {
	r.Handle("/", getIndex())
}

func getIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		views.Page(views.Prop{
			Title:   "Test",
			Content: nil,
		}).Render(w)
	}
}
