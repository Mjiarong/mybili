
package service

import (
	"mybili/model"
	"mybili/serializer"
	"mybili/util/errmsg"
)

// CreateVideoService 视频投稿的服务
type CreateVideoService struct {
	Title  string `form:"title" json:"title" binding:"required,min=2,max=100"`
	Info   string `form:"info" json:"info" binding:"max=3000"`
	URL    string `form:"url" json:"url"`
	Avatar string `form:"avatar" json:"avatar"`
}

// Create 创建视频
func (service *CreateVideoService) Create(user *model.User) serializer.Response {
	video := model.Video{
		Title:  service.Title,
		Info:   service.Info,
		URL:    service.URL,
		Avatar: service.Avatar,
		UserID:user.ID,
	}

	err := model.DB.Create(&video).Error
	if err != nil {
		return serializer.Response{
			Code: errmsg.ERROR_VIDEO_CREATING,
			Msg:  errmsg.GetErrMsg(errmsg.ERROR_VIDEO_CREATING),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Code: errmsg.SUCCESS,
		Msg:  errmsg.GetErrMsg(errmsg.SUCCESS),
		Data: serializer.BuildVideo(video),
	}
}


// ShowVideoService 投稿详情的服务
type ShowVideoService struct {
}

// Show 视频
func (service *ShowVideoService) Show(id string) serializer.Response {
	var video model.Video
	err := model.DB.First(&video, id).Error
	if err != nil {
		return serializer.Response{
			Code: errmsg.ERROR_VIDEO_NOEXIST,
			Msg:  errmsg.GetErrMsg(errmsg.ERROR_VIDEO_NOEXIST),
			Error:  err.Error(),
		}
	}

	//处理视频被观看的一系问题
	video.AddView()

	return serializer.Response{
		Data: serializer.BuildVideo(video),
	}
}


// ListVideoService 视频列表服务
type ListVideoService struct {
	Limit int `form:"limit"`
	Start int `form:"start"`
}

// List 视频列表
func (service *ListVideoService) List() serializer.Response {
	videos := []model.Video{}
	var total int64

	if service.Limit == 0 {
		service.Limit = 6
	}

	if err := model.DB.Model(model.Video{}).Count(&total).Error; err != nil {
		return serializer.Response{
			Code: errmsg.DB_CONNECT_FAILED,
			Msg:  errmsg.GetErrMsg(errmsg.DB_CONNECT_FAILED),
			Error:  err.Error(),
		}
	}

	if err := model.DB.Limit(service.Limit).Offset(service.Start).Find(&videos).Error; err != nil {
		return serializer.Response{
			Code: errmsg.DB_CONNECT_FAILED,
			Msg:  errmsg.GetErrMsg(errmsg.DB_CONNECT_FAILED),
			Error:  err.Error(),
		}
	}

	return serializer.BuildListResponse(serializer.BuildVideos(videos), uint(total))
}

// UpdateVideoService 更新视频的服务
type UpdateVideoService struct {
	Title string `form:"title" json:"title" binding:"required,min=2,max=30"`
	Info  string `form:"info" json:"info" binding:"max=300"`
}

// Update 更新视频
func (service *UpdateVideoService) Update(id string) serializer.Response {
	var video model.Video
	err := model.DB.First(&video, id).Error
	if err != nil {
		return serializer.Response{
			Code: errmsg.ERROR_VIDEO_NOEXIST,
			Msg:  errmsg.GetErrMsg(errmsg.ERROR_VIDEO_NOEXIST),
		}
	}

	video.Title = service.Title
	video.Info = service.Info
	err = model.DB.Save(&video).Error
	if err != nil {
		return serializer.Response{
			Code: errmsg.ERROR_VIDEO_SAVE_FAILED,
			Msg:  errmsg.GetErrMsg(errmsg.ERROR_VIDEO_SAVE_FAILED),
		}
	}

	return serializer.Response{
		Code: errmsg.SUCCESS,
		Msg:  errmsg.GetErrMsg(errmsg.SUCCESS),
		Data: serializer.BuildVideo(video),
	}
}

// DeleteVideoService 删除投稿的服务
type DeleteVideoService struct {
}

// Delete 删除视频
func (service *DeleteVideoService) Delete(id string) serializer.Response {
	var video model.Video
	err := model.DB.First(&video, id).Error
	if err != nil {
		return serializer.Response{
			Code: errmsg.ERROR_VIDEO_NOEXIST,
			Msg:  errmsg.GetErrMsg(errmsg.ERROR_VIDEO_NOEXIST),
			Error:  err.Error(),
		}
	}

	err = model.DB.Delete(&video).Error
	if err != nil {
		return serializer.Response{
			Code: errmsg.ERROR_VIDEO_DELETE_FAILED,
			Msg:  errmsg.GetErrMsg(errmsg.ERROR_VIDEO_DELETE_FAILED),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Code: errmsg.SUCCESS,
		Msg:  errmsg.GetErrMsg(errmsg.SUCCESS),
	}
}