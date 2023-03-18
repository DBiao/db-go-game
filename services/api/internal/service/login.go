package service

import (
	"db-go-game/pkg/xhttp"
	"db-go-game/services/api/internal/dto"
)

type IUserService interface {
	Login(req *dto.LoginReq) *xhttp.Resp
}

type userService struct {
}

func NewUserService() IUserService {
	return &userService{}
}

func (u *userService) Login(req *dto.LoginReq) *xhttp.Resp {
	resp := new(xhttp.Resp)

	resp.Data = "token"
	return resp
}
