package service

import (
	"mybili/model"
	"mybili/serializer"
	"mybili/util/errmsg"
)

// UserRegisterService 管理用户注册服务
type UserRegisterService struct {
	Nickname        string `form:"nickname" json:"nickname" binding:"required,min=2,max=30"`
	UserName        string `form:"user_name" json:"user_name" binding:"required,min=5,max=30"`
	Password        string `form:"password" json:"password" binding:"required,min=8,max=40"`
	PasswordConfirm string `form:"password_confirm" json:"password_confirm" binding:"required,min=8,max=40"`
}

// valid 验证表单
func (service *UserRegisterService) valid() *serializer.Response {
	if service.PasswordConfirm != service.Password {
		return &serializer.Response{
			Code: errmsg.PASSWORD_ENTERED_DIFFERENT,
			Msg:  errmsg.GetErrMsg(errmsg.PASSWORD_ENTERED_DIFFERENT),
		}
	}

	count := int64(0)
	model.DB.Model(&model.User{}).Where("nickname = ?", service.Nickname).Count(&count)
	if count > 0 {
		return &serializer.Response{
			Code: errmsg.NICKNAME_OCCUPIED,
			Msg:  errmsg.GetErrMsg(errmsg.NICKNAME_OCCUPIED),
		}
	}

	count = 0
	model.DB.Model(&model.User{}).Where("user_name = ?", service.UserName).Count(&count)
	if count > 0 {
		return &serializer.Response{
			Code: errmsg.USERNAME_REGISTERED,
			Msg:  errmsg.GetErrMsg(errmsg.USERNAME_REGISTERED),
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
	}

	// 表单验证
	if err := service.valid(); err != nil {
		return *err
	}

	// 加密密码
	if err := user.SetPassword(service.Password); err != nil {
		return serializer.Response{
			Code: errmsg.USERNAME_REGISTERED,
			Msg:  errmsg.GetErrMsg(errmsg.USERNAME_REGISTERED),
			Error:err.Error(),
		}
	}

	// 创建用户
	if err := model.DB.Create(&user).Error; err != nil {
		return serializer.Response{
			Code: errmsg.USERNAME_REGISTERED,
			Msg:  errmsg.GetErrMsg(errmsg.USERNAME_REGISTERED),
			Error:err.Error(),
		}
	}

	return serializer.BuildUserResponse(user)
}
