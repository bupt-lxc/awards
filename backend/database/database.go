package database

import (
	"awards-management/config"
	"awards-management/models"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() error {
	dbConfig := config.AppConfig.Database

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
		dbConfig.Charset,
		dbConfig.ParseTime,
		dbConfig.Loc,
	)

	log.Printf("正在连接数据库: %s:%d", dbConfig.Host, dbConfig.Port)

	// 配置GORM日志
	var gormLogger logger.Interface
	if config.AppConfig.Server.Mode == "debug" {
		gormLogger = logger.New(
			log.New(log.Writer(), "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
			},
		)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// 连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Printf("数据库连接失败: %v", err)
		log.Printf("DSN: %s:***@tcp(%s:%d)/%s", dbConfig.Username, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 获取底层sql.DB对象进行连接池配置
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("获取sql.DB失败: %v", err)
		return fmt.Errorf("获取sql.DB失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Second)

	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		log.Printf("数据库Ping测试失败: %v", err)
		return fmt.Errorf("数据库Ping测试失败: %v", err)
	}

	log.Printf("数据库连接成功")

	// 自动迁移数据库表
	if err := models.AutoMigrate(DB); err != nil {
		log.Printf("数据库迁移失败: %v", err)
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	log.Printf("数据库迁移完成")

	// 创建默认管理员账户
	if err := createDefaultAdmin(); err != nil {
		log.Printf("创建默认管理员失败: %v", err)
		return fmt.Errorf("创建默认管理员失败: %v", err)
	}

	log.Println("数据库初始化完成")
	return nil
}

// createDefaultAdmin 创建初始管理员账户。
// 出于安全考虑，首次启动必须通过环境变量显式提供密码。
func createDefaultAdmin() error {
	// 检查是否已存在admin用户
	var count int64
	DB.Model(&models.Admin{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		return nil // 已存在，不需要创建
	}

	// 初始密码只允许通过环境变量提供，避免新部署使用公开默认密码。
	defaultPassword := os.Getenv("ADMIN_DEFAULT_PASSWORD")
	if defaultPassword == "" {
		return fmt.Errorf("未设置 ADMIN_DEFAULT_PASSWORD，无法创建初始管理员")
	}

	// 使用 bcrypt 生成密码哈希 (cost=12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), 12)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// 创建默认管理员
	admin := models.Admin{
		Username:     "admin",
		PasswordHash: string(hashedPassword),
		Role:         "admin",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return DB.Create(&admin).Error
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
