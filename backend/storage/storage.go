package storage

import (
	"io"
	"sync"
)

// Service 对象存储服务接口
type Service interface {
	Upload(fileName string, reader io.Reader, size int64, contentType string) (string, error)
	Delete(filePath string) error
	Ping() error
	Get(filePath string) (io.ReadCloser, string, error) // 获取文件内容和 Content-Type
}

// COSConfig 腾讯云 COS 配置
type COSConfig struct {
	SecretID  string // 腾讯云 SecretId
	SecretKey string // 腾讯云 SecretKey
	Region    string // 地域，如 ap-guangzhou
	Bucket    string // 存储桶名称，格式 bucket-appid
	PublicURL string // 对外访问 URL（CDN 加速域名或 COS 默认域名）
}

var (
	globalSvc Service
	mu        sync.Mutex
)

// SetGlobal 设置全局存储服务
func SetGlobal(svc Service) {
	mu.Lock()
	defer mu.Unlock()
	globalSvc = svc
}

// Get 获取全局存储服务
func Get() Service {
	return globalSvc
}
