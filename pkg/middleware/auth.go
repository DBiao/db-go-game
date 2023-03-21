package middleware

import (
	"db-go-game/pkg/common/djwt"
	"db-go-game/pkg/common/dlog"
	"db-go-game/pkg/common/dredis"
	"db-go-game/pkg/constant"
	"db-go-game/pkg/dhttp"
	"db-go-game/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"strings"
)

func JwtAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			err          error
			token        *jwt.Token
			ok           bool
			uid          interface{}
			sessionId    interface{}
			sessionIdVal string
			sessionIdKey string
		)
		token, err = djwt.ParseFromCookie(ctx)
		if err != nil {
			ctx.Abort()
			dhttp.Error(ctx, dhttp.ERROR_CODE_HTTP_JWT_TOKEN_ERR, err.Error())
			return
		}
		claims := jwt.MapClaims{}
		for key, value := range token.Claims.(jwt.MapClaims) {
			claims[key] = value
		}
		if uid, ok = claims[constant.USER_UID]; ok == false {
			ctx.Abort()
			dhttp.Error(ctx, dhttp.ERROR_CODE_HTTP_USER_ID_DOESNOT_EXIST, dhttp.ERROR_HTTP_USER_ID_DOESNOT_EXIST)
			return
		}

		ctx.Set(constant.USER_UID, uid)
		if strings.HasPrefix(ctx.FullPath(), constant.API_PUBLIC) && ctx.Request.Method == constant.HTTP_REQUEST_METHOD_GET {
			return
		}
		if sessionId, ok = claims[constant.USER_JWT_SESSION_ID]; ok == false {
			ctx.Abort()
			dhttp.Error(ctx, dhttp.ERROR_CODE_HTTP_JWT_TOKEN_UUID_DOESNOT_EXIST, dhttp.ERROR_HTTP_JWT_TOKEN_UUID_DOESNOT_EXIST)
			return
		}
		sessionIdKey = constant.RK_SYNC_USER_ACCESS_TOKEN_SESSION_ID + utils.ToString(uid)
		if sessionIdVal, err = dredis.Get(sessionIdKey); err != nil {
			ctx.Abort()
			dhttp.Error(ctx, dhttp.ERROR_CODE_HTTP_TOKEN_AUTHENTICATION_FAILED, dhttp.ERROR_HTTP_TOKEN_AUTHENTICATION_FAILED)
			dlog.Warn(dhttp.ERROR_CODE_HTTP_TOKEN_AUTHENTICATION_FAILED, dhttp.ERROR_HTTP_TOKEN_AUTHENTICATION_FAILED, err.Error())
			return
		}
		if sessionIdVal != utils.ToString(sessionId) {
			ctx.Abort()
			dhttp.Error(ctx, dhttp.ERROR_CODE_HTTP_TOKEN_AUTHENTICATION_FAILED, dhttp.ERROR_HTTP_TOKEN_AUTHENTICATION_FAILED)
			return
		}
	}
}
