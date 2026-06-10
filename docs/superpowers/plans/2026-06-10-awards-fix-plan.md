# 获奖管理系统 — 安全与存储修复 实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 修复 23 个已发现问题中的代码层面问题，补齐 Docker/nginx/systemd 运维基础设施，使项目达到可上线标准。

**Architecture:** 分层修复 — 底层依赖（bcrypt、工具函数）→ 数据模型 → 配置 → 认证 → Service → Controller → 中间件 → 运维文件。每层修复后上层自然生效。

**Tech Stack:** Go 1.21, Gin, GORM, bcrypt, MySQL 8.0, Docker, Nginx

---

## ⚠️ 执行前必读：关键约定

1. **文件替换 vs 精确编辑**：标记"用以下内容完全替换"的步骤 → 使用 Write 工具覆盖整个文件。未标记"完全替换"的步骤 → 使用 Edit 工具做精确修改。
2. **中间态编译失败是预期的**：Task 3 删除 `Members` 字段后、Task 8/9/10 执行前，项目会有编译错误。每个 Task 正常 commit，后续 Task 会修复。
3. **go mod tidy 时机**：每个添加/修改依赖的 Task 末尾必须执行 `go mod tidy`。
4. **import 整理**：替换文件时使用给定 import；编辑文件时确保移除不再使用的 import、添加新使用的 import。
5. **路径基准**：所有路径相对于仓库根目录 `d:/projects/奖状查询/奖状/`。

---

## 文件变更地图

| 操作 | 文件 |
|------|------|
| **新增** | `backend/utils/image.go`, `backend/middleware/rate_limit.go`, `backend/Dockerfile`, `backend/docker-compose.yml`, `backend/nginx.conf`, `backend/.env.example` |
| **修改** | `backend/go.mod`, `backend/models/models.go`, `backend/config/config.go`, `backend/auth/auth.go`, `backend/database/database.go`, `backend/services/admin_service.go`, `backend/services/person_award_service.go`, `backend/services/group_award_service.go`, `backend/controllers/award_controller.go`, `backend/routes/routes.go`, `backend/main.go`, `backend/.gitignore` |
| **不变** | 前端全部文件、`backend/dto/dto.go`、`backend/middleware/monitor_middleware.go`、`backend/controllers/monitor_controller.go`、`backend/services/monitor_service.go`、`backend/config.yaml` |

---

### Task 1: 添加 bcrypt 依赖

**Files:**
- Modify: `backend/go.mod`

- [ ] **Step 1: 添加 golang.org/x/crypto 依赖**

`backend/go.mod` 的 require 块中已有 `golang.org/x/crypto v0.17.0 // indirect`，需要将其改为 direct 依赖。在 require 块中添加：

```
golang.org/x/crypto v0.17.0
```

同时从 indirect 块中移除对应行。

- [ ] **Step 2: 运行 go mod tidy 确认依赖正确**

```bash
cd backend && go mod tidy
```

Expected: 无错误输出。

- [ ] **Step 3: Commit**

```bash
git add backend/go.mod backend/go.sum
git commit -m "deps: add bcrypt (golang.org/x/crypto) as direct dependency"
```

---

### Task 2: 新增图片转 PNG 工具函数

**Files:**
- Create: `backend/utils/image.go`

- [ ] **Step 1: 创建 utils 目录和 image.go**

```bash
mkdir -p backend/utils
```

- [ ] **Step 2: 编写 image.go**

```go
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
```

- [ ] **Step 3: 添加 golang.org/x/image 依赖（webp + bmp 解码器）**

```bash
cd backend && go get golang.org/x/image@latest && go mod tidy
```

> **说明**：`golang.org/x/image` 包含 `webp` 和 `bmp` 两个子包的解码器注册。匿名导入后 `image.Decode` 即可识别 WebP 和 BMP 格式。

- [ ] **Step 4: 编译验证**

```bash
cd backend && go build ./...
```

Expected: 无编译错误。

- [ ] **Step 5: Commit**

```bash
git add backend/utils/ backend/go.mod backend/go.sum
git commit -m "feat: add image-to-PNG conversion utility"
```

---

### Task 3: 更新数据模型 — 字段类型 + 移除冗余

**Files:**
- Modify: `backend/models/models.go`

- [ ] **Step 1: 使用 Edit 工具对 models.go 做三项精确修改**

> ⚠️ **使用 Edit 工具逐项修改，不要 Write 覆盖整个文件**（文件还包含 `Member` 结构体、`TableName()` 方法、`AutoMigrate()` 函数，不可丢失）。

**修改 A — PersonAward.CertificateImage 字段 gorm tag：**

将 `CertificateImage string \`json:"certificate_image" gorm:"type:longtext"\`` 改为 `CertificateImage string \`json:"certificate_image" gorm:"size:500"\``

**修改 B — GroupAward 结构体：删除 Members 字段 + CertificateImage gorm tag：**

- 删除 `Members string \`json:"members" gorm:"type:longtext;not null"\`  // JSON格式存储成员信息` 这一整行及其上方的注释行
- 将 `CertificateImage string \`json:"certificate_image" gorm:"type:longtext"\`` 改为 `CertificateImage string \`json:"certificate_image" gorm:"size:500"\``

**修改 C — GroupMember.CertificateImage 字段 gorm tag：**

将 `CertificateImage string \`json:"certificate_image" gorm:"column:certificate_image"\`` 改为 `CertificateImage string \`json:"certificate_image" gorm:"size:500"\``

- [ ] **Step 2: 编译验证**

```bash
cd backend && go build ./...
```

Expected: 编译失败（因为 service 层引用了已删除的 `Members` 字段）。这是预期的——后续 task 修复。

- [ ] **Step 3: Commit**

```bash
git add backend/models/models.go
git commit -m "refactor: change certificate_image to file path, remove redundant members JSON field"
```

---

### Task 4: 配置改为环境变量优先

**Files:**
- Modify: `backend/config/config.go`

- [ ] **Step 1: 重写 config.go，添加环境变量绑定**

用以下内容完全替换 `backend/config/config.go`：

```go
package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	Charset         string        `mapstructure:"charset"`
	ParseTime       bool          `mapstructure:"parseTime"`
	Loc             string        `mapstructure:"loc"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expireHours"`
}

type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allowOrigins"`
	AllowMethods     []string `mapstructure:"allowMethods"`
	AllowHeaders     []string `mapstructure:"allowHeaders"`
	AllowCredentials bool     `mapstructure:"allowCredentials"`
}

var AppConfig *Config

func LoadConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置默认值
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8000)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.parseTime", true)
	viper.SetDefault("database.loc", "Local")
	viper.SetDefault("database.maxIdleConns", 10)
	viper.SetDefault("database.maxOpenConns", 100)
	viper.SetDefault("database.connMaxLifetime", 3600)
	viper.SetDefault("jwt.expireHours", 8)
	viper.SetDefault("cors.allowCredentials", true)

	// 环境变量绑定 — 优先级高于 config.yaml
	// Server
	_ = viper.BindEnv("server.mode", "GIN_MODE")
	// Database
	_ = viper.BindEnv("database.host", "DB_HOST")
	_ = viper.BindEnv("database.port", "DB_PORT")
	_ = viper.BindEnv("database.username", "DB_USER")
	_ = viper.BindEnv("database.password", "DB_PASSWORD")
	_ = viper.BindEnv("database.dbname", "DB_NAME")
	// JWT
	_ = viper.BindEnv("jwt.secret", "JWT_SECRET")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: 无法读取配置文件 %s: %v", configPath, err)
		log.Println("将仅使用环境变量和默认值")
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Printf("Error unmarshaling config: %v", err)
		return err
	}

	AppConfig.Database.ConnMaxLifetime *= time.Second

	log.Println("配置文件加载成功")
	return nil
}
```

- [ ] **Step 2: 编译验证**

```bash
cd backend && go build ./...
```

Expected: 无编译错误（config 包不依赖其他被改动的包）。

- [ ] **Step 3: Commit**

```bash
git add backend/config/config.go
git commit -m "feat: add environment variable overrides for sensitive config"
```

---

### Task 5: 认证模块 SHA-256 → bcrypt

**Files:**
- Modify: `backend/auth/auth.go`

- [ ] **Step 1: 修改 auth.go 的 AuthenticateUser 方法**

文件头部添加 import：

```go
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
```

修改 `AuthenticateUser` 方法：

```go
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
```

- [ ] **Step 2: 编译验证**

```bash
cd backend && go build ./...
```

- [ ] **Step 3: Commit**

```bash
git add backend/auth/auth.go
git commit -m "security: replace SHA-256 with bcrypt for password verification"
```

---

### Task 6: 数据库初始化改为 bcrypt

**Files:**
- Modify: `backend/database/database.go`

- [ ] **Step 1: 修改 createDefaultAdmin 使用 bcrypt**

在 `backend/database/database.go` 顶部添加 import：

```go
import (
	"awards-management/config"
	"awards-management/models"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)
```

移除 `crypto/sha256` import（不再需要）。

重写 `createDefaultAdmin` 函数：

```go
// createDefaultAdmin 创建默认管理员账户
func createDefaultAdmin() error {
	var count int64
	DB.Model(&models.Admin{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		return nil
	}

	// 使用 bcrypt 生成密码哈希 (cost=12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), 12)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	admin := models.Admin{
		Username:     "admin",
		PasswordHash: string(hashedPassword),
		Role:         "admin",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return DB.Create(&admin).Error
}
```

- [ ] **Step 2: 编译验证**

```bash
cd backend && go build ./...
```

Expected: 无编译错误。

- [ ] **Step 3: Commit**

```bash
git add backend/database/database.go
git commit -m "security: use bcrypt for default admin password hashing"
```

---

### Task 7: Admin Service 改为 bcrypt

**Files:**
- Modify: `backend/services/admin_service.go`

- [ ] **Step 1: 重写 admin_service.go**

用以下内容完全替换：

```go
package services

import (
	"awards-management/database"
	"awards-management/models"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminService struct{}

func NewAdminService() *AdminService {
	return &AdminService{}
}

func (s *AdminService) ChangePassword(username, oldPassword, newPassword string) error {
	_, err := s.AuthenticateAdmin(username, oldPassword)
	if err != nil {
		return fmt.Errorf("旧密码错误")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	err = database.DB.Model(&models.Admin{}).Where("username = ?", username).
		Update("password_hash", string(hashedPassword)).Error
	if err != nil {
		return fmt.Errorf("密码更新失败: %v", err)
	}

	return nil
}

func (s *AdminService) GetAdminByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	err := database.DB.Where("username = ?", username).First(&admin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}
	return &admin, nil
}

func (s *AdminService) AuthenticateAdmin(username, password string) (*models.Admin, error) {
	var admin models.Admin
	err := database.DB.Where("username = ?", username).First(&admin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户名或密码错误")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	return &admin, nil
}
```

- [ ] **Step 2: 编译验证**

```bash
cd backend && go build ./...
```

Expected: 无编译错误。

- [ ] **Step 3: Commit**

```bash
git add backend/services/admin_service.go
git commit -m "security: use bcrypt in admin service for password operations"
```

---

### Task 8: 重写 PersonAwardService — 移除 base64 + 简化批量上传

**Files:**
- Modify: `backend/services/person_award_service.go`

- [ ] **Step 1: 用以下内容完全替换 person_award_service.go**

```go
package services

import (
	"awards-management/database"
	"awards-management/dto"
	"awards-management/models"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type PersonAwardService struct{}

func NewPersonAwardService() *PersonAwardService {
	return &PersonAwardService{}
}

func (s *PersonAwardService) QueryPersonAwards(query dto.PersonAwardQuery) ([]models.PersonAward, error) {
	var awards []models.PersonAward
	db := database.DB.Model(&models.PersonAward{})

	db = db.Where("person_name = ? AND id_number = ?", query.PersonName, query.IDNumber)

	if query.Year != nil {
		db = db.Where("year = ?", *query.Year)
	}
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}

	err := db.Find(&awards).Error
	return awards, err
}

func (s *PersonAwardService) CreatePersonAward(record dto.PersonAwardRecord) (*models.PersonAward, error) {
	award := models.PersonAward{
		Year:             record.Year,
		PersonName:       record.PersonName,
		IDNumber:         record.IDNumber,
		Category:         record.Category,
		AwardLevel:       record.AwardLevel,
		Description:      record.Description,
		CertificateImage: record.CertificateImage,
	}

	err := database.DB.Create(&award).Error
	if err != nil {
		return nil, err
	}

	return &award, nil
}

func (s *PersonAwardService) BatchCreatePersonAwards(records []dto.PersonAwardRecord) (int, int, []int, error) {
	successful := 0
	failed := 0
	var failedItems []int

	for i, record := range records {
		_, err := s.CreatePersonAward(record)
		if err != nil {
			failed++
			failedItems = append(failedItems, i)
			log.Printf("创建个人记录失败 (ID: %s): %v", record.IDNumber, err)
		} else {
			successful++
		}
	}

	return successful, failed, failedItems, nil
}

func (s *PersonAwardService) UpdatePersonAward(update dto.PersonAwardUpdate) error {
	var award models.PersonAward
	err := database.DB.First(&award, update.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("获奖记录不存在")
		}
		return err
	}

	updateData := make(map[string]interface{})
	if update.Year != nil {
		updateData["year"] = *update.Year
	}
	if update.PersonName != "" {
		updateData["person_name"] = update.PersonName
	}
	if update.IDNumber != "" {
		updateData["id_number"] = update.IDNumber
	}
	if update.Category != "" {
		updateData["category"] = update.Category
	}
	if update.AwardLevel != "" {
		updateData["award_level"] = update.AwardLevel
	}
	if update.Description != "" {
		updateData["description"] = update.Description
	}
	if update.CertificateImage != "" {
		updateData["certificate_image"] = update.CertificateImage
	}

	if len(updateData) == 0 {
		return fmt.Errorf("没有提供更新数据")
	}

	return database.DB.Model(&award).Updates(updateData).Error
}

func (s *PersonAwardService) DeletePersonAward(id uint) error {
	result := database.DB.Delete(&models.PersonAward{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("获奖记录不存在")
	}
	return nil
}

func (s *PersonAwardService) GetPersonAwardByID(id uint) (*models.PersonAward, error) {
	var award models.PersonAward
	err := database.DB.First(&award, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("获奖记录不存在")
		}
		return nil, err
	}
	return &award, nil
}

func (s *PersonAwardService) GetAllPersonAwards() ([]models.PersonAward, error) {
	var awards []models.PersonAward
	err := database.DB.Order("created_at DESC").Find(&awards).Error
	return awards, err
}

func (s *PersonAwardService) GetAllPersonAwardsLight() ([]map[string]interface{}, error) {
	var awards []models.PersonAward
	err := database.DB.Select("id, year, person_name, id_number, category, award_level, description, created_at, updated_at").
		Order("created_at DESC").
		Find(&awards).Error
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, award := range awards {
		results = append(results, map[string]interface{}{
			"id":              award.ID,
			"year":            award.Year,
			"person_name":     award.PersonName,
			"id_number":       award.IDNumber,
			"category":        award.Category,
			"award_level":     award.AwardLevel,
			"description":     award.Description,
			"created_at":      award.CreatedAt,
			"updated_at":      award.UpdatedAt,
			"has_certificate": award.CertificateImage != "",
		})
	}

	return results, nil
}
```

- [ ] **Step 2: 编译验证**

```bash
cd backend && go build ./...
```

Expected: 可能仍有 group_award_service 相关的编译错误（Task 9 修复）。

- [ ] **Step 3: Commit**

```bash
git add backend/services/person_award_service.go
git commit -m "refactor: remove base64 logic from person award service"
```

---

### Task 9: 重写 GroupAwardService — 移除 base64 + Members 冗余 + FOR UPDATE + 清理日志

**Files:**
- Modify: `backend/services/group_award_service.go`

- [ ] **Step 1: 用以下精简版本完全替换 group_award_service.go**

```go
package services

import (
	"awards-management/database"
	"awards-management/dto"
	"awards-management/models"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

type GroupAwardService struct{}

func NewGroupAwardService() *GroupAwardService {
	return &GroupAwardService{}
}

func (s *GroupAwardService) QueryGroupAwards(query dto.GroupAwardQuery) ([]map[string]interface{}, error) {
	var awards []models.GroupAward

	dbQuery := database.DB.Model(&models.GroupAward{}).
		Preload("GroupMembers").
		Order("created_at DESC")

	dbQuery = dbQuery.Where("group_name LIKE ?", "%"+query.GroupName+"%")

	if query.Year != nil {
		dbQuery = dbQuery.Where("year = ?", *query.Year)
	}
	if query.Category != "" {
		dbQuery = dbQuery.Where("category = ?", query.Category)
	}
	if query.PersonName != "" && query.IDNumber != "" {
		dbQuery = dbQuery.Where("EXISTS (SELECT 1 FROM group_members WHERE group_award_id = group_awards.id AND member_name = ? AND member_id_number = ?)",
			query.PersonName, query.IDNumber)
	}

	if err := dbQuery.Find(&awards).Error; err != nil {
		return nil, fmt.Errorf("查询团队信息失败: %v", err)
	}

	log.Printf("查询到 %d 条团队记录", len(awards))

	resultAwards := make([]map[string]interface{}, 0, len(awards))

	for _, award := range awards {
		certificateImage := strings.TrimSpace(award.CertificateImage)

		awardMap := map[string]interface{}{
			"id":                award.ID,
			"year":              award.Year,
			"group_name":        award.GroupName,
			"category":          award.Category,
			"award_level":       award.AwardLevel,
			"description":       award.Description,
			"certificate_image": certificateImage,
			"created_at":        award.CreatedAt,
			"updated_at":        award.UpdatedAt,
			"members":           []map[string]interface{}{},
		}

		for _, member := range award.GroupMembers {
			if (query.PersonName == "" || query.PersonName == member.MemberName) &&
				(query.IDNumber == "" || query.IDNumber == member.MemberIDNumber) {

				memberCert := strings.TrimSpace(member.CertificateImage)
				if memberCert == "" {
					memberCert = certificateImage
				}

				memberMap := map[string]interface{}{
					"name":              member.MemberName,
					"id_number":         member.MemberIDNumber,
					"certificate_image": memberCert,
				}
				awardMap["members"] = append(awardMap["members"].([]map[string]interface{}), memberMap)
			}
		}

		resultAwards = append(resultAwards, awardMap)
	}

	return resultAwards, nil
}

func (s *GroupAwardService) CreateGroupAward(record dto.GroupAwardRecord) (*models.GroupAward, error) {
	award := models.GroupAward{
		Year:             record.Year,
		GroupName:        record.GroupName,
		Category:         record.Category,
		AwardLevel:       record.AwardLevel,
		Description:      record.Description,
		CertificateImage: record.CertificateImage,
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&award).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, member := range record.Members {
		certImage := member.CertificateImage
		if certImage == "" {
			certImage = record.CertificateImage
		}

		groupMember := models.GroupMember{
			GroupAwardID:     award.ID,
			MemberName:       member.Name,
			MemberIDNumber:   member.IDNumber,
			CertificateImage: certImage,
		}
		if err := tx.Create(&groupMember).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	log.Printf("创建团体记录成功: %s", record.GroupName)
	return &award, nil
}

func (s *GroupAwardService) BatchCreateGroupAwards(records []dto.GroupAwardRecord) (int, int, []int, error) {
	successful := 0
	failed := 0
	var failedItems []int

	for i, record := range records {
		_, err := s.CreateGroupAward(record)
		if err != nil {
			log.Printf("创建团体记录失败 (团队: %s): %v", record.GroupName, err)
			failed++
			failedItems = append(failedItems, i)
		} else {
			successful++
		}
	}
	return successful, failed, failedItems, nil
}

func (s *GroupAwardService) UpdateGroupAward(update dto.GroupAwardUpdate) error {
	var award models.GroupAward
	err := database.DB.First(&award, update.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("获奖记录不存在")
		}
		return err
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updateData := make(map[string]interface{})
	if update.Year != nil {
		updateData["year"] = *update.Year
	}
	if update.GroupName != "" {
		updateData["group_name"] = update.GroupName
	}
	if update.Category != "" {
		updateData["category"] = update.Category
	}
	if update.AwardLevel != "" {
		updateData["award_level"] = update.AwardLevel
	}
	if update.Description != "" {
		updateData["description"] = update.Description
	}
	if update.CertificateImage != "" {
		updateData["certificate_image"] = update.CertificateImage
	}

	if update.Members != nil {
		for i := range update.Members {
			if update.Members[i].CertificateImage == "" {
				update.Members[i].CertificateImage = award.CertificateImage
			}
		}

		if err := tx.Where("group_award_id = ?", update.ID).Delete(&models.GroupMember{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		for _, member := range update.Members {
			gm := models.GroupMember{
				GroupAwardID:     update.ID,
				MemberName:       member.Name,
				MemberIDNumber:   member.IDNumber,
				CertificateImage: member.CertificateImage,
			}
			if err := tx.Create(&gm).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if len(updateData) > 0 {
		if err := tx.Model(&award).Updates(updateData).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (s *GroupAwardService) DeleteGroupAward(id uint) error {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var award models.GroupAward
	if err := tx.First(&award, id).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("获奖记录不存在")
		}
		return err
	}

	var members []models.GroupMember
	if err := tx.Where("group_award_id = ?", id).Find(&members).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 收集需要删除的文件
	var filesToDelete []string
	if award.CertificateImage != "" && strings.HasPrefix(award.CertificateImage, "/static/certificates/") {
		filesToDelete = append(filesToDelete, strings.TrimPrefix(award.CertificateImage, "/static/certificates/"))
	}
	for _, m := range members {
		if m.CertificateImage != "" && strings.HasPrefix(m.CertificateImage, "/static/certificates/") {
			filesToDelete = append(filesToDelete, strings.TrimPrefix(m.CertificateImage, "/static/certificates/"))
		}
	}

	if err := tx.Exec("DELETE FROM group_members WHERE group_award_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&award).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	// 异步清理文件
	go func(files []string) {
		for _, file := range files {
			filePath := filepath.Join("uploads/certificates", file)
			if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
				log.Printf("清理证书文件失败 [%s]: %v", filePath, err)
			}
		}
	}(filesToDelete)

	log.Printf("删除团体记录成功: ID=%d", id)
	return nil
}

func (s *GroupAwardService) GetGroupAwardByID(id uint) (map[string]interface{}, error) {
	var award models.GroupAward
	err := database.DB.Preload("GroupMembers").First(&award, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("获奖记录不存在")
		}
		return nil, err
	}

	result := map[string]interface{}{
		"id":                award.ID,
		"year":              award.Year,
		"group_name":        award.GroupName,
		"category":          award.Category,
		"award_level":       award.AwardLevel,
		"description":       award.Description,
		"certificate_image": award.CertificateImage,
		"created_at":        award.CreatedAt,
		"updated_at":        award.UpdatedAt,
		"members":           []map[string]interface{}{},
	}

	for _, member := range award.GroupMembers {
		result["members"] = append(result["members"].([]map[string]interface{}), map[string]interface{}{
			"name":              member.MemberName,
			"id_number":         member.MemberIDNumber,
			"certificate_image": member.CertificateImage,
		})
	}

	return result, nil
}

func (s *GroupAwardService) GetAllGroupAwards() ([]map[string]interface{}, error) {
	var awards []models.GroupAward

	err := database.DB.Preload("GroupMembers").
		Order("created_at DESC").
		Find(&awards).Error

	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	var results []map[string]interface{}
	for _, award := range awards {
		if award.GroupName == "" || award.Year == 0 {
			continue
		}

		result := map[string]interface{}{
			"id":                award.ID,
			"year":              award.Year,
			"group_name":        award.GroupName,
			"category":          award.Category,
			"award_level":       award.AwardLevel,
			"description":       award.Description,
			"certificate_image": award.CertificateImage,
			"created_at":        award.CreatedAt,
			"updated_at":        award.UpdatedAt,
			"members":           []map[string]interface{}{},
		}

		for _, member := range award.GroupMembers {
			if member.MemberName == "" || member.MemberIDNumber == "" {
				continue
			}
			result["members"] = append(result["members"].([]map[string]interface{}), map[string]interface{}{
				"name":              member.MemberName,
				"id_number":         member.MemberIDNumber,
				"certificate_image": member.CertificateImage,
			})
		}

		results = append(results, result)
	}

	return results, nil
}

func (s *GroupAwardService) GetAllGroupAwardsLight() ([]map[string]interface{}, error) {
	var awards []models.GroupAward

	err := database.DB.Select("id, year, group_name, category, award_level, description, created_at, updated_at").
		Preload("GroupMembers").
		Order("created_at DESC").
		Find(&awards).Error

	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	var results []map[string]interface{}
	for _, award := range awards {
		if award.GroupName == "" || award.Year == 0 {
			continue
		}

		result := map[string]interface{}{
			"id":              award.ID,
			"year":            award.Year,
			"group_name":      award.GroupName,
			"category":        award.Category,
			"award_level":     award.AwardLevel,
			"description":     award.Description,
			"created_at":      award.CreatedAt,
			"updated_at":      award.UpdatedAt,
			"has_certificate": award.CertificateImage != "",
			"members":         []map[string]interface{}{},
		}

		for _, member := range award.GroupMembers {
			if member.MemberName == "" || member.MemberIDNumber == "" {
				continue
			}
			result["members"] = append(result["members"].([]map[string]interface{}), map[string]interface{}{
				"name":            member.MemberName,
				"id_number":       member.MemberIDNumber,
				"has_certificate": member.CertificateImage != "",
			})
		}

		results = append(results, result)
	}

	return results, nil
}

func (s *GroupAwardService) GetGroupMembersWithCertificates(awardID uint) ([]models.GroupMember, error) {
	var members []models.GroupMember
	err := database.DB.Where("group_award_id = ?", awardID).Find(&members).Error
	return members, err
}

func (s *GroupAwardService) UpdateMemberCertificate(update dto.MemberCertificateUpdate) error {
	result := database.DB.Model(&models.GroupMember{}).
		Where("group_award_id = ? AND member_name = ? AND member_id_number = ?",
			update.AwardID, update.MemberName, update.MemberIDNumber).
		Update("certificate_image", update.CertificateImage)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到匹配的成员记录")
	}

	return nil
}

func (s *GroupAwardService) UpdateGroupCertificate(awardID uint, certificateImage string) error {
	var award models.GroupAward
	if err := database.DB.First(&award, awardID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("未找到团队记录")
		}
		return err
	}

	award.CertificateImage = certificateImage
	if err := database.DB.Save(&award).Error; err != nil {
		return fmt.Errorf("更新证书失败: %v", err)
	}

	log.Printf("更新团队证书成功: ID=%d", award.ID)
	return nil
}

func (s *GroupAwardService) UpdateGroupCertificateByName(groupName string, certificateImage string) error {
	var award models.GroupAward
	if err := database.DB.Where("group_name = ?", groupName).First(&award).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("未找到团队记录: %s", groupName)
		}
		return err
	}

	return s.UpdateGroupCertificate(award.ID, certificateImage)
}
```

- [ ] **Step 2: 编译验证**

```bash
cd backend && go build ./...
```

Expected: 编译错误（award_controller 引用的旧 base64 逻辑尚未被 Task 10 替换）。这是预期的，Task 10 将修复。

- [ ] **Step 3: Commit（尽管编译可能暂时失败，Task 10 会修复）**

```bash
git add backend/services/group_award_service.go
git commit -m "refactor: remove base64 logic, redundant members, FOR UPDATE lock, debug logs from group service"
```

---

### Task 10: 重写 AwardController — 文件上传校验 + PNG 转换

**Files:**
- Modify: `backend/controllers/award_controller.go`

- [ ] **Step 1: 用以下内容完全替换 award_controller.go**

```go
package controllers

import (
	"awards-management/auth"
	"awards-management/dto"
	"awards-management/services"
	"awards-management/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AwardController struct {
	personAwardService *services.PersonAwardService
	groupAwardService  *services.GroupAwardService
}

func NewAwardController() *AwardController {
	return &AwardController{
		personAwardService: services.NewPersonAwardService(),
		groupAwardService:  services.NewGroupAwardService(),
	}
}

func (ac *AwardController) QueryPersonAwards(c *gin.Context) {
	var query dto.PersonAwardQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "请求参数错误: " + err.Error()})
		return
	}

	awards, err := ac.personAwardService.QueryPersonAwards(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "查询失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.QueryResponse{
		Awards: awards,
		Total:  len(awards),
		QueryParams: map[string]interface{}{
			"year": query.Year, "person_name": query.PersonName,
			"id_number": query.IDNumber, "category": query.Category,
		},
	})
}

func (ac *AwardController) QueryGroupAwards(c *gin.Context) {
	var query dto.GroupAwardQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "请求参数错误: " + err.Error()})
		return
	}

	if query.Year == nil && query.PersonName == "" && query.IDNumber == "" && query.GroupName == "" && query.Category == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "请至少提供一个查询条件"})
		return
	}

	awards, err := ac.groupAwardService.QueryGroupAwards(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "查询失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.QueryResponse{
		Awards: awards,
		Total:  len(awards),
		QueryParams: map[string]interface{}{
			"person_name": query.PersonName, "id_number": query.IDNumber,
			"year": query.Year, "group_name": query.GroupName, "category": query.Category,
		},
	})
}

func (ac *AwardController) UploadPersonAward(c *gin.Context) {
	var record dto.PersonAwardRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "请求参数错误: " + err.Error()})
		return
	}

	award, err := ac.personAwardService.CreatePersonAward(record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "个人获奖记录上传失败: " + err.Error()})
		return
	}

	username, _ := auth.GetCurrentUsername(c)
	c.JSON(http.StatusOK, dto.UploadResponse{
		Message: "个人获奖记录上传成功", AwardID: award.ID, UploadedBy: username,
	})
}

func (ac *AwardController) UploadGroupAward(c *gin.Context) {
	var record dto.GroupAwardRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "请求参数错误: " + err.Error()})
		return
	}

	award, err := ac.groupAwardService.CreateGroupAward(record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "团体获奖记录上传失败: " + err.Error()})
		return
	}

	username, _ := auth.GetCurrentUsername(c)
	c.JSON(http.StatusOK, dto.UploadResponse{
		Message: "团体获奖记录上传成功", AwardID: award.ID, UploadedBy: username,
	})
}

func (ac *AwardController) BatchUploadPersonAwards(c *gin.Context) {
	var records []dto.PersonAwardRecord
	if err := c.ShouldBindJSON(&records); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "请求参数错误: " + err.Error()})
		return
	}

	successful, failed, failedItems, err := ac.personAwardService.BatchCreatePersonAwards(records)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "批量上传失败: " + err.Error()})
		return
	}

	username, _ := auth.GetCurrentUsername(c)
	c.JSON(http.StatusOK, dto.BatchUploadResponse{
		Message: "批量上传完成", Successful: successful, Failed: failed,
		UploadedBy: username, FailedItems: failedItems,
	})
}

func (ac *AwardController) BatchUploadGroupAwards(c *gin.Context) {
	var records []dto.GroupAwardRecord
	if err := c.ShouldBindJSON(&records); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "请求参数错误: " + err.Error()})
		return
	}

	successful, failed, failedItems, err := ac.groupAwardService.BatchCreateGroupAwards(records)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "批量上传失败: " + err.Error()})
		return
	}

	username, _ := auth.GetCurrentUsername(c)
	c.JSON(http.StatusOK, dto.BatchUploadResponse{
		Message: "批量上传完成", Successful: successful, Failed: failed,
		UploadedBy: username, FailedItems: failedItems,
	})
}

func (ac *AwardController) UpdatePersonAward(c *gin.Context) {
	var update dto.PersonAwardUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "请求参数错误: " + err.Error()})
		return
	}

	if err := ac.personAwardService.UpdatePersonAward(update); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "更新失败: " + err.Error()})
		return
	}

	username, _ := auth.GetCurrentUsername(c)
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "个人获奖记录更新成功",
		Data:    map[string]interface{}{"updated_by": username, "award_id": update.ID},
	})
}

func (ac *AwardController) UpdateGroupAward(c *gin.Context) {
	var update dto.GroupAwardUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "请求参数错误: " + err.Error()})
		return
	}

	if err := ac.groupAwardService.UpdateGroupAward(update); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "更新失败: " + err.Error()})
		return
	}

	username, _ := auth.GetCurrentUsername(c)
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "团体获奖记录更新成功",
		Data:    map[string]interface{}{"updated_by": username, "award_id": update.ID},
	})
}

func (ac *AwardController) DeleteAward(c *gin.Context) {
	awardID, err := strconv.ParseUint(c.Param("award_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "无效的获奖记录ID"})
		return
	}

	awardType := c.Query("award_type")
	if awardType == "" {
		awardType = "person"
	}

	var deleteErr error
	if awardType == "group" {
		deleteErr = ac.groupAwardService.DeleteGroupAward(uint(awardID))
	} else {
		deleteErr = ac.personAwardService.DeletePersonAward(uint(awardID))
	}

	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "删除失败: " + deleteErr.Error()})
		return
	}

	username, _ := auth.GetCurrentUsername(c)
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "获奖记录删除成功",
		Data:    map[string]interface{}{"deleted_by": username, "award_id": awardID, "award_type": awardType},
	})
}

func (ac *AwardController) GetAwardDetail(c *gin.Context) {
	awardID, err := strconv.ParseUint(c.Param("award_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "无效的获奖记录ID"})
		return
	}

	personAward, err := ac.personAwardService.GetPersonAwardByID(uint(awardID))
	if err == nil {
		c.JSON(http.StatusOK, map[string]interface{}{"type": "person", "award": personAward})
		return
	}

	groupAward, err := ac.groupAwardService.GetGroupAwardByID(uint(awardID))
	if err == nil {
		c.JSON(http.StatusOK, map[string]interface{}{"type": "group", "award": groupAward})
		return
	}

	c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "获奖记录不存在"})
}

func (ac *AwardController) GetAllAwardsForManagement(c *gin.Context) {
	personAwards, err := ac.personAwardService.GetAllPersonAwards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "获取个人获奖记录失败: " + err.Error()})
		return
	}

	groupAwards, err := ac.groupAwardService.GetAllGroupAwards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "获取团体获奖记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"person_awards": personAwards, "group_awards": groupAwards,
		"total_person": len(personAwards), "total_group": len(groupAwards),
		"total": len(personAwards) + len(groupAwards),
	})
}

func (ac *AwardController) GetAllAwardsForManagementLight(c *gin.Context) {
	personAwards, err := ac.personAwardService.GetAllPersonAwardsLight()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "获取个人获奖记录失败: " + err.Error()})
		return
	}

	groupAwards, err := ac.groupAwardService.GetAllGroupAwardsLight()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "获取团体获奖记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"person_awards": personAwards, "group_awards": groupAwards,
		"total_person": len(personAwards), "total_group": len(groupAwards),
		"total": len(personAwards) + len(groupAwards), "light_mode": true,
	})
}

func (ac *AwardController) GetGroupMembersWithCertificates(c *gin.Context) {
	awardID, err := strconv.ParseUint(c.Param("award_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "无效的获奖记录ID"})
		return
	}

	members, err := ac.groupAwardService.GetGroupMembersWithCertificates(uint(awardID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "获取成员信息失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"award_id": awardID, "members": members, "total": len(members),
	})
}

func (ac *AwardController) UpdateMemberCertificate(c *gin.Context) {
	var update dto.MemberCertificateUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "请求参数错误: " + err.Error()})
		return
	}

	if err := ac.groupAwardService.UpdateMemberCertificate(update); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "更新失败: " + err.Error()})
		return
	}

	username, _ := auth.GetCurrentUsername(c)
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "成员证书更新成功",
		Data: map[string]interface{}{
			"updated_by": username, "award_id": update.AwardID,
			"member_name": update.MemberName, "member_id_number": update.MemberIDNumber,
		},
	})
}

// UploadCertificateImage 上传证书图片 — 校验 → 转PNG → 存文件
func (ac *AwardController) UploadCertificateImage(c *gin.Context) {
	file, err := c.FormFile("certificate")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "获取文件失败: " + err.Error()})
		return
	}

	// 校验文件大小
	if file.Size > utils.MaxFileSize() {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: fmt.Sprintf("文件过大，最大允许 %dMB", utils.MaxFileSize()/(1<<20)),
		})
		return
	}

	// 打开文件读取内容
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "读取文件失败: " + err.Error()})
		return
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "读取文件失败: " + err.Error()})
		return
	}

	// 校验 MIME 类型（通过文件内容 magic bytes）
	contentType := http.DetectContentType(fileBytes)
	if !utils.AllowedMIMETypes()[contentType] {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: fmt.Sprintf("不支持的文件类型: %s，仅允许图片格式", contentType),
		})
		return
	}

	// 转为 PNG
	pngBytes, err := utils.ConvertToPNG(fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "图片转换失败: " + err.Error()})
		return
	}

	// 保存到 uploads/certificates/
	certificatesDir := "uploads/certificates"
	if err := os.MkdirAll(certificatesDir, 0750); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "创建目录失败: " + err.Error()})
		return
	}

	filename := uuid.New().String() + ".png"
	filePath := filepath.Join(certificatesDir, filename)
	if err := os.WriteFile(filePath, pngBytes, 0640); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "保存文件失败: " + err.Error()})
		return
	}

	staticPath := "/static/certificates/" + filename
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "证书图片上传成功",
		Data:    map[string]interface{}{"file_path": staticPath},
	})
}

// BatchUploadCertificateImages 批量上传证书图片
func (ac *AwardController) BatchUploadCertificateImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "获取表单失败: " + err.Error()})
		return
	}

	files := form.File["certificates[]"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "未找到上传的文件"})
		return
	}

	certificatesDir := "uploads/certificates"
	if err := os.MkdirAll(certificatesDir, 0750); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "创建目录失败: " + err.Error()})
		return
	}

	results := make([]map[string]interface{}, 0)
	for _, file := range files {
		itemResult := map[string]interface{}{"filename": file.Filename, "success": false}

		// 校验文件大小
		if file.Size > utils.MaxFileSize() {
			itemResult["error"] = fmt.Sprintf("文件过大，最大允许 %dMB", utils.MaxFileSize()/(1<<20))
			results = append(results, itemResult)
			continue
		}

		// 读取文件
		src, err := file.Open()
		if err != nil {
			itemResult["error"] = "读取文件失败: " + err.Error()
			results = append(results, itemResult)
			continue
		}

		fileBytes := make([]byte, file.Size)
		if _, err := src.Read(fileBytes); err != nil {
			src.Close()
			itemResult["error"] = "读取文件失败: " + err.Error()
			results = append(results, itemResult)
			continue
		}
		src.Close()

		// 校验 MIME 类型
		contentType := http.DetectContentType(fileBytes)
		if !utils.AllowedMIMETypes()[contentType] {
			itemResult["error"] = fmt.Sprintf("不支持的文件类型: %s", contentType)
			results = append(results, itemResult)
			continue
		}

		// 转 PNG
		pngBytes, err := utils.ConvertToPNG(fileBytes)
		if err != nil {
			itemResult["error"] = "图片转换失败: " + err.Error()
			results = append(results, itemResult)
			continue
		}

		// 从原始文件名提取身份证号
		base := filepath.Base(file.Filename)
		idNumber := strings.TrimSuffix(base, filepath.Ext(base))

		filename := uuid.New().String() + ".png"
		filePath := filepath.Join(certificatesDir, filename)
		if err := os.WriteFile(filePath, pngBytes, 0640); err != nil {
			itemResult["error"] = "保存文件失败: " + err.Error()
			results = append(results, itemResult)
			continue
		}

		itemResult["success"] = true
		itemResult["file_path"] = "/static/certificates/" + filename
		itemResult["id_number"] = idNumber
		results = append(results, itemResult)
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: fmt.Sprintf("成功上传 %d 个文件", len(results)),
		Data:    map[string]interface{}{"results": results},
	})
}
```

- [ ] **Step 2: 添加 uuid 依赖**

```bash
cd backend && go get github.com/google/uuid@latest && go mod tidy
```

- [ ] **Step 3: 编译验证**

```bash
cd backend && go build ./...
```

Expected: 无编译错误。

- [ ] **Step 4: Commit**

```bash
git add backend/controllers/award_controller.go backend/services/group_award_service.go backend/go.mod backend/go.sum
git commit -m "refactor: PNG file storage, remove base64, add upload validation"
```

---

### Task 11: 新增请求频率限制中间件

**Files:**
- Create: `backend/middleware/rate_limit.go`

- [ ] **Step 1: 创建 rate_limit.go**

```go
package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 基于令牌桶的内存限流器
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*tokenBucket
	rate     float64 // 每秒允许的请求数
	capacity int     // 桶容量
}

type tokenBucket struct {
	tokens   float64
	lastTime time.Time
}

// NewRateLimiter 创建限流器
func NewRateLimiter(requestsPerMinute, capacity int) *RateLimiter {
	return &RateLimiter{
		buckets:  make(map[string]*tokenBucket),
		rate:     float64(requestsPerMinute) / 60.0,
		capacity: capacity,
	}
}

func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[key]
	now := time.Now()

	if !exists {
		rl.buckets[key] = &tokenBucket{
			tokens:   float64(rl.capacity) - 1,
			lastTime: now,
		}
		return true
	}

	// 计算补充的令牌
	elapsed := now.Sub(bucket.lastTime).Seconds()
	bucket.tokens += elapsed * rl.rate
	if bucket.tokens > float64(rl.capacity) {
		bucket.tokens = float64(rl.capacity)
	}
	bucket.lastTime = now

	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}

	return false
}

// 清理过期桶（每5分钟清理一次）
func (rl *RateLimiter) startCleanup() {
	go func() {
		rl.mu.Lock()
		rl.buckets = make(map[string]*tokenBucket)
		rl.mu.Unlock()
	}()
	time.AfterFunc(5*time.Minute, rl.startCleanup)
}

var (
	ipLimiter    *RateLimiter
	tokenLimiter *RateLimiter
	onceRL       sync.Once
)

func initLimiters() {
	onceRL.Do(func() {
		ipLimiter = NewRateLimiter(30, 30)    // 30 req/min, burst 30
		tokenLimiter = NewRateLimiter(60, 60) // 60 req/min, burst 60

		time.AfterFunc(5*time.Minute, ipLimiter.startCleanup)
		time.AfterFunc(5*time.Minute, tokenLimiter.startCleanup)
	})
}

// IPRateLimit 基于 IP 的限流中间件（用于公开查询接口）
func IPRateLimit() gin.HandlerFunc {
	initLimiters()
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !ipLimiter.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// TokenRateLimit 基于 Token 的限流中间件（用于认证接口）
func TokenRateLimit() gin.HandlerFunc {
	initLimiters()
	return func(c *gin.Context) {
		key, exists := c.Get("username")
		if !exists {
			key = c.ClientIP()
		}
		keyStr, ok := key.(string)
		if !ok {
			keyStr = c.ClientIP()
		}

		if !tokenLimiter.allow(keyStr) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
```

- [ ] **Step 2: 编译验证**

```bash
cd backend && go build ./...
```

Expected: 无编译错误。

- [ ] **Step 3: Commit**

```bash
git add backend/middleware/rate_limit.go
git commit -m "feat: add token-bucket rate limiting middleware"
```

---

### Task 12: 确认 main.go 无需改动

**Files:**
- 无需修改

> **说明**：限流中间件全部在 Task 13 的 `routes/routes.go` 中按路由组注册（`auth.Use(middleware.IPRateLimit())` 等），main.go 无需全局注册。本 Task 跳过。

- [ ] **Step 1: 验证 main.go 编译通过**

```bash
cd backend && go build ./...
```

Expected: 无编译错误。

---

### Task 13: 更新路由 — 添加限流 + 增强健康检查

**Files:**
- Modify: `backend/routes/routes.go`

- [ ] **Step 1: 重写 routes.go**

```go
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
```

- [ ] **Step 2: 移除不再需要的 import**

确保 routes.go 不再 import 不需要的包（http 和 gin 仍需要）。

- [ ] **Step 3: 编译验证**

```bash
cd backend && go build ./...
```

Expected: 无编译错误。

- [ ] **Step 4: Commit**

```bash
git add backend/routes/routes.go
git commit -m "feat: add rate limiting to routes, enhance health check with DB ping"
```

---

### Task 14: 新建 Dockerfile

**Files:**
- Create: `backend/Dockerfile`

- [ ] **Step 1: 创建 Dockerfile**

```dockerfile
# Stage 1: 编译
FROM golang:1.21-alpine AS builder

RUN apk --no-cache add git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o awards-management .

# Stage 2: 运行
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/awards-management .
COPY config.yaml .
RUN mkdir -p uploads/certificates

EXPOSE 8000

CMD ["./awards-management"]
```

- [ ] **Step 2: Commit**

```bash
git add backend/Dockerfile
git commit -m "feat: add multi-stage Dockerfile"
```

---

### Task 15: 新建 docker-compose.yml

**Files:**
- Create: `backend/docker-compose.yml`

- [ ] **Step 1: 创建 docker-compose.yml**

```yaml
version: "3.8"

services:
  mysql:
    image: mysql:8.0
    container_name: awards-mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_ROOT_PASSWORD}
      MYSQL_DATABASE: ${DB_NAME:-awards_management}
      MYSQL_USER: ${DB_USER:-awards_user}
      MYSQL_PASSWORD: ${DB_PASSWORD}
    ports:
      - "127.0.0.1:3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - awards-net

  backend:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: awards-backend
    restart: always
    environment:
      GIN_MODE: ${GIN_MODE:-release}
      DB_HOST: mysql
      DB_PORT: "3306"
      DB_USER: ${DB_USER:-awards_user}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: ${DB_NAME:-awards_management}
      JWT_SECRET: ${JWT_SECRET}
    volumes:
      - uploads_data:/app/uploads
    depends_on:
      mysql:
        condition: service_healthy
    networks:
      - awards-net

  nginx:
    image: nginx:1.25-alpine
    container_name: awards-nginx
    restart: always
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/conf.d/default.conf:ro
      - ../frontend/dist:/usr/share/nginx/html:ro
    depends_on:
      - backend
    networks:
      - awards-net

volumes:
  mysql_data:
  uploads_data:

networks:
  awards-net:
    driver: bridge
```

- [ ] **Step 2: Commit**

```bash
git add backend/docker-compose.yml
git commit -m "feat: add docker-compose with MySQL, backend, nginx"
```

---

### Task 16: 新建 nginx.conf

**Files:**
- Create: `backend/nginx.conf`

- [ ] **Step 1: 创建 nginx.conf**

```nginx
server {
    listen 80;
    server_name _;

    client_max_body_size 20m;

    # gzip
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml;

    # 前端静态文件
    location / {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /index.html;
    }

    location /admin {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /admin.html;
    }

    location /verify {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /verify.html;
    }

    # API 反向代理
    location /api/ {
        proxy_pass http://backend:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 120s;
    }

    # 证书静态资源
    location /static/ {
        proxy_pass http://backend:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_read_timeout 30s;
    }
}
```

- [ ] **Step 2: Commit**

```bash
git add backend/nginx.conf
git commit -m "feat: add nginx reverse proxy configuration"
```

---

### Task 17: 新建 .env.example + 更新 .gitignore

**Files:**
- Create: `backend/.env.example`
- Modify: `backend/.gitignore`

- [ ] **Step 1: 创建 .env.example**

```
# 数据库
DB_ROOT_PASSWORD=change-me-root-password
DB_USER=awards_user
DB_PASSWORD=change-me-db-password
DB_NAME=awards_management

# JWT
JWT_SECRET=change-me-to-a-random-256-bit-key

# Gin 模式: debug / release
GIN_MODE=release
```

- [ ] **Step 2: 更新 .gitignore**

确保 `backend/.gitignore` 包含以下行（如果没有则追加）：

```
# 实际环境变量（含密钥，不可提交）
.env
```

同时移除或注释掉 `config/local.yaml` 和 `config/production.yaml` 相关行（如存在），并在末尾添加：

```
# 可执行文件
awards-management
awards-management.exe
```

- [ ] **Step 3: 删除仓库中的可执行文件**

```bash
rm -f backend/awards-management.exe
rm -f backend/data.db
```

- [ ] **Step 4: Commit**

```bash
git add backend/.env.example backend/.gitignore
git rm --cached backend/awards-management.exe backend/data.db 2>/dev/null || true
git commit -m "feat: add .env.example, update gitignore, remove binaries"
```

---

### Task 18: 最终编译验证 + 完整构建测试

**Files:** 全部

- [ ] **Step 1: 完整编译**

```bash
cd backend && go build -o awards-management-test ./...
```

Expected: 无编译错误，生成二进制文件。

- [ ] **Step 2: 清理测试二进制**

```bash
rm -f backend/awards-management-test
```

- [ ] **Step 3: 构建前端（确认无破坏）**

```bash
cd frontend && npm install && npm run build
```

Expected: 构建成功，dist/ 目录生成。

- [ ] **Step 4: 删除仓库中不该存在的文件**

```bash
# 删除 .DS_Store
find . -name ".DS_Store" -delete
```

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "chore: final build verification, clean up .DS_Store"
```

---

### Task 19: 创建 MySQL 初始化脚本

**Files:**
- Create: `backend/init.sql`

- [ ] **Step 1: 创建 init.sql**

```sql
CREATE DATABASE IF NOT EXISTS awards_management CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

- [ ] **Step 2: Commit**

```bash
git add backend/init.sql
git commit -m "feat: add MySQL init script for docker-compose"
```

---

## 部署操作步骤

以上任务完成后，按以下步骤部署：

```bash
# 1. 复制环境变量模板并填写真实值
cd backend
cp .env.example .env
# 编辑 .env，填写强密码和随机 JWT Secret

# 2. 构建前端
cd ../frontend
npm install && npm run build

# 3. 启动服务
cd ../backend
docker compose up -d

# 4. 检查服务状态
docker compose ps
curl http://localhost/health

# 5. 首次登录后修改默认密码
# 默认账户: admin / admin123
# 登录后立即通过 /api/v1/auth/change-password 修改密码
```

---

*计划更新: 2026-06-10*
