package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/r1is/ReportSystem/controller"
	"golang.org/x/time/rate"
	"net/http"
)

// 定义限流器结构体
type RateLimiter struct {
	limiter *rate.Limiter
}

// 创建一个新的限流器
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(r, b),
	}
}

// 定义限流中间件
func (rl *RateLimiter) LimitHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.limiter.Allow() {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}

// 应用限流器到指定路由
func ApplyRateLimit(route *gin.Engine, path string, limiter *RateLimiter) {
	route.POST(path, limiter.LimitHandler(), controller.Register)
}
