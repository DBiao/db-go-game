package dgin

import (
	"db-go-game/pkg/constant"
	"db-go-game/pkg/dhttp"
	"db-go-game/pkg/utils"
	"github.com/gin-gonic/gin"
)

func GetUid(ctx *gin.Context) (uid int64) {
	var (
		value  any
		exists bool
	)
	value, exists = ctx.Get(constant.USER_UID)
	if exists == false {
		dhttp.Error(ctx, dhttp.ERROR_CODE_HTTP_USER_ID_DOESNOT_EXIST, dhttp.ERROR_HTTP_USER_ID_DOESNOT_EXIST)
		return
	}
	uid, _ = utils.ToInt64(value)
	if uid == 0 {
		dhttp.Error(ctx, dhttp.ERROR_CODE_HTTP_USER_ID_DOESNOT_EXIST, dhttp.ERROR_HTTP_USER_ID_DOESNOT_EXIST)
		return
	}
	return
}
