package web

// This is the "catch-all" handler for all non-API routes. It attempts to serve
// files from the document root on disk.

import (
	"fmt"
	"text/template"
	"net/http"
	"os"
)

const DOCUMENT_ROOT string = "http/public"

func (g *appContext) IndexHandler(w http.ResponseWriter, r *http.Request) {
	vars := struct {
		Title       string
		Initialized bool   // This is false until the DB has been initialized
	}{
		Title: "Hello world",
		Initialized: false,
	}

	t, err := template.ParseFiles("http/public/index.html")
	if err != nil {
		panic("Error parsing index template!")
	}
	t.Execute(w, vars)
}

func (g *appContext) StaticHandler(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path

	// File exists?
	filePath := fmt.Sprintf("%s/%s", DOCUMENT_ROOT, uri)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		w.Write([]byte("404 Not Found"))
		return
	}

	http.ServeFile(w, r, filePath)
}
