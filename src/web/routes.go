package web

import (
	"fmt"
	"net/http"
	"github.com/aichaos/scarecrow/src/models"
	"github.com/gorilla/mux"
)

// TODO: move to a different file
func (h *appHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!\n"))
}

func (h *appHandler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Testing...
	test := models.Test{}
	h.db.FirstOrInit(&test, models.Test{Count: 0})
	test.Count = test.Count + 1
	h.db.Save(&test)
	w.Write([]byte(fmt.Sprintf("OK %d\n", test.Count)))
}

// registerRoutes registers all the HTTP routes for the web server.
func registerRoutes(h *appHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", h.IndexHandler)
	r.HandleFunc("/v1/status", h.StatusHandler)
	return r
}
