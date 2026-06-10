package controllers

import (
	"awards-management/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// serverError 根据运行模式返回错误：debug 模式返回详细信息，release 模式返回通用消息
func serverError(c *gin.Context, prefix string, err error) {
	msg := "服务器内部错误，请稍后重试"
	if config.AppConfig != nil && config.AppConfig.Server.Mode == "debug" {
		msg = prefix + ": " + err.Error()
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
}
