package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/HaikalRFadhilahh/course-golang/helper"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Read .env File
	godotenv.Load()

	// Read ENV Value
	PORT := helper.GetEnv("PORT", "3000")
	HOST := helper.GetEnv("HOST", "127.0.0.1")
	APP_MODE := helper.GetEnv("APP_MODE", "production")
	connectionString := fmt.Sprintf("%s:%s", HOST, PORT)

	if APP_MODE != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Router
	router := gin.Default()
	// Static Route for Attachment and Public File
	router.Static("/public", "./attachments")
	// API Routing
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Welcome to Golang App")
	})

	// Running Gin Server
	fmt.Print("Go Gin Gonic Running in ", connectionString)
	err := router.Run(connectionString)
	if err != nil {
		log.Fatal(err.Error())
	}
}
