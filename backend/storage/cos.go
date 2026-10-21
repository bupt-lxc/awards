package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

// cosService 基于腾讯云 COS 的对象存储
type cosService struct {
	client        *cos.Client
	cfg           COSConfig
	localFallback *localFallback // COS 不可用时的本地回退
}

// NewCOSService 创建 COS 存储服务
func NewCOSService(cfg COSConfig) (Service, error) {
	if cfg.PublicURL == "" {
		cfg.PublicURL = "/static/certificates"
	}
	if cfg.SecretID == "" || cfg.SecretKey == "" || cfg.Region == "" || cfg.Bucket == "" {
		log.Printf("[Storage] COS 配置不完整，使用本地存储")
		return newLocalFallback("/static/certificates"), nil
	}

	bucketURL, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", cfg.Bucket, cfg.Region))
	if err != nil {
		log.Printf("[Storage] COS URL 解析失败: %v", err)
		return newLocalFallback("/static/certificates"), nil
	}

	client := cos.NewClient(&cos.BaseURL{BucketURL: bucketURL}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
		Timeout: 30 * time.Second,
	})

	// 验证连通性
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = client.Bucket.Head(ctx)
	if err != nil {
		log.Printf("[Storage] COS 连通性检查失败: %v，回退本地存储", err)
		return newLocalFallback("/static/certificates"), nil
	}

	log.Printf("[Storage] 腾讯云 COS 已连接: %s (region: %s)", cfg.Bucket, cfg.Region)
	return &cosService{
		client:        client,
		cfg:           cfg,
		localFallback: newLocalFallback("/static/certificates"),
	}, nil
}

func (s *cosService) Upload(fileName string, reader io.Reader, size int64, contentType string) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("读取文件内容失败: %w", err)
	}

	if s.client != nil {
		publicPath, err := s.uploadToCOS(fileName, data, size, contentType)
		if err == nil {
			return publicPath, nil
		}
		log.Printf("[Storage] COS 上传失败: %v，回退本地", err)
	}

	return s.localFallback.Upload(fileName, bytes.NewReader(data), int64(len(data)), contentType)
}

func (s *cosService) uploadToCOS(fileName string, data []byte, _ int64, contentType string) (string, error) {
	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = ".png"
	}
	objName := "certificates/" + uuid.New().String() + ext

	ctx := context.Background()
	opts := &cos.ObjectPutOptions{
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			XCosACL: "public-read",
		},
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType:   contentType,
			ContentLength: int64(len(data)),
		},
	}

	_, err := s.client.Object.Put(ctx, objName, bytes.NewReader(data), opts)
	if err != nil {
		return "", fmt.Errorf("COS PutObject 失败: %w", err)
	}

	return buildPublicPath(s.cfg.PublicURL, objName), nil
}

func (s *cosService) Get(filePath string) (io.ReadCloser, string, error) {
	// 优先从 COS 获取
	if s.client != nil {
		objName := extractObjectName(filePath)
		if objName != "" {
			ctx := context.Background()
			resp, err := s.client.Object.Get(ctx, objName, nil)
			if err == nil {
				return resp.Body, resp.Header.Get("Content-Type"), nil
			}
			log.Printf("[Storage] COS Get 失败: %v，尝试本地", err)
		}
	}
	// 回退本地
	return s.localFallback.Get(filePath)
}

func (s *cosService) Delete(filePath string) error {
	objName := extractObjectName(filePath)
	if objName == "" {
		return fmt.Errorf("无效的文件路径: %s", filePath)
	}

	if s.client != nil {
		ctx := context.Background()
		_, err := s.client.Object.Delete(ctx, objName)
		if err == nil {
			return nil
		}
		log.Printf("[Storage] COS 删除失败: %v，尝试本地删除", err)
	}

	return s.localFallback.Delete(filePath)
}

func (s *cosService) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := s.client.Bucket.Head(ctx)
	return err
}

// === 本地存储回退 ===

type localFallback struct {
	publicURL string
	dir       string
}

func newLocalFallback(publicURL string) *localFallback {
	return &localFallback{
		publicURL: publicURL,
		dir:       "uploads/certificates",
	}
}

func (s *localFallback) Upload(fileName string, reader io.Reader, _ int64, _ string) (string, error) {
	if err := os.MkdirAll(s.dir, 0750); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = ".png"
	}
	savedName := uuid.New().String() + ext
	savePath := filepath.Join(s.dir, savedName)

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}
	if err := os.WriteFile(savePath, data, 0640); err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}
	return s.publicURL + "/" + savedName, nil
}

func (s *localFallback) Get(filePath string) (io.ReadCloser, string, error) {
	fileName := filepath.Base(filePath)
	if fileName == "" || fileName == "." {
		return nil, "", fmt.Errorf("无效的文件路径: %s", filePath)
	}
	localPath := filepath.Join(s.dir, fileName)
	f, err := os.Open(localPath)
	if err != nil {
		return nil, "", err
	}
	return f, "image/png", nil
}

func (s *localFallback) Delete(filePath string) error {
	objName := filepath.Base(filePath)
	if objName == "" || objName == "." {
		return nil
	}
	localPath := filepath.Join(s.dir, objName)
	if err := os.Remove(localPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

func (s *localFallback) Ping() error {
	if _, err := os.Stat(s.dir); os.IsNotExist(err) {
		return fmt.Errorf("本地存储目录不存在: %s", s.dir)
	}
	return nil
}

// === 工具函数 ===

func extractObjectName(filePath string) string {
	// 从 URL 路径提取对象名
	// "/static/certificates/xxx.png" → "certificates/xxx.png"
	// "https://cdn.example.com/certificates/xxx.png" → "certificates/xxx.png"
	base := filepath.Base(filePath)
	if base == "" || base == "." {
		return ""
	}
	return "certificates/" + base
}

func buildPublicPath(publicURL, objectName string) string {
	base := strings.TrimRight(publicURL, "/")
	if base == "" || base == "/static/certificates" {
		return "/static/certificates/" + filepath.Base(objectName)
	}
	return base + "/" + strings.TrimLeft(objectName, "/")
}
