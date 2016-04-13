package web

import (
	"fmt"
	"github.com/aichaos/scarecrow/src/models"
	"github.com/gorilla/mux"
	"net/http"
)

func (g *appContext) StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Testing...
	test := models.Test{}
	g.db.FirstOrInit(&test, models.Test{Count: 0})
	test.Count = test.Count + 1
	g.db.Save(&test)
	w.Write([]byte(fmt.Sprintf("OK %d\n", test.Count)))
}

func (g *appContext) AdminSetupHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf(`{"status": "ok"}`)));
}

// registerRoutes registers all the HTTP routes for the web server.
func registerRoutes(g *appContext) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/static/{rest:.*}", g.StaticHandler)
	r.HandleFunc("/v1/status", g.StatusHandler)
	r.HandleFunc("/v1/admin/setup", g.AdminSetupHandler)

	// The index route and all other routes are handled by the React app
	r.HandleFunc("/", g.IndexHandler)
	r.HandleFunc("/{rest:.*}", g.IndexHandler)
	return r
}
