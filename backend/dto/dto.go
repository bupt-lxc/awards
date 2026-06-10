package dto

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminChangePasswordRequest 修改管理员密码请求
type AdminChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// TokenResponse 令牌响应
type TokenResponse struct {
	AccessToken string                 `json:"access_token"`
	TokenType   string                 `json:"token_type"`
	ExpiresIn   int                    `json:"expires_in"`
	UserInfo    map[string]interface{} `json:"user_info"`
}

// PersonAwardQuery 个人获奖查询请求
type PersonAwardQuery struct {
	Year       *int   `json:"year,omitempty"`
	PersonName string `json:"person_name" binding:"required"`
	IDNumber   string `json:"id_number" binding:"required"`
	Category   string `json:"category,omitempty"`
}

// GroupAwardQuery 团体获奖查询请求
type GroupAwardQuery struct {
	PersonName string `json:"person_name,omitempty"`
	IDNumber   string `json:"id_number,omitempty"`
	Year       *int   `json:"year,omitempty"`
	GroupName  string `json:"group_name" binding:"required"`
	Category   string `json:"category,omitempty"`
}

// PersonAwardRecord 个人获奖记录
type PersonAwardRecord struct {
	Year             int    `json:"year" binding:"required"`
	PersonName       string `json:"person_name" binding:"required"`
	IDNumber         string `json:"id_number" binding:"required"`
	Category         string `json:"category" binding:"required"`
	AwardLevel       string `json:"award_level" binding:"required"`
	Description      string `json:"description,omitempty"`
	CertificateImage string `json:"certificate_image,omitempty"`
}

// GroupAwardRecord 团体获奖记录
type GroupAwardRecord struct {
	Year             int      `json:"year" binding:"required"`
	GroupName        string   `json:"group_name" binding:"required"`
	Category         string   `json:"category" binding:"required"`
	AwardLevel       string   `json:"award_level" binding:"required"`
	Description      string   `json:"description,omitempty"`
	Members          []Member `json:"members" binding:"required,dive"`
	CertificateImage string   `json:"certificate_image,omitempty"` // 新增：团队奖状图片
}

// PersonAwardUpdate 个人获奖更新请求
type PersonAwardUpdate struct {
	ID               uint   `json:"id" binding:"required"`
	Year             *int   `json:"year,omitempty"`
	PersonName       string `json:"person_name,omitempty"`
	IDNumber         string `json:"id_number,omitempty"`
	Category         string `json:"category,omitempty"`
	AwardLevel       string `json:"award_level,omitempty"`
	Description      string `json:"description,omitempty"`
	CertificateImage string `json:"certificate_image,omitempty"`
}

// GroupAwardUpdate 团体获奖更新请求
type GroupAwardUpdate struct {
	ID               uint     `json:"id" binding:"required"`
	Year             *int     `json:"year,omitempty"`
	GroupName        string   `json:"group_name,omitempty"`
	Category         string   `json:"category,omitempty"`
	AwardLevel       string   `json:"award_level,omitempty"`
	Description      string   `json:"description,omitempty"`
	CertificateImage string   `json:"certificate_image,omitempty"`
	Members          []Member `json:"members,omitempty"`
}

// MemberCertificateUpdate 成员证书更新请求
type MemberCertificateUpdate struct {
	AwardID          uint   `json:"award_id" binding:"required"`
	MemberName       string `json:"member_name" binding:"required"`
	MemberIDNumber   string `json:"member_id_number" binding:"required"`
	CertificateImage string `json:"certificate_image" binding:"required"`
}

// QueryResponse 查询响应
type QueryResponse struct {
	Awards      interface{}            `json:"awards"`
	Total       int                    `json:"total"`
	QueryParams map[string]interface{} `json:"query_params"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error string `json:"error"`
}

// UploadResponse 上传响应
type UploadResponse struct {
	Message    string `json:"message"`
	AwardID    uint   `json:"award_id,omitempty"`
	UploadedBy string `json:"uploaded_by,omitempty"`
}

// BatchUploadResponse 批量上传响应
type BatchUploadResponse struct {
	Message     string `json:"message"`
	Successful  int    `json:"successful"`
	Failed      int    `json:"failed"`
	UploadedBy  string `json:"uploaded_by"`
	FailedItems []int  `json:"failed_items,omitempty"`
}

type Member struct {
	Name             string `json:"name"`
	IDNumber         string `json:"id_number"`
	CertificateImage string `json:"certificate_image"`
}
