package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"mybili/serializer"
	"mybili/service"
	"mybili/utils"
	"net/http"
)

// UserRegister 用户注册接口
func UserRegister(c *gin.Context) {
	var service service.UserRegisterService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Register()
		c.JSON(200, res)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse(err))
	}
}

// UserLogin 用户登录接口
func UserLogin(c *gin.Context) {
	var service service.UserLoginService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Login(c)
		c.JSON(200, res)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse(err))
	}
}

// UserInfo 查询用户接口
func UserInfo(c *gin.Context) {
	var service service.UserInfoService
	res := service.GetInfo(c.Param("name"))
	c.JSON(200, res)
}

// UserLogout 用户登出
func UserLogout(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
	s.Save()
	c.JSON(200, serializer.Response{
		Code: utils.SUCCESS,
		Msg:  utils.GetErrMsg(utils.SUCCESS),
	})
}
