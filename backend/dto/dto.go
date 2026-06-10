package dto

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required,max=50"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}

// AdminChangePasswordRequest 修改管理员密码请求
type AdminChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=128"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=128"`
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
	PersonName string `json:"person_name" binding:"required,max=100"`
	IDNumber   string `json:"id_number" binding:"required,max=50"`
	Category   string `json:"category,omitempty" binding:"max=100"`
	Page       int    `json:"page,omitempty"`      // 页码，从1开始；为0时不分页
	PageSize   int    `json:"page_size,omitempty"` // 每页条数，默认20，最大100
}

// GroupAwardQuery 团体获奖查询请求
type GroupAwardQuery struct {
	PersonName string `json:"person_name,omitempty" binding:"max=100"`
	IDNumber   string `json:"id_number,omitempty" binding:"max=50"`
	Year       *int   `json:"year,omitempty"`
	GroupName  string `json:"group_name,omitempty" binding:"max=200"`
	Category   string `json:"category,omitempty" binding:"max=100"`
	Page       int    `json:"page,omitempty"`      // 页码，从1开始；为0时不分页
	PageSize   int    `json:"page_size,omitempty"` // 每页条数，默认20，最大100
}

// PersonAwardRecord 个人获奖记录
type PersonAwardRecord struct {
	Year             int    `json:"year" binding:"required"`
	PersonName       string `json:"person_name" binding:"required,max=100"`
	IDNumber         string `json:"id_number" binding:"required,max=50"`
	Category         string `json:"category" binding:"required,max=100"`
	AwardLevel       string `json:"award_level" binding:"required,max=50"`
	Description      string `json:"description,omitempty"`
	CertificateImage string `json:"certificate_image,omitempty" binding:"max=500"`
}

// GroupAwardRecord 团体获奖记录
type GroupAwardRecord struct {
	Year             int      `json:"year" binding:"required"`
	GroupName        string   `json:"group_name" binding:"required,max=200"`
	Category         string   `json:"category" binding:"required,max=100"`
	AwardLevel       string   `json:"award_level" binding:"required,max=50"`
	Description      string   `json:"description,omitempty"`
	Members          []Member `json:"members" binding:"required,dive"`
	CertificateImage string   `json:"certificate_image,omitempty" binding:"max=500"`
}

// PersonAwardUpdate 个人获奖更新请求
type PersonAwardUpdate struct {
	ID               uint   `json:"id" binding:"required"`
	Year             *int   `json:"year,omitempty"`
	PersonName       string `json:"person_name,omitempty" binding:"max=100"`
	IDNumber         string `json:"id_number,omitempty" binding:"max=50"`
	Category         string `json:"category,omitempty" binding:"max=100"`
	AwardLevel       string `json:"award_level,omitempty" binding:"max=50"`
	Description      string `json:"description,omitempty"`
	CertificateImage string `json:"certificate_image,omitempty" binding:"max=500"`
}

// GroupAwardUpdate 团体获奖更新请求
type GroupAwardUpdate struct {
	ID               uint     `json:"id" binding:"required"`
	Year             *int     `json:"year,omitempty"`
	GroupName        string   `json:"group_name,omitempty" binding:"max=200"`
	Category         string   `json:"category,omitempty" binding:"max=100"`
	AwardLevel       string   `json:"award_level,omitempty" binding:"max=50"`
	Description      string   `json:"description,omitempty"`
	CertificateImage string   `json:"certificate_image,omitempty" binding:"max=500"`
	Members          []Member `json:"members,omitempty"`
}

// MemberCertificateUpdate 成员证书更新请求
type MemberCertificateUpdate struct {
	AwardID          uint   `json:"award_id" binding:"required"`
	MemberName       string `json:"member_name" binding:"required,max=100"`
	MemberIDNumber   string `json:"member_id_number" binding:"required,max=50"`
	CertificateImage string `json:"certificate_image" binding:"required,max=500"`
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
	Name             string `json:"name" binding:"required,max=100"`
	IDNumber         string `json:"id_number" binding:"required,max=50"`
	CertificateImage string `json:"certificate_image,omitempty" binding:"max=500"`
}
