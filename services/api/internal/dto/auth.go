package dto

type SignUpReq struct {
	RegPlatform uint8  `json:"reg_platform" binding:"required,oneof=1 2 3 4 5"` // 注册平台 1:iOS 2:安卓
	Nickname    string `json:"nickname" binding:"required,min=1,max=20"`        // 昵称
	Account     string `json:"account" binding:"required,min=6,max=20"`         // 密码
	Password    string `json:"password" binding:"required,len=32"`              // 密码
	Gender      int32  `json:"gender" binding:"omitempty,oneof=0 1 2"`          // 性别
	BirthTs     int64  `json:"birth_ts" binding:"omitempty,gt=0"`               // 生日
	Email       string `json:"email" binding:"omitempty,email"`                 // Email
	Mobile      string `json:"mobile" binding:"required,min=8,max=20"`          // 手机号
	AvatarKey   string `json:"avatar_key" binding:"omitempty"`                  // 头像(暂时弃用)
	CityId      int64  `json:"city_id" binding:"omitempty,gte=0"`               // 城市ID
	Uuid        string `json:"uuid" binding:"required,len=40"`                  // UDID
}

type SignInReq struct {
	Account  string `json:"account" binding:"required,min=1,max=20"`  // 账户
	Password string `json:"password" binding:"required,min=8,max=20"` // 密码
}
