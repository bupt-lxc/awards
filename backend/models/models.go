package models

import (
	"time"

	"gorm.io/gorm"
)

// Admin 管理员模型
type Admin struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string    `json:"username" gorm:"uniqueIndex;not null;size:50"`
	PasswordHash string    `json:"-" gorm:"not null;size:255"`
	Role         string    `json:"role" gorm:"default:admin;size:20"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// PersonAward 个人获奖记录模型
type PersonAward struct {
	ID                uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Year              int       `json:"year" gorm:"not null;index"`
	PersonName        string    `json:"person_name" gorm:"not null;size:100;index"`
	IDNumber          string    `json:"id_number" gorm:"not null;size:50;index"`
	Category          string    `json:"category" gorm:"not null;size:100;index"`
	AwardLevel        string    `json:"award_level" gorm:"not null;size:50"`
	Description       string    `json:"description" gorm:"type:longtext"`
	CertificateImage  string    `json:"certificate_image" gorm:"size:500"`
	CertificateNumber string    `json:"certificate_number" gorm:"size:100"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// GroupAward 团体获奖记录模型
type GroupAward struct {
	ID                uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	Year              int           `json:"year" gorm:"not null;index:idx_award_info,priority:1"`
	GroupName         string        `json:"group_name" gorm:"not null;size:200;index:idx_award_info,priority:2"`
	Category          string        `json:"category" gorm:"not null;size:100;index:idx_award_info,priority:3"`
	AwardLevel        string        `json:"award_level" gorm:"not null;size:50;index:idx_award_level"`
	Description       string        `json:"description" gorm:"type:text"`
	CertificateImage  string        `json:"certificate_image" gorm:"size:500"`
	CertificateNumber string        `json:"certificate_number" gorm:"size:100"`
	GroupMembers      []GroupMember `json:"group_members" gorm:"foreignKey:GroupAwardID;constraint:OnDelete:CASCADE"`
	CreatedAt         time.Time     `json:"created_at" gorm:"autoCreateTime;index:idx_created_at"`
	UpdatedAt         time.Time     `json:"updated_at" gorm:"autoUpdateTime"`

	// 添加复合索引说明
	// 1. idx_award_info: 年份+团队名称+组别的复合索引，用于快速查找获奖信息
	// 2. idx_award_level: 奖项等级索引，用于按奖项等级筛选
	// 3. idx_created_at: 创建时间索引，用于按时间排序和筛选
}

// GroupMember 团体成员模型（用于高效查询）
type GroupMember struct {
	ID                uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupAwardID      uint   `json:"group_award_id" gorm:"not null;index:idx_group_award_id"`
	MemberName        string `json:"member_name" gorm:"not null;size:100;index:idx_member_info,priority:2"`
	MemberIDNumber    string `json:"member_id_number" gorm:"not null;size:50;index:idx_member_info,priority:1"`
	CertificateImage  string `json:"certificate_image" gorm:"size:500"`
	CertificateNumber string `json:"certificate_number" gorm:"size:100"`

	// 添加复合索引
	// 1. 团队ID和成员信息的复合索引，用于快速查找团队成员
	// 2. 成员姓名和证件号的复合索引，用于快速查找成员
	// gorm:"index:idx_group_member,priority:1;index:idx_member_info,priority:1,2"
}

// Member 成员结构体（用于JSON序列化）
type Member struct {
	Name             string `json:"name" binding:"required"`
	IDNumber         string `json:"id_number" binding:"required"`
	CertificateImage string `json:"certificate_image,omitempty"`
}

// TableName 设置表名
func (Admin) TableName() string {
	return "admins"
}

func (PersonAward) TableName() string {
	return "person_awards"
}

func (GroupAward) TableName() string {
	return "group_awards"
}

func (GroupMember) TableName() string {
	return "group_members"
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Admin{},
		&PersonAward{},
		&GroupAward{},
		&GroupMember{},
	)
}
