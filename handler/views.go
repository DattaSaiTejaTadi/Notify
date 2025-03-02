package handler

import "net/http"

func (h *handler) Index(w http.ResponseWriter, r *http.Request) {

	renderTemplate(w, "index.html")
}
func (h *handler) TestOnline(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "test2.html")
}
func (h *handler) AdminHome(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home-admin.html")
}
func (h *handler) L1Homepage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home-senior.html")
}
func (h *handler) L2Homepage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home-member.html")
}
func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl, nil)

	if err != nil {
		http.Error(w, "Error handling template", http.StatusInternalServerError)
	}
}
