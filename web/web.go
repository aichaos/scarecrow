package web

import (
	"fmt"
	"net/http"
	"github.com/aichaos/scarecrow/types"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func StartServer(config types.WebConfig) {
	r := registerRoutes()
	http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), r)
}
