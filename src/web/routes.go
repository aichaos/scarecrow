package web

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func StatusHandler(c *gin.Context) {
	session := sessions.Default(c)
	var count int
	v := session.Get("count")
	if v == nil {
		count = 0
	} else {
		count = v.(int)
		count += 1
	}
	session.Set("count", count)
	session.Save()
	c.JSON(200, gin.H{"status": "ok", "count": count})
}
