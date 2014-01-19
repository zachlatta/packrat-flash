package packrat

import (
	"html/template"
	"io"
	"net/http"
	"strings"

	"appengine"
	"appengine/blobstore"
)

var templates = template.Must(template.ParseFiles(
	"packrat/application.html",
	"packrat/index.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, i interface{}) {
	if err := templates.ExecuteTemplate(w, tmpl+".html", i); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func serveError(c appengine.Context, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "Internal Server Error")
	c.Errorf("%v", err)
}

func init() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/new", handleSubmit)
	http.HandleFunc("/game/", handleGame)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	uploadURL, err := blobstore.UploadURL(c, "/new", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text/html")
	renderTemplate(w, "index", uploadURL)
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	blobs, _, err := blobstore.ParseUpload(r)
	if err != nil {
		serveError(c, w, err)
		return
	}
	file := blobs["file"]
	if len(file) == 0 {
		c.Errorf("No file uploaded.")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/game/"+string(file[0].BlobKey), http.StatusFound)
}

func handleGame(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	blobKey := parts[len(parts)-1]
	blobstore.Send(w, appengine.BlobKey(blobKey))
}
