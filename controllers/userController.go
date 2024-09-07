package controllers

import (
	"net/http"

	"github.com/HaikalRFadhilahh/course-golang/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

type userResponse struct {
	StatusCode int         `json:"statusCode"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func (db *UserController) Login(ctx *gin.Context) {
	user := models.User{}
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"statusCode": 500,
			"status":     "error",
			"message":    err.Error(),
		})
		return
	}

	password := user.Password

	err = db.DB.Where("email=?", user.Email).Take(&user).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"statusCode": 500,
			"status":     "error",
			"message":    err.Error(),
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
			"statusCode": 401,
			"status":     "error",
			"message":    "Credentials invalid",
		})
		return
	}

	ctx.JSON(http.StatusOK, userResponse{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "Authenticated",
		Data:       user,
	})
}

func (db *UserController) Register(ctx *gin.Context) {
	user := models.User{}
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, userResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    err.Error(),
		})
		return
	}

	checkUsers := db.DB.Table("users").Where("email=?", user.Email).Take(&user).RowsAffected > 0
	if checkUsers {
		ctx.JSON(http.StatusConflict, userResponse{
			StatusCode: http.StatusConflict,
			Status:     "error",
			Message:    "Email exist",
		})
		return
	}

	user.Role = "Employees"

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, userResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    err.Error(),
		})
		return
	}

	user.Password = string(hashPassword)

	err = db.DB.Table("users").Create(&user).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, userResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, userResponse{
		StatusCode: http.StatusCreated,
		Status:     "success",
		Message:    "Created",
	})
}
