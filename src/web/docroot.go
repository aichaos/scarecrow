package web

// This is the "catch-all" handler for all non-API routes. It attempts to serve
// files from the document root on disk.

import (
	"bytes"
	"github.com/aichaos/scarecrow/src/db"
	"github.com/aichaos/scarecrow/src/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"text/template"
)

const DOCUMENT_ROOT string = "http/public"

func IndexHandler(c *gin.Context) {
	DB := db.GetInstance()
	appconfig := models.AppConfig{}
	DB.Driver.First(&appconfig)

	// Get their session info.
	session := sessions.Default(c)
	var loggedIn bool = false

	v := session.Get("loggedIn")
	if v != nil {
		loggedIn = v.(bool)
	}

	vars := struct {
		Title       string
		Initialized bool // This is false until the DB has been initialized
		LoggedIn    bool
	}{
		Title:       "Hello world",
		Initialized: appconfig.Initialized,
		LoggedIn:    loggedIn,
	}

	t, err := template.ParseFiles("http/public/index.html")
	if err != nil {
		panic("Error parsing index template!")
	}

	buf := new(bytes.Buffer)
	t.Execute(buf, vars)
	c.Render(200, render.Data{
		ContentType: "text/html",
		Data:        buf.Bytes(),
	})
}
