package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/HaikalRFadhilahh/course-golang/config"
	"github.com/HaikalRFadhilahh/course-golang/controllers"
	"github.com/HaikalRFadhilahh/course-golang/helper"
	"github.com/HaikalRFadhilahh/course-golang/models"
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

	// Database Connection
	db, _ := config.DatabaseConnection()
	db.AutoMigrate(&models.User{}, &models.Task{})
	config.CreateOwnerAccount(db)

	// Controller
	userController := controllers.UserController{DB: db}

	// Router
	router := gin.Default()
	// Static Route for Attachment and Public File
	router.Static("/public", "./attachments")
	// API Routing
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://haik.my.id")
	})

	// Users Routing
	router.POST("/users/login", userController.Login)
	router.POST("/users/register", userController.Register)
	router.DELETE("/users/delete", userController.Delete)
	router.PUT("/users/update", userController.Update)
	router.GET("/users", userController.GetAllUsers)

	// Running Gin Server
	fmt.Print("Go Gin Gonic Running in ", connectionString)
	err := router.Run(connectionString)
	if err != nil {
		log.Fatal(err.Error())
	}
}
