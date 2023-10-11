package service

import (
	"fmt"
	"math/rand"
	"mybili/cache"
	"mybili/model"
	"mybili/serializer"
	"mybili/utils"
)

// CreateVideoService 视频投稿的服务
type CreateVideoService struct {
	Title     string  `form:"title" json:"title" binding:"required,min=2,max=30"`
	Info      string  `form:"info" json:"info" binding:"max=100"`
	VideoKey  string  `form:"video_key" json:"video_key"`
	AvatarKey string  `form:"avatar_key" json:"avatar_key"`
	Duration  float32 `form:"duration" json:"duration"  binding:"required,min=0.0"`
	Creator   string  `form:"creator" json:"creator"  binding:"required"`
}

// Create 创建视频
func (service *CreateVideoService) Create() serializer.Response {
	video := model.Video{
		Title:     service.Title,
		Info:      service.Info,
		VideoKey:  service.VideoKey,
		AvatarKey: service.AvatarKey,
		Duration:  service.Duration,
		Creator:   service.Creator,
	}

	err := model.DB.Create(&video).Error
	if err != nil {
		return serializer.Response{
			Code:  utils.ERROR_VIDEO_CREATING,
			Msg:   utils.GetErrMsg(utils.ERROR_VIDEO_CREATING),
			Error: err.Error(),
		}
	}

	//创建视频时把视频的评论数初始化
	// 第三个参数代表key的过期时间，0代表不会过期。
	err = cache.RedisClient.SetNX(cache.VideoCommentKey(video.ID), "0", 0).Err()
	if err != nil {
		return serializer.Response{
			Code:  utils.ERROR_VIDEO_CREATING,
			Msg:   utils.GetErrMsg(utils.ERROR_VIDEO_CREATING),
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code: utils.SUCCESS,
		Msg:  utils.GetErrMsg(utils.SUCCESS),
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
			Code:  utils.ERROR_VIDEO_NOEXIST,
			Msg:   utils.GetErrMsg(utils.ERROR_VIDEO_NOEXIST),
			Error: err.Error(),
		}
	}

	//处理视频被观看的一系问题
	video.AddView()

	return serializer.Response{
		Code: utils.SUCCESS,
		Msg:  utils.GetErrMsg(utils.SUCCESS),
		Data: serializer.BuildVideo(video),
	}
}

// ListVideoService 视频列表服务
type ListVideoService struct {
	Limit    int    `form:"limit"`
	Start    int    `form:"start"`
	UserName string `form:"username"`
}

// List 视频列表
func (service *ListVideoService) List() serializer.Response {
	videos := []model.Video{}
	var total int64

	//Count 用于获取匹配的记录数
	if err := model.DB.Model(model.Video{}).Count(&total).Error; err != nil {
		return serializer.Response{
			Code:  utils.DB_CONNECT_FAILED,
			Msg:   utils.GetErrMsg(utils.DB_CONNECT_FAILED),
			Error: err.Error(),
		}
	}

	if service.Limit >= int(total) {
		service.Limit = int(total)
	}

	ranOffset := rand.Intn(int(total) - service.Limit + 1)

	if err := model.DB.Limit(service.Limit).Offset(ranOffset).Find(&videos).Error; err != nil {
		return serializer.Response{
			Code:  utils.DB_CONNECT_FAILED,
			Msg:   utils.GetErrMsg(utils.DB_CONNECT_FAILED),
			Error: err.Error(),
		}
	}

	return serializer.BuildListResponse(serializer.BuildVideos(videos), uint(total))
}

func (service *ListVideoService) ListByUser(user_name string) serializer.Response {
	videos := []model.Video{}
	var total int64

	if err := model.DB.Model(model.Video{}).Where("creator = ?", user_name).Count(&total).Error; err != nil {
		return serializer.Response{
			Code:  utils.DB_CONNECT_FAILED,
			Msg:   utils.GetErrMsg(utils.DB_CONNECT_FAILED),
			Error: err.Error(),
		}
	}

	if err := model.DB.Where("creator = ?", user_name).Find(&videos).Error; err != nil {
		return serializer.Response{
			Code:  utils.DB_CONNECT_FAILED,
			Msg:   utils.GetErrMsg(utils.DB_CONNECT_FAILED),
			Error: err.Error(),
		}
	}

	return serializer.BuildListResponse(serializer.BuildVideos(videos), uint(total))
}

// UpdateVideoService 更新视频的服务
type UpdateVideoService struct {
	Title     string  `form:"title" json:"title" binding:"required,min=2,max=100"`
	Info      string  `form:"info" json:"info" binding:"max=300"`
	VideoKey  string  `form:"video_key" json:"video_key"`
	AvatarKey string  `form:"avatar_key" json:"avatar_key"`
	Duration  float32 `form:"duration" json:"duration"  binding:"required,min=0.0"`
}

// Update 更新视频
func (service *UpdateVideoService) Update(id string) serializer.Response {
	var video model.Video
	err := model.DB.First(&video, id).Error
	if err != nil {
		return serializer.Response{
			Code: utils.ERROR_VIDEO_NOEXIST,
			Msg:  utils.GetErrMsg(utils.ERROR_VIDEO_NOEXIST),
		}
	}

	video.Title = service.Title
	video.Info = service.Info
	video.VideoKey = service.VideoKey
	video.AvatarKey = service.AvatarKey
	video.Duration = service.Duration

	err = model.DB.Save(&video).Error
	if err != nil {
		return serializer.Response{
			Code: utils.ERROR_VIDEO_SAVE_FAILED,
			Msg:  utils.GetErrMsg(utils.ERROR_VIDEO_SAVE_FAILED),
		}
	}

	return serializer.Response{
		Code: utils.SUCCESS,
		Msg:  utils.GetErrMsg(utils.SUCCESS),
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
			Code:  utils.ERROR_VIDEO_NOEXIST,
			Msg:   utils.GetErrMsg(utils.ERROR_VIDEO_NOEXIST),
			Error: err.Error(),
		}
	}

	videoID := video.ID
	err = model.DB.Delete(&video).Error
	if err != nil {
		return serializer.Response{
			Code:  utils.ERROR_VIDEO_DELETE_FAILED,
			Msg:   utils.GetErrMsg(utils.ERROR_VIDEO_DELETE_FAILED),
			Error: err.Error(),
		}
	}

	var comments []model.Comment //查询视频相关的评论
	db := model.DB.Where("video_id = ?", videoID).Find(&comments)
	if err = db.Error; err != nil {
		return serializer.Response{
			Code:  utils.ERROR_VIDEO_DELETE_FAILED,
			Msg:   utils.GetErrMsg(utils.ERROR_VIDEO_DELETE_FAILED),
			Error: err.Error(),
		}
	}

	if db.RowsAffected > 0 { //删除视频相关的评论(如果有的话)
		//删除相关评论的redis数据库数据
		for _, comment := range comments {
			if comment.ParentId == 0 { //找出根评论，删除根评论及其子评论
				ser := DeleteCommentService{}
				ser.Delete(comment.ID)
			}
		}

		db = model.DB.Where("video_id = ?", videoID).Delete(&model.Comment{})
		if err = db.Error; err != nil {
			return serializer.Response{
				Code:  utils.ERROR_VIDEO_DELETE_FAILED,
				Msg:   utils.GetErrMsg(utils.ERROR_VIDEO_DELETE_FAILED),
				Error: err.Error(),
			}
		}
	}

	//删除redis数据库相关的数据
	//删除视频的播放量数据
	video.DeletedVideoViewNum(videoID)
	if err != nil {
		return serializer.Response{
			Code:  utils.REDIS_OPERATE_FAILED,
			Msg:   utils.GetErrMsg(utils.REDIS_OPERATE_FAILED),
			Error: err.Error(),
		}
	}

	//删除视频的点赞数据
	video.DeletedVideoLikesNum(videoID)
	if err != nil {
		return serializer.Response{
			Code:  utils.REDIS_OPERATE_FAILED,
			Msg:   utils.GetErrMsg(utils.REDIS_OPERATE_FAILED),
			Error: err.Error(),
		}
	}

	//删除视频的评论数量数据
	video.DeletedVideoCommentsNum(videoID)
	if err != nil {
		return serializer.Response{
			Code:  utils.REDIS_OPERATE_FAILED,
			Msg:   utils.GetErrMsg(utils.REDIS_OPERATE_FAILED),
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code: utils.SUCCESS,
		Msg:  utils.GetErrMsg(utils.SUCCESS),
	}
}

// VideoInfoService 视频点赞服务
type VideoInfoService struct {
	VideoID uint `form:"video_id" json:"video_id"`
	UserID  uint `form:"user_id" json:"user_id"`
}

func (service *VideoInfoService) AddLikes() serializer.Response {
	// 检测用户id是否包含在点赞集合中
	ok, err := cache.RedisClient.SIsMember(cache.VideoLikesKey(service.VideoID), service.UserID).Result()
	if err != nil {
		return serializer.Response{
			Code:  utils.REDIS_OPERATE_FAILED,
			Msg:   utils.GetErrMsg(utils.REDIS_OPERATE_FAILED),
			Error: err.Error(),
		}
	}

	var num int64
	if ok { // 用户重复点赞视为取消点赞
		//删除集合中的元素，返回被删元素的个数
		num, err = cache.RedisClient.SRem(cache.VideoLikesKey(service.VideoID), service.UserID).Result()

	} else { // 增加视频点赞数
		//添加元素，支持批量，返回添加成功的个数
		num, err = cache.RedisClient.SAdd(cache.VideoLikesKey(service.VideoID), service.UserID).Result()
	}
	if err != nil {
		return serializer.Response{
			Code:  utils.REDIS_OPERATE_FAILED,
			Msg:   utils.GetErrMsg(utils.REDIS_OPERATE_FAILED),
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code: utils.SUCCESS,
		Msg:  utils.GetErrMsg(utils.SUCCESS),
		Data: num,
	}
}

func (service *VideoInfoService) QueryLikes() serializer.Response {
	// 检测用户id是否包含在点赞集合中
	ok, err := cache.RedisClient.SIsMember(cache.VideoLikesKey(service.VideoID), service.UserID).Result()
	if err != nil {
		fmt.Println(err)
	}

	return serializer.Response{
		Code: utils.SUCCESS,
		Msg:  utils.GetErrMsg(utils.SUCCESS),
		Data: ok,
	}
}

func (service *VideoInfoService) QueryInfo() serializer.Response {
	// 检测UserID是否包含在点赞集合中
	ok, err := cache.RedisClient.SIsMember(cache.VideoLikesKey(service.VideoID), service.UserID).Result()
	if err != nil {
		fmt.Println(err)
	}

	return serializer.Response{
		Code: utils.SUCCESS,
		Msg:  utils.GetErrMsg(utils.SUCCESS),
		Data: struct { // 定义部分
			Liked bool
		}{ // 值初始化部分
			Liked: ok,
		},
	}
}
