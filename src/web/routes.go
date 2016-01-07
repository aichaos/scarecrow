package web

import (
	"fmt"
	"net/http"
	"github.com/aichaos/scarecrow/src/models"
	"github.com/gorilla/mux"
)

func (g *appContext) StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Testing...
	test := models.Test{}
	g.db.FirstOrInit(&test, models.Test{Count: 0})
	test.Count = test.Count + 1
	g.db.Save(&test)
	w.Write([]byte(fmt.Sprintf("OK %d\n", test.Count)))
}

// registerRoutes registers all the HTTP routes for the web server.
func registerRoutes(g *appContext) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", g.IndexHandler)
	r.HandleFunc("/static/{rest:.*}", g.StaticHandler)
	r.HandleFunc("/v1/status", g.StatusHandler)
	return r
}
