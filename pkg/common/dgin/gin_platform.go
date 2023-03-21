package dgin

import (
	"db-go-game/pkg/constant"
	"db-go-game/pkg/dhttp"
	"db-go-game/pkg/utils"
	"github.com/gin-gonic/gin"
)

func GetPlatform(ctx *gin.Context) (platform int32) {
	var (
		value  any
		exists bool
	)
	value, exists = ctx.Get(constant.USER_PLATFORM)
	if exists == false {
		dhttp.Error(ctx, dhttp.ERROR_CODE_HTTP_PLATFORM_DOESNOT_EXIST, dhttp.ERROR_HTTP_PLATFORM_DOESNOT_EXIST)
		return
	}
	platform, _ = utils.ToInt32(value)
	if platform == 0 {
		dhttp.Error(ctx, dhttp.ERROR_CODE_HTTP_PLATFORM_DOESNOT_EXIST, dhttp.ERROR_HTTP_PLATFORM_DOESNOT_EXIST)
		return
	}
	return
}
