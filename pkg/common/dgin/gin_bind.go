package dgin

import (
	"db-go-game/pkg/dhttp"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindJSON(ctx *gin.Context, params interface{}) (err error) {
	if err = ctx.BindJSON(params); err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			if len(err.(validator.ValidationErrors)) > 0 {
				dhttp.Error(ctx, dhttp.ERROR_CODE_HTTP_REQ_PARAM_ERR, "parameter validation failed on the "+err.(validator.ValidationErrors)[0].StructField())
				return
			}
		}
		dhttp.Error(ctx, dhttp.ERROR_CODE_HTTP_REQ_PARAM_ERR, err.Error())
		return
	}
	return
}

func ShouldBindQuery(ctx *gin.Context, params interface{}) (err error) {
	if err = ctx.ShouldBindQuery(params); err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			if len(err.(validator.ValidationErrors)) > 0 {
				dhttp.Error(ctx, dhttp.ERROR_CODE_HTTP_REQ_PARAM_ERR, "parameter validation failed on the "+err.(validator.ValidationErrors)[0].StructField())
				return
			}
		}
		dhttp.Error(ctx, dhttp.ERROR_CODE_HTTP_REQ_PARAM_ERR, err.Error())
		return
	}
	return
}
