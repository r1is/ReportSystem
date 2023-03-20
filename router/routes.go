package router

import (
	"github.com/gin-gonic/gin"
	"github.com/r1is/ReportSystem/controller"
	"github.com/r1is/ReportSystem/middleware"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.UserInfo)
	return r
}
