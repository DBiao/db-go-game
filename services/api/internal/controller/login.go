package controller

import (
	"db-go-game/pkg/common/xgin"
	"db-go-game/pkg/xhttp"
	"db-go-game/services/api/internal/dto"
	"db-go-game/services/api/internal/service"
	"github.com/gin-gonic/gin"
)

type IUserController interface {
	Login(ctx *gin.Context)
}

type userController struct {
	userService service.IUserService
}

func NewUserController(userService service.IUserService) IUserController {
	return &userController{userService: userService}
}

func (u *userController) Login(ctx *gin.Context) {
	var req *dto.LoginReq
	if err := xgin.BindJSON(ctx, req); err != nil {
		xhttp.Error(ctx, xhttp.ERROR_CODE_HTTP_REQ_PARAM_ERR, err.Error())
		return
	}

	resp := u.userService.Login(req)
	xhttp.Success(ctx, resp)
}
