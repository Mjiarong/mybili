package service

import (
	"fmt"
	"mybili/cache"
	"mybili/model"
	"mybili/serializer"
	"mybili/utils"
)

type CreateCommentService struct {
	//结构体成员必须以大写开头
	Content       string `form:"content" json:"content" binding:"required,max=300"`
	UserId        uint   `form:"user_id" json:"user_id"  binding:"required"`
	UserName      string `form:"user_name" json:"user_name"  binding:"required"`
	Nickname      string `form:"nickname" json:"nickname"  binding:"required"`
	UserAvatarKey string `form:"user_avatar_key" json:"user_avatar_key"`
	VideoId       uint   `form:"video_id" json:"video_id"  binding:"required"`
	ParentId      uint   `form:"parent_id" json:"parent_id"`
	ReplyUserName string `form:"reply_user_name" json:"reply_user_name"`
}

// Create 创建评论
func (service *CreateCommentService) Create() serializer.Response {
	comment := model.Comment{
		Content:       service.Content,
		UserId:        service.UserId,
		UserName:      service.UserName,
		Nickname:      service.Nickname,
		UserAvatarKey: service.UserAvatarKey,
		VideoId:       service.VideoId,
		ParentId:      service.ParentId,
		ReplyUserName: service.ReplyUserName,
	}

	err := model.DB.Create(&comment).Error
	if err != nil {
		return serializer.Response{
			Code:  utils.ERROR_COMMENT_CREATING,
			Msg:   utils.GetErrMsg(utils.ERROR_COMMENT_CREATING),
			Error: err.Error(),
		}
	}

	//对应的视频评论数+1
	comment.AddCommentsNum()

	return serializer.Response{
		Code: utils.SUCCESS,
		Msg:  utils.GetErrMsg(utils.SUCCESS),
		Data: serializer.BuildComment(comment),
	}
}

type ListCommentService struct {
	VideoID uint `form:"video_id" json:"video_id"`
	UserID  uint `form:"user_id" json:"user_id"`
}

// List 评论列表
func (service *ListCommentService) List() serializer.Response {
	comments := []model.Comment{}

	if err := model.DB.Where("video_id = ?", service.VideoID).Find(&comments).Error; err != nil {
		return serializer.Response{
			Code:  utils.DB_CONNECT_FAILED,
			Msg:   utils.GetErrMsg(utils.DB_CONNECT_FAILED),
			Error: err.Error(),
		}
	}

	items, total := serializer.BuildComments(comments, service.UserID)

	return serializer.BuildListResponse(items, total)
}

// CommentLikesService 评论点赞服务
type CommentLikesService struct {
	CommentID uint `form:"comment_id"  json:"comment_id"`
	UserID    uint `form:"user_id" json:"user_id"`
}

func (service *CommentLikesService) AddLikes() serializer.Response {
	// 检测用户id是否已经包含在点踩集合中
	ok, err := cache.RedisClient.SIsMember(cache.CommentDislikesKey(service.CommentID), service.UserID).Result()
	if err != nil {
		return serializer.Response{
			Code:  utils.REDIS_OPERATE_FAILED,
			Msg:   utils.GetErrMsg(utils.REDIS_OPERATE_FAILED),
			Error: err.Error(),
		}
	}

	var num int64
	if ok { // 用户点赞视为取消点踩
		num, err = cache.RedisClient.SRem(cache.CommentDislikesKey(service.CommentID), service.UserID).Result()
	}

	// 检测用户id是否包含在点赞集合中
	ok, err = cache.RedisClient.SIsMember(cache.CommentLikesKey(service.CommentID), service.UserID).Result()
	if err != nil {
		return serializer.Response{
			Code:  utils.REDIS_OPERATE_FAILED,
			Msg:   utils.GetErrMsg(utils.REDIS_OPERATE_FAILED),
			Error: err.Error(),
		}
	}

	if ok { // 用户重复点赞视为取消点赞
		//删除集合中的元素，返回被删元素的个数
		num, err = cache.RedisClient.SRem(cache.CommentLikesKey(service.CommentID), service.UserID).Result()
	} else { // 增加评论点赞数
		//添加元素，支持批量，返回添加成功的个数
		num, err = cache.RedisClient.SAdd(cache.CommentLikesKey(service.CommentID), service.UserID).Result()
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

// 视频点踩
func (service *CommentLikesService) AddDislikes() serializer.Response {
	// 检测用户id是否包含在点赞集合中
	ok, err := cache.RedisClient.SIsMember(cache.CommentLikesKey(service.CommentID), service.UserID).Result()
	if err != nil {
		return serializer.Response{
			Code:  114,
			Msg:   utils.GetErrMsg(utils.REDIS_OPERATE_FAILED),
			Error: err.Error(),
		}
	}

	var num int64
	if ok { // 用户点踩视为取消点赞
		num, err = cache.RedisClient.SRem(cache.CommentLikesKey(service.CommentID), service.UserID).Result()
	}

	// 检测用户id是否已经包含在点踩集合中
	ok, err = cache.RedisClient.SIsMember(cache.CommentDislikesKey(service.CommentID), service.UserID).Result()
	if err != nil {
		return serializer.Response{
			Code:  514,
			Msg:   utils.GetErrMsg(utils.REDIS_OPERATE_FAILED),
			Error: err.Error(),
		}
	}

	if ok { // 用户重复点踩视为取消点踩
		num, err = cache.RedisClient.SRem(cache.CommentDislikesKey(service.CommentID), service.UserID).Result()
	} else { // 增加评论点踩数
		//添加元素，支持批量，返回添加成功的个数
		num, err = cache.RedisClient.SAdd(cache.CommentDislikesKey(service.CommentID), service.UserID).Result()
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

type DeleteCommentService struct {
	CommentID uint `form:"comment_id" json:"comment_id"`
	VideoID   uint `form:"video_id" json:"video_id"`
}

// Delete 删除评论
func (service *DeleteCommentService) Delete(CommentID uint) serializer.Response {
	var comment model.Comment
	err := model.DB.First(&comment, CommentID).Error
	if err != nil {
		return serializer.Response{
			Code:  utils.ERROR_COMMENT_NOEXIST,
			Msg:   utils.GetErrMsg(utils.ERROR_COMMENT_NOEXIST),
			Error: err.Error(),
		}
	}

	var isRootRep bool
	if comment.ParentId == 0 { //判断是不是根评论
		isRootRep = true
	} else {
		isRootRep = false
	}

	//删除评论
	err = model.DB.Delete(&comment).Error
	if err != nil {
		return serializer.Response{
			Code:  utils.ERROR_VIDEO_DELETE_FAILED,
			Msg:   utils.GetErrMsg(utils.ERROR_VIDEO_DELETE_FAILED),
			Error: err.Error(),
		}
	}

	var subRepNum int64 //子评论数目
	if isRootRep {      //删除相关的子评论(如果有的话)
		var comments []model.Comment
		db := model.DB.Where("parent_id = ?", CommentID).Find(&comments)
		if err = db.Error; err != nil {
			return serializer.Response{
				Code:  utils.ERROR_VIDEO_DELETE_FAILED,
				Msg:   utils.GetErrMsg(utils.ERROR_VIDEO_DELETE_FAILED),
				Error: err.Error(),
			}
		}
		//执行结果受影响的行数
		subRepNum = db.RowsAffected
		fmt.Printf("db.RowsAffected=%v\n", db.RowsAffected)
		fmt.Println(comments)

		db = model.DB.Where("parent_id = ?", CommentID).Delete(&model.Comment{})
		if err = db.Error; err != nil {
			return serializer.Response{
				Code:  utils.ERROR_VIDEO_DELETE_FAILED,
				Msg:   utils.GetErrMsg(utils.ERROR_VIDEO_DELETE_FAILED),
				Error: err.Error(),
			}
		}

		//删除相关子评论的redis数据库数据
		for _, subRep := range comments {
			comment.DeletedCommentsLikes(subRep.ID)
			comment.DeletedCommentsDislikes(subRep.ID)
		}

	}

	//删除redis数据库相关的数据
	comment.DecVideoCommentsNum(subRepNum + 1)
	if err != nil {
		return serializer.Response{
			Code:  utils.REDIS_OPERATE_FAILED,
			Msg:   utils.GetErrMsg(utils.REDIS_OPERATE_FAILED),
			Error: err.Error(),
		}
	}

	comment.DeletedCommentsLikes(CommentID)
	if err != nil {
		return serializer.Response{
			Code:  utils.REDIS_OPERATE_FAILED,
			Msg:   utils.GetErrMsg(utils.REDIS_OPERATE_FAILED),
			Error: err.Error(),
		}
	}

	comment.DeletedCommentsDislikes(CommentID)
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
