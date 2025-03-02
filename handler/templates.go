package handler

import (
	"net/http"
	"text/template"
)

var templates = template.Must(template.ParseGlob("web/templates/*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", nil)
	if err != nil {
		http.Error(w, "Failed to parse files", http.StatusInternalServerError)
	}
}

func (u *UserHandler) Home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}
