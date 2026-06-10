package middleware

import (
	"awards-management/services"
	"time"

	"github.com/gin-gonic/gin"
)

// MonitorMiddleware 监控中间件
func MonitorMiddleware() gin.HandlerFunc {
	monitorService := services.GetMonitorService()

	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 计算响应时间
		responseTime := float64(time.Since(start).Nanoseconds()) / 1e6 // 转换为毫秒

		// 记录请求
		monitorService.RecordRequest(
			c.Request.Method,
			c.Request.URL.Path,
			responseTime,
			c.Writer.Status(),
			c.ClientIP(),
		)
	}
}

// UserActivityMiddleware 用户活动监控中间件
func UserActivityMiddleware() gin.HandlerFunc {
	monitorService := services.GetMonitorService()

	return func(c *gin.Context) {
		// 检查是否有用户信息
		if userID, exists := c.Get("user_id"); exists {
			monitorService.UpdateUserActivity(userID.(string))
		}

		c.Next()
	}
}
