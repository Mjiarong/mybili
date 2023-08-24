package service

import (
	"github.com/gin-gonic/gin"
	"mybili/middleware"
	"mybili/model"
	"mybili/serializer"
	"mybili/utils"
)

// UserLoginService 管理用户登录的服务
type UserLoginService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=5,max=30"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=40"`
}

// 设置token
func (service *UserLoginService) setToken(UserName string) (string, int) {
	token, code := middleware.SetToken(UserName)
	if code != utils.SUCCESS {
		return utils.GetErrMsg(code), code
	}
	return token, code
}

// Login 用户登录函数
func (service *UserLoginService) Login(c *gin.Context) serializer.Response {
	var user model.User

	//Check UserName
	if err := model.DB.Where("user_name = ?", service.UserName).First(&user).Error; err != nil {
		return serializer.Response{
			Code: utils.ACCOUNT_INCORRECT,
			Msg:  utils.GetErrMsg(utils.ACCOUNT_INCORRECT),
		}
	}

	if user.CheckPassword(service.Password) == false {
		return serializer.Response{
			Code: utils.PASSWORD_INCORRECT,
			Msg:  utils.GetErrMsg(utils.PASSWORD_INCORRECT),
		}
	}

	// 设置token
	token, code := service.setToken(user.UserName)
	return serializer.BuildUserResponse(user, code, token)
}
