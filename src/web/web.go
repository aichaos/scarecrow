package web

import (
	"fmt"
	"github.com/aichaos/scarecrow/src/types"
	"github.com/gin-gonic/gin"
)

type appContext struct {
	// db *gorm.DB
}

func StartServer(config types.WebConfig) {
	r := gin.Default()

	r.Static("/static", "./http/public/static")

	r.GET("/v1/status", StatusHandler)
	r.POST("/v1/admin/setup", AdminSetupHandler)

	// The index and catch-all handlers go to the index.html page.
	r.GET("/", IndexHandler)
	r.NoRoute(IndexHandler)

	r.Run(fmt.Sprintf("%s:%d", config.Host, config.Port))
}
