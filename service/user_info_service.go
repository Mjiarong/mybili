package service

import (
	"mybili/model"
	"mybili/serializer"
	"mybili/utils"
)

// UserInfoService 管理查询用户信息服务
type UserInfoService struct {
}

// getInfo 查询用户信息函数
func (service *UserInfoService) GetInfo(name string) serializer.Response {
	var user model.User

	code := utils.SUCCESS
	//Check UserName
	if err := model.DB.Where("user_name = ?", name).First(&user).Error; err != nil {
		return serializer.Response{
			Code: utils.ACCOUNT_INCORRECT,
			Msg:  utils.GetErrMsg(utils.ACCOUNT_INCORRECT),
		}
	}

	return serializer.BuildUserResponse(user, code, utils.GetErrMsg(code))
}
