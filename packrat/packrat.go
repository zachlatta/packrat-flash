package packrat

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles(
	"packrat/application.html",
	"packrat/index.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, i interface{}) {
	if err := templates.ExecuteTemplate(w, tmpl+".html", i); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func init() {
	http.HandleFunc("/", root)
}

func root(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}
