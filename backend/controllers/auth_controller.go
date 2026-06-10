package controllers

import (
	"awards-management/auth"
	"awards-management/config"
	"awards-management/dto"
	"awards-management/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	adminService   *services.AdminService
	authManager    *auth.AuthManager
	monitorService *services.MonitorService
}

// NewAuthController 创建认证控制器实例
func NewAuthController() *AuthController {
	return &AuthController{
		adminService:   services.NewAdminService(),
		authManager:    auth.GetAuthManager(),
		monitorService: services.GetMonitorService(),
	}
}

// Login 管理员登录
func (ac *AuthController) Login(c *gin.Context) {
	var req dto.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证用户身份
	admin, err := ac.authManager.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "用户名或密码错误",
		})
		return
	}

	// 创建访问令牌
	accessToken, err := ac.authManager.CreateAccessToken(admin.Username, admin.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "创建令牌失败",
		})
		return
	}

	// 添加在线用户到监控系统
	userID := fmt.Sprintf("%d", admin.ID)
	ac.monitorService.AddOnlineUser(
		userID,
		admin.Username,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		admin.Role == "admin",
	)

	// 返回响应
	response := dto.TokenResponse{
		AccessToken: accessToken,
		TokenType:   "bearer",
		ExpiresIn:   config.AppConfig.JWT.ExpireHours * 3600, // 转换为秒
		UserInfo: map[string]interface{}{
			"username":   admin.Username,
			"role":       admin.Role,
			"created_at": admin.CreatedAt,
		},
	}

	c.JSON(http.StatusOK, response)
}

// ChangePassword 修改管理员密码
func (ac *AuthController) ChangePassword(c *gin.Context) {
	var req dto.AdminChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 从JWT中获取当前用户名
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "未找到用户信息",
		})
		return
	}

	// 修改密码
	err := ac.adminService.ChangePassword(username.(string), req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "密码修改成功",
	})
}

// GetCurrentUserInfo 获取当前用户信息
func (ac *AuthController) GetCurrentUserInfo(c *gin.Context) {
	currentUser, exists := auth.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "未认证",
		})
		return
	}

	userInfo := map[string]interface{}{
		"username":   currentUser.Username,
		"role":       currentUser.Role,
		"created_at": currentUser.CreatedAt,
	}

	c.JSON(http.StatusOK, userInfo)
}

// Logout 管理员登出
func (ac *AuthController) Logout(c *gin.Context) {
	// 从JWT中获取当前用户ID
	userID, exists := c.Get("user_id")
	if exists {
		ac.monitorService.RemoveOnlineUser(userID.(string))
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "登出成功",
	})
}
