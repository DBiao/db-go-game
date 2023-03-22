package controller

import (
	"db-go-game/pkg/common/dgin"
	"db-go-game/pkg/dhttp"
	"db-go-game/services/api/internal/dto"
	"db-go-game/services/api/internal/service"
	"github.com/gin-gonic/gin"
)

type IAuthController interface {
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	SignOut(ctx *gin.Context)
}

type authController struct {
	authService service.IAuthService
}

func NewAuthController(AuthService service.IAuthService) IAuthController {
	return &authController{authService: AuthService}
}

func (u *authController) SignUp(ctx *gin.Context) {
	var req *dto.SignUpReq
	if err := dgin.BindJSON(ctx, req); err != nil {
		dhttp.Error(ctx, dhttp.ERROR_CODE_HTTP_REQ_PARAM_ERR, err.Error())
		return
	}

	resp := u.authService.SignUp(req)
	dhttp.Success(ctx, resp)
}

func (u *authController) SignIn(ctx *gin.Context) {
	var req *dto.SignInReq
	if err := dgin.BindJSON(ctx, req); err != nil {
		dhttp.Error(ctx, dhttp.ERROR_CODE_HTTP_REQ_PARAM_ERR, err.Error())
		return
	}

	resp := u.authService.SignIn(req)
	dhttp.Success(ctx, resp)
}

func (u *authController) RefreshToken(ctx *gin.Context) {
	var token string
	resp := u.authService.RefreshToken(token)
	dhttp.Success(ctx, resp)
}

func (u *authController) SignOut(ctx *gin.Context) {
	var token string
	resp := u.authService.SignOut(token)
	dhttp.Success(ctx, resp)
}
