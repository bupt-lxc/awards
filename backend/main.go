package main

import (
	"awards-management/config"
	"awards-management/database"
	"awards-management/middleware"
	"awards-management/routes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	if err := config.LoadConfig("config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 创建Gin引擎
	r := gin.Default()

	// 设置请求体大小限制为200MB，以支持批量上传
	r.MaxMultipartMemory = 200 << 20 // 200MB
	// 设置JSON请求体大小限制
	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 100<<20) // 100MB
		c.Next()
	})

	// 配置CORS
	corsConfig := cors.Config{
		AllowOrigins:     config.AppConfig.CORS.AllowOrigins,
		AllowMethods:     config.AppConfig.CORS.AllowMethods,
		AllowHeaders:     config.AppConfig.CORS.AllowHeaders,
		AllowCredentials: config.AppConfig.CORS.AllowCredentials,
	}
	r.Use(cors.New(corsConfig))

	// 添加监控中间件
	r.Use(middleware.MonitorMiddleware())
	r.Use(middleware.UserActivityMiddleware())

	// 静态资源服务：将uploads/certificates映射为/static/certificates
	r.Static("/static/certificates", "uploads/certificates")

	// 设置路由
	routes.SetupRoutes(r)

	// 启动信息
	printStartupInfo()

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	log.Printf("服务器启动在: %s", addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 在 goroutine 中启动，主线程等待信号
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	// 30 秒超时优雅关闭 HTTP 服务
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// 关闭数据库连接
	if err := database.CloseDatabase(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("服务器已关闭")
}

func printStartupInfo() {
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("获奖信息管理系统后端API (Go版本)")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("启动中...")
	fmt.Printf("API地址: http://%s:%d\n", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	fmt.Printf("健康检查: http://%s:%d/health\n", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	fmt.Println(strings.Repeat("=", 50))
}
