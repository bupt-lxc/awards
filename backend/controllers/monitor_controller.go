package controllers

import (
	"awards-management/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MonitorController 监控控制器
type MonitorController struct {
	monitorService *services.MonitorService
}

// NewMonitorController 创建监控控制器实例
func NewMonitorController() *MonitorController {
	return &MonitorController{
		monitorService: services.GetMonitorService(),
	}
}

// GetMonitorData 获取监控数据
func (mc *MonitorController) GetMonitorData(c *gin.Context) {
	data := mc.monitorService.GetMonitorData()
	c.JSON(http.StatusOK, data)
}

// GetDatabaseStats 获取数据库统计
func (mc *MonitorController) GetDatabaseStats(c *gin.Context) {
	stats, err := mc.monitorService.GetDatabaseStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取数据库统计失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// ResetStats 重置统计
func (mc *MonitorController) ResetStats(c *gin.Context) {
	mc.monitorService.ResetStats()
	c.JSON(http.StatusOK, gin.H{
		"message": "统计数据已重置",
	})
}

// GetOnlineUsers 获取在线用户
func (mc *MonitorController) GetOnlineUsers(c *gin.Context) {
	data := mc.monitorService.GetMonitorData()
	c.JSON(http.StatusOK, gin.H{
		"online_users": data.OnlineUsers,
		"online_count": data.OnlineCount,
		"timestamp":    data.Timestamp,
	})
}

// KickUser 踢出用户
func (mc *MonitorController) KickUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户ID不能为空",
		})
		return
	}

	mc.monitorService.RemoveOnlineUser(userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "用户已踢出",
	})
}

// GetSystemHealth 获取系统健康状态
func (mc *MonitorController) GetSystemHealth(c *gin.Context) {
	data := mc.monitorService.GetMonitorData()

	// 计算健康分数
	healthScore := 100.0
	issues := []string{}

	// 检查内存使用率
	if data.SystemMetrics.MemoryUsage > 80 {
		healthScore -= 20
		issues = append(issues, "内存使用率过高")
	}

	// 检查错误率
	if data.RequestStats.TotalRequests > 0 {
		errorRate := float64(data.RequestStats.ErrorRequests) / float64(data.RequestStats.TotalRequests) * 100
		if errorRate > 10 {
			healthScore -= 30
			issues = append(issues, "错误率过高")
		}
	}

	// 检查响应时间
	if data.RequestStats.AvgResponseTime > 2000 { // 2秒
		healthScore -= 15
		issues = append(issues, "响应时间过长")
	}

	// 检查数据库连接
	if data.SystemMetrics.DBConnections > 50 {
		healthScore -= 10
		issues = append(issues, "数据库连接数过多")
	}

	// 检查Goroutine数量
	if data.SystemMetrics.Goroutines > 1000 {
		healthScore -= 10
		issues = append(issues, "Goroutine数量过多")
	}

	healthStatus := "healthy"
	if healthScore < 70 {
		healthStatus = "unhealthy"
	} else if healthScore < 90 {
		healthStatus = "warning"
	}

	c.JSON(http.StatusOK, gin.H{
		"health_score":   healthScore,
		"health_status":  healthStatus,
		"issues":         issues,
		"system_metrics": data.SystemMetrics,
		"request_stats":  data.RequestStats,
		"timestamp":      data.Timestamp,
	})
}

// GetRequestHistory 获取请求历史
func (mc *MonitorController) GetRequestHistory(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 200 {
		limit = 50
	}

	data := mc.monitorService.GetMonitorData()
	history := data.RequestStats.RequestHistory

	// 限制返回数量
	if len(history) > limit {
		history = history[len(history)-limit:]
	}

	c.JSON(http.StatusOK, gin.H{
		"request_history": history,
		"total_count":     len(data.RequestStats.RequestHistory),
		"returned_count":  len(history),
		"timestamp":       data.Timestamp,
	})
}
