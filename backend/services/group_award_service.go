package services

import (
	"awards-management/database"
	"awards-management/dto"
	"awards-management/models"
	"awards-management/storage"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

type GroupAwardService struct{}

func NewGroupAwardService() *GroupAwardService {
	return &GroupAwardService{}
}

func (s *GroupAwardService) QueryGroupAwards(query dto.GroupAwardQuery) ([]map[string]interface{}, int64, error) {
	var awards []models.GroupAward
	var total int64

	dbQuery := database.DB.Model(&models.GroupAward{}).
		Preload("GroupMembers").
		Order("created_at DESC")

	if query.GroupName != "" {
		dbQuery = dbQuery.Where("group_name LIKE ?", "%"+query.GroupName+"%")
	}
	if query.Year != nil {
		dbQuery = dbQuery.Where("year = ?", *query.Year)
	}
	if query.Category != "" {
		dbQuery = dbQuery.Where("category = ?", query.Category)
	}
	if query.PersonName != "" {
		dbQuery = dbQuery.Where("EXISTS (SELECT 1 FROM group_members WHERE group_award_id = group_awards.id AND member_name = ?)",
			query.PersonName)
	}
	if query.IDNumber != "" {
		dbQuery = dbQuery.Where("EXISTS (SELECT 1 FROM group_members WHERE group_award_id = group_awards.id AND member_id_number = ?)",
			query.IDNumber)
	}

	// 先统计总数
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询团队信息失败: %v", err)
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
		dbQuery = dbQuery.Offset((query.Page - 1) * pageSize).Limit(pageSize)
	}

	if err := dbQuery.Find(&awards).Error; err != nil {
		return nil, 0, fmt.Errorf("查询团队信息失败: %v", err)
	}

	log.Printf("查询到 %d 条团队记录 (总计 %d)", len(awards), total)

	resultAwards := make([]map[string]interface{}, 0, len(awards))

	for _, award := range awards {
		certificateImage := strings.TrimSpace(award.CertificateImage)

		awardMap := map[string]interface{}{
			"id":                 award.ID,
			"year":               award.Year,
			"group_name":         award.GroupName,
			"category":           award.Category,
			"award_level":        award.AwardLevel,
			"description":        award.Description,
			"certificate_image":  certificateImage,
			"certificate_number": strings.TrimSpace(award.CertificateNumber),
			"created_at":         award.CreatedAt,
			"updated_at":         award.UpdatedAt,
			"members":            []map[string]interface{}{},
		}

		for _, member := range award.GroupMembers {
			if (query.PersonName == "" || query.PersonName == member.MemberName) &&
				(query.IDNumber == "" || query.IDNumber == member.MemberIDNumber) {

				memberMap := map[string]interface{}{
					"name":               member.MemberName,
					"id_number":          member.MemberIDNumber,
					"certificate_image":  strings.TrimSpace(member.CertificateImage),
					"certificate_number": strings.TrimSpace(member.CertificateNumber),
				}
				awardMap["members"] = append(awardMap["members"].([]map[string]interface{}), memberMap)
			}
		}

		resultAwards = append(resultAwards, awardMap)
	}

	return resultAwards, total, nil
}

func (s *GroupAwardService) CreateGroupAward(record dto.GroupAwardRecord) (*models.GroupAward, error) {
	award := models.GroupAward{
		Year:              record.Year,
		GroupName:         record.GroupName,
		Category:          record.Category,
		AwardLevel:        record.AwardLevel,
		Description:       record.Description,
		CertificateImage:  record.CertificateImage,
		CertificateNumber: record.CertificateNumber,
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
		groupMember := models.GroupMember{
			GroupAwardID:      award.ID,
			MemberName:        member.Name,
			MemberIDNumber:    member.IDNumber,
			CertificateImage:  member.CertificateImage,
			CertificateNumber: member.CertificateNumber,
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
	if update.CertificateNumber != nil {
		updateData["certificate_number"] = *update.CertificateNumber
	}

	if update.Members != nil {
		if err := tx.Where("group_award_id = ?", update.ID).Delete(&models.GroupMember{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		for _, member := range update.Members {
			gm := models.GroupMember{
				GroupAwardID:      update.ID,
				MemberName:        member.Name,
				MemberIDNumber:    member.IDNumber,
				CertificateImage:  member.CertificateImage,
				CertificateNumber: member.CertificateNumber,
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

	// 收集需要删除的文件（保留完整路径传递给存储服务）
	filesToDelete := make(map[string]struct{})
	if award.CertificateImage != "" {
		filesToDelete[award.CertificateImage] = struct{}{}
	}
	for _, m := range members {
		if m.CertificateImage != "" {
			filesToDelete[m.CertificateImage] = struct{}{}
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

	// 异步清理文件（通过存储服务，MinIO 或本地）
	go func(files map[string]struct{}) {
		store := storage.Get()
		if store == nil {
			return
		}
		for file := range files {
			if err := store.Delete(file); err != nil {
				log.Printf("清理证书文件失败 [%s]: %v", file, err)
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
		"id":                 award.ID,
		"year":               award.Year,
		"group_name":         award.GroupName,
		"category":           award.Category,
		"award_level":        award.AwardLevel,
		"description":        award.Description,
		"certificate_image":  award.CertificateImage,
		"certificate_number": award.CertificateNumber,
		"created_at":         award.CreatedAt,
		"updated_at":         award.UpdatedAt,
		"members":            []map[string]interface{}{},
	}

	for _, member := range award.GroupMembers {
		result["members"] = append(result["members"].([]map[string]interface{}), map[string]interface{}{
			"name":               member.MemberName,
			"id_number":          member.MemberIDNumber,
			"certificate_image":  member.CertificateImage,
			"certificate_number": member.CertificateNumber,
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
			"id":                 award.ID,
			"year":               award.Year,
			"group_name":         award.GroupName,
			"category":           award.Category,
			"award_level":        award.AwardLevel,
			"description":        award.Description,
			"certificate_image":  award.CertificateImage,
			"certificate_number": award.CertificateNumber,
			"created_at":         award.CreatedAt,
			"updated_at":         award.UpdatedAt,
			"members":            []map[string]interface{}{},
		}

		for _, member := range award.GroupMembers {
			if member.MemberName == "" || member.MemberIDNumber == "" {
				continue
			}
			result["members"] = append(result["members"].([]map[string]interface{}), map[string]interface{}{
				"name":               member.MemberName,
				"id_number":          member.MemberIDNumber,
				"certificate_image":  member.CertificateImage,
				"certificate_number": member.CertificateNumber,
			})
		}

		results = append(results, result)
	}

	return results, nil
}

func (s *GroupAwardService) GetAllGroupAwardsLight() ([]map[string]interface{}, error) {
	var awards []models.GroupAward

	err := database.DB.Select("id, year, group_name, category, award_level, description, certificate_image, certificate_number, created_at, updated_at").
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
			"id":                 award.ID,
			"year":               award.Year,
			"group_name":         award.GroupName,
			"category":           award.Category,
			"award_level":        award.AwardLevel,
			"description":        award.Description,
			"certificate_number": award.CertificateNumber,
			"created_at":         award.CreatedAt,
			"updated_at":         award.UpdatedAt,
			"has_certificate":    award.CertificateImage != "",
			"members":            []map[string]interface{}{},
		}

		for _, member := range award.GroupMembers {
			if member.MemberName == "" || member.MemberIDNumber == "" {
				continue
			}
			result["members"] = append(result["members"].([]map[string]interface{}), map[string]interface{}{
				"name":               member.MemberName,
				"id_number":          member.MemberIDNumber,
				"has_certificate":    member.CertificateImage != "",
				"certificate_number": member.CertificateNumber,
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
