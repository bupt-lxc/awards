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
		serverError(c, "查询失败", err)
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
		serverError(c, "查询失败", err)
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
		serverError(c, "个人获奖记录上传失败", err)
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
		serverError(c, "团体获奖记录上传失败", err)
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
		serverError(c, "批量上传失败", err)
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
		serverError(c, "批量上传失败", err)
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
		serverError(c, "更新失败", err)
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
		serverError(c, "更新失败", err)
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
		serverError(c, "删除失败", deleteErr)
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
		serverError(c, "获取个人获奖记录失败", err)
		return
	}

	groupAwards, err := ac.groupAwardService.GetAllGroupAwards()
	if err != nil {
		serverError(c, "获取团体获奖记录失败", err)
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
		serverError(c, "获取个人获奖记录失败", err)
		return
	}

	groupAwards, err := ac.groupAwardService.GetAllGroupAwardsLight()
	if err != nil {
		serverError(c, "获取团体获奖记录失败", err)
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
		serverError(c, "获取成员信息失败", err)
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
		serverError(c, "更新失败", err)
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

	// 打开文件读取全部内容
	src, err := file.Open()
	if err != nil {
		serverError(c, "读取文件失败", err)
		return
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		serverError(c, "读取文件失败", err)
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
		serverError(c, "图片转换失败", err)
		return
	}

	// 保存到 uploads/certificates/
	certificatesDir := "uploads/certificates"
	if err := os.MkdirAll(certificatesDir, 0750); err != nil {
		serverError(c, "创建目录失败", err)
		return
	}

	filename := uuid.New().String() + ".png"
	filePath := filepath.Join(certificatesDir, filename)
	if err := os.WriteFile(filePath, pngBytes, 0640); err != nil {
		serverError(c, "保存文件失败", err)
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
		serverError(c, "创建目录失败", err)
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

		fileBytes, err := io.ReadAll(src)
		src.Close()
		if err != nil {
			itemResult["error"] = "读取文件失败: " + err.Error()
			results = append(results, itemResult)
			continue
		}

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
