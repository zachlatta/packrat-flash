package packrat

import (
	"html/template"
	"io"
	"net/http"

	"appengine"
	"appengine/blobstore"
)

type Game struct {
	Name  string
	Email string
	File  []byte
}

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
	http.HandleFunc("/", root)
	http.HandleFunc("/new", submitGame)
	http.HandleFunc("/serve/", handleServe)
}

func root(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	uploadURL, err := blobstore.UploadURL(c, "/new", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text/html")
	renderTemplate(w, "index", uploadURL)
}

func submitGame(w http.ResponseWriter, r *http.Request) {
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
	http.Redirect(w, r, "/serve/?blobKey="+string(file[0].BlobKey), http.StatusFound)
}

func handleServe(w http.ResponseWriter, r *http.Request) {
	blobstore.Send(w, appengine.BlobKey(r.FormValue("blobKey")))
}
