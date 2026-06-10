package auth

import (
	"awards-management/config"
	"awards-management/database"
	"awards-management/models"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Claims JWT声明结构
type Claims struct {
	Username string `json:"sub"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthManager 认证管理器
type AuthManager struct{}

// CreateAccessToken 创建访问令牌
func (am *AuthManager) CreateAccessToken(username, role string) (string, error) {
	claims := Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.AppConfig.JWT.ExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

// VerifyToken 验证令牌
func (am *AuthManager) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// AuthenticateUser 验证用户身份
func (am *AuthManager) AuthenticateUser(username, password string) (*models.Admin, error) {
	var admin models.Admin
	err := database.DB.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}

	// 使用 bcrypt 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &admin, nil
}

var authManager = &AuthManager{}

// GetAuthManager 获取认证管理器实例
func GetAuthManager() *AuthManager {
	return authManager
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "缺少认证头",
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的认证格式",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := authManager.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的认证凭据",
			})
			c.Abort()
			return
		}

		// 从数据库获取用户信息
		var admin models.Admin
		err = database.DB.Where("username = ?", claims.Username).First(&admin).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "用户不存在",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("current_user", &admin)
		c.Set("username", admin.Username)
		c.Set("user_id", fmt.Sprintf("%d", admin.ID)) // 添加用户ID用于监控
		c.Next()
	}
}

// RequireAdminRole 要求管理员角色中间件
func RequireAdminRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("current_user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未认证",
			})
			c.Abort()
			return
		}

		admin, ok := user.(*models.Admin)
		if !ok || admin.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "权限不足，需要管理员权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetCurrentUser 获取当前用户
func GetCurrentUser(c *gin.Context) (*models.Admin, bool) {
	user, exists := c.Get("current_user")
	if !exists {
		return nil, false
	}

	admin, ok := user.(*models.Admin)
	return admin, ok
}

// GetCurrentUsername 获取当前用户名
func GetCurrentUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}

	usernameStr, ok := username.(string)
	return usernameStr, ok
}
