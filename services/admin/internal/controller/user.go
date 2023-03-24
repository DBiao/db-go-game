package controller

import "db-go-game/services/admin/internal/service"

type IUserController interface {
}

type userController struct {
}

func NewUserController(userService service.IUserService) IUserController {
	return &userController{}
}
