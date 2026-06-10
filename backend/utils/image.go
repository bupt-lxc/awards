package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/png"

	// 注册各图片格式的解码器（image.Decode 依赖这些匿名导入）
	_ "image/gif"
	_ "image/jpeg"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
)

// ConvertToPNG 将任意支持的图片格式转为 PNG 字节（无损，仅做格式转换不压缩）
func ConvertToPNG(input []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("图片解码失败: %w", err)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("PNG编码失败: %w", err)
	}

	return buf.Bytes(), nil
}

// AllowedMIMETypes 返回允许的图片 MIME 类型集合
func AllowedMIMETypes() map[string]bool {
	return map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
}

// MaxFileSize 返回最大允许的文件大小（10MB）
func MaxFileSize() int64 {
	return 10 << 20
}
