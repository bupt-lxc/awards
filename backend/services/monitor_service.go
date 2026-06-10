package services

import (
	"awards-management/database"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

// MonitorService 监控服务
type MonitorService struct {
	onlineUsers   map[string]*OnlineUser
	requestStats  *RequestStats
	systemMetrics *SystemMetrics
	mu            sync.RWMutex
	startTime     time.Time
}

// OnlineUser 在线用户
type OnlineUser struct {
	UserID     string    `json:"user_id"`
	Username   string    `json:"username"`
	IP         string    `json:"ip"`
	UserAgent  string    `json:"user_agent"`
	LastActive time.Time `json:"last_active"`
	LoginTime  time.Time `json:"login_time"`
	IsAdmin    bool      `json:"is_admin"`
}

// RequestStats 请求统计
type RequestStats struct {
	TotalRequests   int64         `json:"total_requests"`
	SuccessRequests int64         `json:"success_requests"`
	ErrorRequests   int64         `json:"error_requests"`
	AvgResponseTime float64       `json:"avg_response_time"`
	QPS             float64       `json:"qps"`
	ResponseTimes   []float64     `json:"-"` // 不序列化到JSON
	LastResetTime   time.Time     `json:"last_reset_time"`
	RequestHistory  []RequestData `json:"request_history"`
}

// RequestData 请求数据
type RequestData struct {
	Timestamp    time.Time `json:"timestamp"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	ResponseTime float64   `json:"response_time"`
	Status       int       `json:"status"`
	IP           string    `json:"ip"`
}

// SystemMetrics 系统指标
type SystemMetrics struct {
	CPUUsage      float64   `json:"cpu_usage"`
	MemoryUsage   float64   `json:"memory_usage"`
	MemoryTotal   uint64    `json:"memory_total"`
	MemoryUsed    uint64    `json:"memory_used"`
	Goroutines    int       `json:"goroutines"`
	DBConnections int       `json:"db_connections"`
	Uptime        int64     `json:"uptime"`
	Timestamp     time.Time `json:"timestamp"`
}

// MonitorData 监控数据
type MonitorData struct {
	OnlineUsers   []OnlineUser   `json:"online_users"`
	OnlineCount   int            `json:"online_count"`
	RequestStats  *RequestStats  `json:"request_stats"`
	SystemMetrics *SystemMetrics `json:"system_metrics"`
	Timestamp     time.Time      `json:"timestamp"`
}

var (
	monitorService *MonitorService
	once           sync.Once
)

// NewMonitorService 创建监控服务实例
func NewMonitorService() *MonitorService {
	once.Do(func() {
		monitorService = &MonitorService{
			onlineUsers:   make(map[string]*OnlineUser),
			requestStats:  &RequestStats{LastResetTime: time.Now()},
			systemMetrics: &SystemMetrics{},
			startTime:     time.Now(),
		}

		// 启动定期更新系统指标
		go monitorService.updateSystemMetrics()

		// 启动定期清理过期用户
		go monitorService.cleanupExpiredUsers()
	})
	return monitorService
}

// GetMonitorService 获取监控服务实例
func GetMonitorService() *MonitorService {
	return NewMonitorService()
}

// AddOnlineUser 添加在线用户
func (ms *MonitorService) AddOnlineUser(userID, username, ip, userAgent string, isAdmin bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.onlineUsers[userID] = &OnlineUser{
		UserID:     userID,
		Username:   username,
		IP:         ip,
		UserAgent:  userAgent,
		LastActive: time.Now(),
		LoginTime:  time.Now(),
		IsAdmin:    isAdmin,
	}

	log.Printf("[监控] 用户上线: %s (%s)", username, ip)
}

// RemoveOnlineUser 移除在线用户
func (ms *MonitorService) RemoveOnlineUser(userID string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if user, exists := ms.onlineUsers[userID]; exists {
		log.Printf("[监控] 用户下线: %s", user.Username)
		delete(ms.onlineUsers, userID)
	}
}

// UpdateUserActivity 更新用户活动时间
func (ms *MonitorService) UpdateUserActivity(userID string) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	if user, exists := ms.onlineUsers[userID]; exists {
		user.LastActive = time.Now()
	}
}

// RecordRequest 记录请求
func (ms *MonitorService) RecordRequest(method, path string, responseTime float64, status int, ip string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	now := time.Now()

	// 更新请求统计
	ms.requestStats.TotalRequests++
	if status >= 200 && status < 400 {
		ms.requestStats.SuccessRequests++
	} else {
		ms.requestStats.ErrorRequests++
	}

	// 更新响应时间
	ms.requestStats.ResponseTimes = append(ms.requestStats.ResponseTimes, responseTime)
	if len(ms.requestStats.ResponseTimes) > 1000 { // 保持最近1000个请求的响应时间
		ms.requestStats.ResponseTimes = ms.requestStats.ResponseTimes[1:]
	}

	// 计算平均响应时间
	total := 0.0
	for _, rt := range ms.requestStats.ResponseTimes {
		total += rt
	}
	ms.requestStats.AvgResponseTime = total / float64(len(ms.requestStats.ResponseTimes))

	// 计算QPS（每秒请求数）
	elapsed := now.Sub(ms.requestStats.LastResetTime).Seconds()
	if elapsed > 0 {
		ms.requestStats.QPS = float64(ms.requestStats.TotalRequests) / elapsed
	}

	// 记录请求历史
	requestData := RequestData{
		Timestamp:    now,
		Method:       method,
		Path:         path,
		ResponseTime: responseTime,
		Status:       status,
		IP:           ip,
	}

	ms.requestStats.RequestHistory = append(ms.requestStats.RequestHistory, requestData)
	if len(ms.requestStats.RequestHistory) > 100 { // 保持最近100个请求
		ms.requestStats.RequestHistory = ms.requestStats.RequestHistory[1:]
	}
}

// GetMonitorData 获取监控数据
func (ms *MonitorService) GetMonitorData() *MonitorData {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	// 复制在线用户列表
	onlineUsers := make([]OnlineUser, 0, len(ms.onlineUsers))
	for _, user := range ms.onlineUsers {
		onlineUsers = append(onlineUsers, *user)
	}

	// 复制请求统计
	requestStats := &RequestStats{
		TotalRequests:   ms.requestStats.TotalRequests,
		SuccessRequests: ms.requestStats.SuccessRequests,
		ErrorRequests:   ms.requestStats.ErrorRequests,
		AvgResponseTime: ms.requestStats.AvgResponseTime,
		QPS:             ms.requestStats.QPS,
		LastResetTime:   ms.requestStats.LastResetTime,
		RequestHistory:  make([]RequestData, len(ms.requestStats.RequestHistory)),
	}
	copy(requestStats.RequestHistory, ms.requestStats.RequestHistory)

	// 复制系统指标
	systemMetrics := &SystemMetrics{
		CPUUsage:      ms.systemMetrics.CPUUsage,
		MemoryUsage:   ms.systemMetrics.MemoryUsage,
		MemoryTotal:   ms.systemMetrics.MemoryTotal,
		MemoryUsed:    ms.systemMetrics.MemoryUsed,
		Goroutines:    ms.systemMetrics.Goroutines,
		DBConnections: ms.systemMetrics.DBConnections,
		Uptime:        ms.systemMetrics.Uptime,
		Timestamp:     ms.systemMetrics.Timestamp,
	}

	return &MonitorData{
		OnlineUsers:   onlineUsers,
		OnlineCount:   len(onlineUsers),
		RequestStats:  requestStats,
		SystemMetrics: systemMetrics,
		Timestamp:     time.Now(),
	}
}

// updateSystemMetrics 更新系统指标
func (ms *MonitorService) updateSystemMetrics() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ms.mu.Lock()

		// 获取内存使用情况
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		ms.systemMetrics.MemoryUsed = m.Alloc
		ms.systemMetrics.MemoryTotal = m.Sys
		ms.systemMetrics.MemoryUsage = float64(m.Alloc) / float64(m.Sys) * 100
		ms.systemMetrics.Goroutines = runtime.NumGoroutine()
		ms.systemMetrics.Uptime = int64(time.Since(ms.startTime).Seconds())
		ms.systemMetrics.Timestamp = time.Now()

		// 获取数据库连接数
		if db := database.GetDB(); db != nil {
			if sqlDB, err := db.DB(); err == nil {
				stats := sqlDB.Stats()
				ms.systemMetrics.DBConnections = stats.OpenConnections
			}
		}

		// 简单的CPU使用率估算（基于GC频率）
		ms.systemMetrics.CPUUsage = float64(m.NumGC) * 0.1 // 简化计算

		ms.mu.Unlock()
	}
}

// cleanupExpiredUsers 清理过期用户（超过30分钟无活动）
func (ms *MonitorService) cleanupExpiredUsers() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		ms.mu.Lock()
		now := time.Now()
		expiredUsers := make([]string, 0)

		for userID, user := range ms.onlineUsers {
			if now.Sub(user.LastActive) > 30*time.Minute {
				expiredUsers = append(expiredUsers, userID)
			}
		}

		for _, userID := range expiredUsers {
			log.Printf("[监控] 清理过期用户: %s", ms.onlineUsers[userID].Username)
			delete(ms.onlineUsers, userID)
		}

		ms.mu.Unlock()
	}
}

// GetDatabaseStats 获取数据库统计
func (ms *MonitorService) GetDatabaseStats() (map[string]interface{}, error) {
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("数据库连接不可用")
	}

	stats := make(map[string]interface{})

	// 获取表记录数
	var personCount, groupCount, memberCount int64
	db.Model(&struct{}{}).Table("person_awards").Count(&personCount)
	db.Model(&struct{}{}).Table("group_awards").Count(&groupCount)
	db.Model(&struct{}{}).Table("group_members").Count(&memberCount)

	stats["person_awards_count"] = personCount
	stats["group_awards_count"] = groupCount
	stats["group_members_count"] = memberCount
	stats["total_records"] = personCount + groupCount + memberCount

	// 获取数据库大小（MySQL）
	var dbSize string
	db.Raw("SELECT ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) AS 'DB Size in MB' FROM information_schema.tables WHERE table_schema = DATABASE()").Scan(&dbSize)
	stats["database_size_mb"] = dbSize

	// 获取表大小
	var tableSizes []map[string]interface{}
	db.Raw(`
		SELECT 
			table_name,
			ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'Size in MB'
		FROM information_schema.tables 
		WHERE table_schema = DATABASE()
		ORDER BY (data_length + index_length) DESC
	`).Scan(&tableSizes)
	stats["table_sizes"] = tableSizes

	return stats, nil
}

// ResetStats 重置统计
func (ms *MonitorService) ResetStats() {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.requestStats = &RequestStats{LastResetTime: time.Now()}
	log.Printf("[监控] 统计数据已重置")
}
