package router

import (
	"github.com/gin-gonic/gin"
	"github.com/r1is/ReportSystem/controller"
	"github.com/r1is/ReportSystem/middleware"
	"golang.org/x/time/rate"
	"time"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	// 创建一个新的限流器，每分钟最多3个请求
	limiter := middleware.NewRateLimiter(rate.Every(time.Minute), 3)

	// 将限流器应用到register路由
	r.POST("/api/auth/register", limiter.LimitHandler(), controller.Register)

	r.POST("/api/auth/login", controller.Login)
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.UserInfo)
	return r
}
