package dto

type LoginReq struct {
	Account  string `json:"account" validate:"required,min=1,max=20"`  // 账户
	Password string `json:"password" validate:"required,min=8,max=20"` // 密码
}
