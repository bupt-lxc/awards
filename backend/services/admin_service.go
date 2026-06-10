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
