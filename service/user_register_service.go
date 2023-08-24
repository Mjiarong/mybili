package service

import (
	"mybili/model"
	"mybili/serializer"
	"mybili/utils"
)

// UserRegisterService 管理用户注册服务
type UserRegisterService struct {
	Nickname        string `form:"nickname" json:"nickname" binding:"required,min=2,max=30"`
	UserName        string `form:"user_name" json:"user_name" binding:"required,min=5,max=30"`
	Password        string `form:"password" json:"password" binding:"required,min=8,max=40"`
	PasswordConfirm string `form:"password_confirm" json:"password_confirm" binding:"required,min=8,max=40"`
	Avatar          string `form:"avatar" json:"avatar"`
}

// valid 验证表单
func (service *UserRegisterService) valid() *serializer.Response {
	if service.PasswordConfirm != service.Password {
		return &serializer.Response{
			Code: utils.PASSWORD_ENTERED_DIFFERENT,
			Msg:  utils.GetErrMsg(utils.PASSWORD_ENTERED_DIFFERENT),
		}
	}

	count := int64(0)
	model.DB.Model(&model.User{}).Where("nickname = ?", service.Nickname).Count(&count)
	if count > 0 {
		return &serializer.Response{
			Code: utils.NICKNAME_OCCUPIED,
			Msg:  utils.GetErrMsg(utils.NICKNAME_OCCUPIED),
		}
	}

	count = 0
	model.DB.Model(&model.User{}).Where("user_name = ?", service.UserName).Count(&count)
	if count > 0 {
		return &serializer.Response{
			Code: utils.USERNAME_REGISTERED,
			Msg:  utils.GetErrMsg(utils.USERNAME_REGISTERED),
		}
	}

	return nil
}

// Register 用户注册
func (service *UserRegisterService) Register() serializer.Response {
	user := model.User{
		Nickname: service.Nickname,
		UserName: service.UserName,
		Status:   model.Active,
		Avatar:   service.Avatar,
	}

	// 验证用户名是否已经被注册
	if err := service.valid(); err != nil {
		return *err
	}

	// 加密密码
	code := utils.SUCCESS
	if err := user.SetPassword(service.Password); err != nil {
		code = utils.PASSWORD_ENCRYPT_FAILED
		return serializer.Response{
			Code:  code,
			Msg:   utils.GetErrMsg(code),
			Error: err.Error(),
		}
	}

	// 创建用户
	if err := model.DB.Create(&user).Error; err != nil {
		code = utils.USERNAME_REGISTERED
		return serializer.Response{
			Code:  code,
			Msg:   utils.GetErrMsg(code),
			Error: err.Error(),
		}
	}

	return serializer.BuildUserResponse(user, code, utils.GetErrMsg(code))
}
