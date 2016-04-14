package web

// This is the "catch-all" handler for all non-API routes. It attempts to serve
// files from the document root on disk.

import (
	"bytes"
	"text/template"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

const DOCUMENT_ROOT string = "http/public"

func IndexHandler(c *gin.Context) {
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

	buf := new(bytes.Buffer)
	t.Execute(buf, vars)
	c.Render(200, render.Data{
		ContentType: "text/html",
		Data: buf.Bytes(),
	})
}
