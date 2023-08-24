package api

import (
	"github.com/gin-gonic/gin"
	"mybili/service"
)

func CreateComment(c *gin.Context) {
	service := service.CreateCommentService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.Create()
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// ListComment 根据视频ID获取评论功能
func ListComment(c *gin.Context) {
	service := service.ListCommentService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.List()
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

func DeleteComments(c *gin.Context) {
	service := service.DeleteCommentService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.Delete(service.CommentID)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// AddLikes 给评论点赞
func AddCommentLikes(c *gin.Context) {
	service := service.CommentLikesService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.AddLikes()
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// AddCommentDislikes 给评论点赞
func AddCommentDislikes(c *gin.Context) {
	service := service.CommentLikesService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.AddDislikes()
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}
