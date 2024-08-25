package routes

import (
	"github.com/Uttkarsh-raj/PS-1708/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/", controller.DemoRoute())
}
