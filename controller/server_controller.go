package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DemoRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Server is working"})
	}
}
