package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Router
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Welcome to Golang App")
	})

	router.Static("/public", "./attachments")

	router.Run("127.0.0.1:8080")
}
