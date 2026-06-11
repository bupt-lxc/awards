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

func (s *PersonAwardService) QueryPersonAwards(query dto.PersonAwardQuery) ([]models.PersonAward, int64, error) {
	var awards []models.PersonAward
	var total int64

	db := database.DB.Model(&models.PersonAward{})

	if query.PersonName != "" {
		db = db.Where("person_name = ?", query.PersonName)
	}
	if query.IDNumber != "" {
		db = db.Where("id_number = ?", query.IDNumber)
	}
	if query.Year != nil {
		db = db.Where("year = ?", *query.Year)
	}
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}

	// 先统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页：page>0 时分页，否则返回全部（向下兼容）
	if query.Page > 0 {
		pageSize := query.PageSize
		if pageSize <= 0 {
			pageSize = 20
		}
		if pageSize > 100 {
			pageSize = 100
		}
		db = db.Offset((query.Page - 1) * pageSize).Limit(pageSize)
	}

	err := db.Order("created_at DESC").Find(&awards).Error
	return awards, total, err
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
