package controllers

import (
	"fmt"
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
		ctx.JSON(http.StatusInternalServerError, userResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    err.Error(),
		})
		return
	}

	password := user.Password

	err = db.DB.Select("id,name,email,password").Where("email=?", user.Email).Take(&user).Error
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
		ctx.JSON(http.StatusUnauthorized, userResponse{
			StatusCode: http.StatusUnauthorized,
			Status:     "error",
			Message:    "Invalid Credentials!",
		})
		return
	}

	// // Generating Token JWT
	// tokenJWT, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"id":    user.Id,
	// 	"name":  user.Name,
	// 	"email": user.Email,
	// 	"iat":   time.Now().Unix(),
	// 	"exp":   time.Now().Add(time.Minute * 5).Unix(),
	// }).SignedString([]byte(helper.GetEnv("JWT_SECRET", "")))
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, userResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Status:     "success",
	// 		Message:    err.Error(),
	// 	})
	// 	return
	// }

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

func (db *UserController) Delete(ctx *gin.Context) {
	params := ctx.Query("id")
	if params == "" {
		ctx.JSON(http.StatusBadRequest, userResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "error",
			Message:    "request not valid",
		})
		return
	}

	var checkUsers int64
	db.DB.Table("users").Where("id=?", params).Count(&checkUsers)
	if checkUsers <= 0 {
		ctx.JSON(http.StatusNotFound, userResponse{
			StatusCode: http.StatusNotFound,
			Status:     "error",
			Message:    "Users not Found",
		})
		return
	}

	db.DB.Table("users").Where("id=?", params).Delete(&models.User{})

	ctx.JSON(http.StatusOK, userResponse{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "Users success Delete!",
	})
}

func (db *UserController) Update(ctx *gin.Context) {
	paramsId := ctx.Query("id")
	user := models.User{}
	if paramsId == "" {
		ctx.JSON(http.StatusBadRequest, userResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "error",
			Message:    "request not valid",
		})
		return
	}

	checkUsers := db.DB.Table("users").Where("id=?", paramsId).First(&user).RowsAffected
	if checkUsers <= 0 {
		ctx.JSON(http.StatusNotFound, userResponse{
			StatusCode: http.StatusNotFound,
			Status:     "error",
			Message:    "Users not Found",
		})
		return
	}

	oldValue := user

	ctx.ShouldBindJSON(&user)

	user.Id = oldValue.Id
	if user.Email != oldValue.Email {
		var checkEmail int64
		db.DB.Table("users").Where("email=?", user.Email).Count(&checkEmail)
		if checkEmail > 0 {
			ctx.JSON(http.StatusConflict, userResponse{
				StatusCode: http.StatusConflict,
				Status:     "error",
				Message:    "update failed,Email exist",
			})
			return
		}
	}

	fmt.Println(user.Password)
	fmt.Println(oldValue.Password)
	if user.Password != oldValue.Password {
		hashNewPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, userResponse{
				StatusCode: http.StatusInternalServerError,
				Status:     "error",
				Message:    err.Error(),
			})
			return
		}
		user.Password = string(hashNewPassword)
	}

	db.DB.Save(&user)

	ctx.JSON(http.StatusOK, userResponse{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "Users updated!",
		Data:       user,
	})
}

func (db *UserController) GetAllUsers(ctx *gin.Context) {
	users := []models.User{}

	err := db.DB.Table("users").Select("id,name,email,created_at,updated_at").Find(&users).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, userResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, userResponse{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "All Data Users",
		Data:       users,
	})

}
