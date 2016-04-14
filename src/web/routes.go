package web

import (
	"fmt"
	"github.com/aichaos/scarecrow/src/db"
	"github.com/aichaos/scarecrow/src/log"
	"github.com/aichaos/scarecrow/src/models"
	"github.com/gin-gonic/gin"
)

func StatusHandler(c *gin.Context) {
	// Testing...
	test := models.Test{}
	// g.db.FirstOrInit(&test, models.Test{Count: 0})
	test.Count = test.Count + 1
	// g.db.Save(&test)
	c.JSON(200, gin.H{"status": "ok"})
}

func AdminSetupHandler(c *gin.Context) {
	// JSON params to this endpoint.
	var params struct {
		DBType    string `json:"dbType"`
		DBString  string `json:"dbString"`
		adminName string `json:"adminName"`
		adminPass string `json:"adminPassword1"`
	}

	if c.BindJSON(&params) == nil {
		// Try the database settings.
		DB := db.GetInstance()
		DB.Config.Type = params.DBType
		DB.Config.ConnString = params.DBString
		log.Info("Trying database parameters given by user...")
		_, err := DB.Connect()
		if err != nil {
			log.Error("Database parameters didn't work: %s", err)
			c.JSON(400, gin.H{
				"status": "error",
				"error": fmt.Sprintf("%s", err),
			})
			return
		}

		// It works! Save the settings to disk.
		DB.SaveConfig()

		c.JSON(200, gin.H{"status": "ok"})
	} else {
		c.JSON(400, gin.H{"status": "error"})
	}
}
