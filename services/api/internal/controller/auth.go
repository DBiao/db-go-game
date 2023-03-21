package controller

import (
	"db-go-game/pkg/common/dgin"
	"db-go-game/pkg/dhttp"
	"db-go-game/services/api/internal/dto"
	"db-go-game/services/api/internal/service"
	"github.com/gin-gonic/gin"
)

type IUserController interface {
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	SignOut(ctx *gin.Context)
}

type userController struct {
	userService service.IUserService
}

func NewUserController(userService service.IUserService) IUserController {
	return &userController{userService: userService}
}

func (u *userController) SignUp(ctx *gin.Context) {
	var req *dto.SignUpReq
	if err := dgin.BindJSON(ctx, req); err != nil {
		dhttp.Error(ctx, dhttp.ERROR_CODE_HTTP_REQ_PARAM_ERR, err.Error())
		return
	}

	resp := u.userService.SignUp(req)
	dhttp.Success(ctx, resp)
}

func (u *userController) SignIn(ctx *gin.Context) {
	var req *dto.SignInReq
	if err := dgin.BindJSON(ctx, req); err != nil {
		dhttp.Error(ctx, dhttp.ERROR_CODE_HTTP_REQ_PARAM_ERR, err.Error())
		return
	}

	resp := u.userService.SignIn(req)
	dhttp.Success(ctx, resp)
}

func (u *userController) RefreshToken(ctx *gin.Context) {
	var token string
	resp := u.userService.RefreshToken(token)
	dhttp.Success(ctx, resp)
}

func (u *userController) SignOut(ctx *gin.Context) {
	var token string
	resp := u.userService.SignOut(token)
	dhttp.Success(ctx, resp)
}
