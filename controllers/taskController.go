package controllers

import (
	"net/http"
	"os"

	"github.com/HaikalRFadhilahh/course-golang/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskController struct {
	DB *gorm.DB
}

type taskResponse struct {
	StatusCode int         `json:"statusCode"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func (u *TaskController) Create(ctx *gin.Context) {
	task := models.Task{}

	err := ctx.BindJSON(&task)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, taskResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "error",
			Message:    err.Error(),
		})
		return
	}

	if task.UserId != nil {
		var checkUsers int64
		u.DB.Table("users").Where("id=?", *task.UserId).Count(&checkUsers)
		if checkUsers <= 0 {
			ctx.JSON(http.StatusBadRequest, taskResponse{
				StatusCode: http.StatusBadRequest,
				Status:     "error",
				Message:    "Users id Not found!",
			})
			return
		}
	}

	err = u.DB.Table("tasks").Create(&task).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, taskResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, taskResponse{
		StatusCode: http.StatusCreated,
		Status:     "success",
		Message:    "Task Success Created",
		Data:       task,
	})
}

func (u *TaskController) Delete(ctx *gin.Context) {
	// Get Wildcard Value
	taskId := ctx.Param("id")

	// Define Models
	task := models.Task{}

	// Validate Data
	if u.DB.Table("tasks").Where("id=?", taskId).First(&task).RowsAffected <= 0 {
		ctx.JSON(http.StatusBadRequest, taskResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "error",
			Message:    "data task not found!",
		})
		return
	}

	// Remove Data
	err := u.DB.Table("tasks").Delete(&task).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, taskResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    err.Error(),
		})
		return
	}

	if task.Attachment != nil {
		os.Remove("attachments/" + *task.Attachment)
	}

	// Return Success
	ctx.JSON(http.StatusOK, taskResponse{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "Data Task success deleted!",
		Data:       task,
	})
}
