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
	StatusCode int         `json:StatusCode`
	Status     string      `json:status`
	Message    string      `json:message`
	Data       interface{} `json:data,omitempty`
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
