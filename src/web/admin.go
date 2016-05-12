package web

import (
	"fmt"
	"github.com/aichaos/scarecrow/src/db"
	"github.com/aichaos/scarecrow/src/log"
	"github.com/aichaos/scarecrow/src/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthUser logs a user in via their session cookie.
func AuthUser(c *gin.Context, loggedIn bool) {
	session := sessions.Default(c)
	session.Set("loggedIn", loggedIn)
	session.Save()
}

// AdminAuthHandler handles logging in as the admin user.
func AdminAuthHandler(c *gin.Context) {
	// JSON params to this endpoint.
	var params struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	authError := gin.H{
		"status":  "error",
		"message": "Access denied.",
	}

	if c.BindJSON(&params) == nil {
		// Get the app config from the DB.
		DB := db.GetInstance()
		appconfig := models.AppConfig{}
		DB.Driver.First(&appconfig)

		// Username must match.
		if params.Username != appconfig.Username {
			c.JSON(200, authError)
			return
		}

		// Check the admin password.
		err := db.CheckPassword(params.Password, appconfig.Password)
		if err == nil {
			// Log them in.
			AuthUser(c, true)
			c.JSON(200, gin.H{"status": "ok"})
		} else {
			c.JSON(200, authError)
		}
	}
}

// AdminDeAuthHandler handles the Log Out event.
func AdminDeAuthHandler(c *gin.Context) {
	AuthUser(c, false)
	c.JSON(200, gin.H{"status": "ok"})
}

func AdminSetupHandler(c *gin.Context) {
	// JSON params to this endpoint.
	var params struct {
		DBType    string `json:"dbType"`
		DBString  string `json:"dbString"`
		AdminName string `json:"adminName"`
		AdminPass string `json:"adminPassword1"`
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
				"error":  fmt.Sprintf("%s", err),
			})
			return
		}

		// It works! Save the settings to disk.
		DB.SaveConfig()

		// Initialize the app configuration.
		appconfig := models.AppConfig{}
		DB.Driver.First(&appconfig)

		// Hash the password.
		hashed, err := db.HashPassword(params.AdminPass)
		if err != nil {
			log.Error("Couldn't hash admin password: %s", err)
			c.JSON(500, gin.H{
				"status": "error",
				"error":  fmt.Sprintf("Password hashing error: %s", err),
			})
			return
		}

		// Set the app config and save it.
		appconfig.Username = params.AdminName
		appconfig.Password = hashed
		appconfig.Name = "scarecrow"
		appconfig.Replies = "replies/standard"

		if DB.Driver.NewRecord(appconfig) {
			log.Debug("Creating NEW AppConfig record.")
			appconfig.Initialized = true
			DB.Driver.Create(&appconfig)
		} else {
			log.Debug("Updating EXISTING AppConfig record.")
			DB.Driver.Save(&appconfig)
		}

		// Log them in as the admin user.
		AuthUser(c, true)
		log.Debug("Sent user a loggedIn session.")

		c.JSON(200, gin.H{"status": "ok"})
	} else {
		c.JSON(400, gin.H{
			"status": "error",
			"error":  "Bad parameters.",
		})
	}
}
