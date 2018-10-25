package ui

import "github.com/gin-gonic/gin"

func RegistManager(serverEngine *gin.Engine) {
	serverEngine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
