package api

import (
	"github.com/gin-gonic/gin"
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
	return serializer.Response{
		Code:  utils.PARAM_ERROR,
		Msg:   utils.GetErrMsg(utils.PARAM_ERROR),
		Error: err.Error(),
	}
}
