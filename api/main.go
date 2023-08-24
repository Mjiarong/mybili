package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	validator "gopkg.in/go-playground/validator.v8"
	"mybili/conf"
	"mybili/serializer"
	"mybili/utils"
)

// Ping 状态检查页面
func Ping(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Code: 0,
		Msg:  "Pong",
	})
}

// CurrentUser 获取当前用户
func CurrentUser(c *gin.Context) (string, bool) {
	username, ok := c.Get("username")
	return username.(string), ok
}

// ErrorResponse 返回错误消息
func ErrorResponse(err error) serializer.Response {
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			field := conf.T(fmt.Sprintf("Field.%s", e.Field))
			tag := conf.T(fmt.Sprintf("Tag.Valid.%s", e.Tag))
			return serializer.Response{
				Code:  utils.VALIDATION_ERROR,
				Msg:   fmt.Sprintf("%s%s", field, tag),
				Error: err.Error(),
			}
		}
	}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return serializer.Response{
			Code:  utils.UNMARSHAL_TYPE_ERROR,
			Msg:   utils.GetErrMsg(utils.UNMARSHAL_TYPE_ERROR),
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code:  utils.PARAM_ERROR,
		Msg:   utils.GetErrMsg(utils.PARAM_ERROR),
		Error: err.Error(),
	}
}
