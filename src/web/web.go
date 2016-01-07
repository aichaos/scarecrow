package web

import (
	"fmt"
	"net/http"
	"github.com/aichaos/scarecrow/src/types"
	"github.com/aichaos/scarecrow/src/models"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type appHandler struct {
	db *gorm.DB
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func StartServer(config types.WebConfig) {
	// TODO: make configurable, etc.
	db, err := gorm.Open("sqlite3", "temp.sqlite")
	if err != nil {
		panic("Can't connect to DB")
	}

	db.CreateTable(&models.Test{})

	handler := &appHandler{db: &db}
	r := registerRoutes(handler)
	http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), r)
}
