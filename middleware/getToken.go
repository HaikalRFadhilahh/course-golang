package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type errorReponse struct {
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

func CheckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestToken := ctx.GetHeader("Authorization")

		if requestToken == "" {
			ctx.JSON(http.StatusForbidden, errorReponse{
				StatusCode: http.StatusForbidden,
				Status:     "error",
				Message:    "Invalid Token!!",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
