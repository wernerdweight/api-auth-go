package routes

import "github.com/gin-gonic/gin"

func Register(r *gin.Engine) {
	r.POST("/authenticate", authenticateHandler)
}

func authenticateHandler(c *gin.Context) {
	// TODO:
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
