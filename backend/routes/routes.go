package routes

import (
	authPkg "awards-management/auth"
	"awards-management/controllers"
	"awards-management/database"
	"awards-management/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	authController := controllers.NewAuthController()
	awardController := controllers.NewAwardController()
	monitorController := controllers.NewMonitorController()

	// 根路径
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":    "获奖信息管理系统API",
			"version":    "1.0.0",
			"docs":       "/docs",
			"redoc":      "/redoc",
			"api_prefix": "/api/v1",
		})
	})

	// 健康检查（含数据库连接检测）
	r.GET("/health", func(c *gin.Context) {
		dbStatus := "healthy"
		if db := database.GetDB(); db != nil {
			if sqlDB, err := db.DB(); err == nil {
				if pingErr := sqlDB.Ping(); pingErr != nil {
					dbStatus = "unhealthy: " + pingErr.Error()
				}
			} else {
				dbStatus = "unhealthy: " + err.Error()
			}
		} else {
			dbStatus = "unhealthy: no database connection"
		}

		status := http.StatusOK
		overall := "healthy"
		if dbStatus != "healthy" {
			status = http.StatusServiceUnavailable
			overall = "unhealthy"
		}

		c.JSON(status, gin.H{
			"status":   overall,
			"service":  "awards-management-api",
			"database": dbStatus,
		})
	})

	// API v1 路由组
	api := r.Group("/api/v1")
	{
		// 认证相关路由（限流：30 req/min per IP）
		auth := api.Group("/auth")
		auth.Use(middleware.IPRateLimit())
		{
			auth.POST("/login", authController.Login)
			auth.POST("/logout", authPkg.AuthMiddleware(), authController.Logout)
			auth.GET("/me", authPkg.AuthMiddleware(), authController.GetCurrentUserInfo)
			auth.PUT("/change-password", authPkg.AuthMiddleware(), authController.ChangePassword)
		}

		// 获奖记录管理相关路由
		awards := api.Group("/awards")
		{
			// 公开接口（限流：30 req/min per IP）
			publicAwards := awards.Use(middleware.IPRateLimit())
			{
				publicAwards.GET("/suggestions", awardController.SuggestNames)
				publicAwards.POST("/query/person", awardController.QueryPersonAwards)
				publicAwards.POST("/query/group", awardController.QueryGroupAwards)
				publicAwards.GET("/:award_id", awardController.GetAwardDetail)
			}

			// 需要认证的接口（限流：60 req/min per token）
			authAwards := awards.Use(authPkg.AuthMiddleware())
			authAwards.Use(middleware.TokenRateLimit())
			{
				authAwards.POST("/upload/certificate", awardController.UploadCertificateImage)
				authAwards.POST("/upload/certificates/batch", awardController.BatchUploadCertificateImages)
				authAwards.POST("/upload/person", awardController.UploadPersonAward)
				authAwards.POST("/upload/group", awardController.UploadGroupAward)
				authAwards.POST("/batch-upload/person", awardController.BatchUploadPersonAwards)
				authAwards.POST("/batch-upload/group", awardController.BatchUploadGroupAwards)
				authAwards.PUT("/update/person", awardController.UpdatePersonAward)
				authAwards.PUT("/update/group", awardController.UpdateGroupAward)
				authAwards.DELETE("/:award_id", awardController.DeleteAward)
				authAwards.GET("/manage/all", awardController.GetAllAwardsForManagement)
				authAwards.GET("/manage/all-light", awardController.GetAllAwardsForManagementLight)
				authAwards.GET("/group/:award_id/members", awardController.GetGroupMembersWithCertificates)
				authAwards.PUT("/group/member/certificate", awardController.UpdateMemberCertificate)
			}
		}

		// 监控相关路由
		monitor := api.Group("/monitor")
		{
			authMonitor := monitor.Use(authPkg.AuthMiddleware())
			{
				authMonitor.GET("/data", monitorController.GetMonitorData)
				authMonitor.GET("/database", monitorController.GetDatabaseStats)
				authMonitor.GET("/health", monitorController.GetSystemHealth)
				authMonitor.GET("/users", monitorController.GetOnlineUsers)
				authMonitor.GET("/requests", monitorController.GetRequestHistory)
				authMonitor.POST("/reset", monitorController.ResetStats)
				authMonitor.DELETE("/users/:user_id", monitorController.KickUser)
			}
		}
	}
}
