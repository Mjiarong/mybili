package api

import (
	"github.com/gin-gonic/gin"
	"mybili/service"
	"net/http"
)

// CreateVideo 视频投稿
func CreateVideo(c *gin.Context) {
	//username, ok := CurrentUser(c)
	//user, _ := model.GetUserByName(username)
	service := service.CreateVideoService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.Create()
		c.JSON(200, res)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse(err))
	}
}

// ShowVideo 视频详情接口
func ShowVideo(c *gin.Context) {
	service := service.ShowVideoService{}
	res := service.Show(c.Param("id"))
	c.JSON(200, res)
}

// ListVideo 视频列表接口
func ListVideo(c *gin.Context) {
	service := service.ListVideoService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.List()
		c.JSON(200, res)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse(err))
	}
}

// ListVideo 根据用户名查询视频列表接口
func ListVideosOfUser(c *gin.Context) {
	service := service.ListVideoService{}
	res := service.ListByUser(c.Param("username"))
	c.JSON(200, res)
}

// UpdateVideo 更新视频的接口
func UpdateVideo(c *gin.Context) {
	service := service.UpdateVideoService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.Update(c.Param("id"))
		c.JSON(200, res)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse(err))
	}
}

// DeleteVideo 删除视频的接口
func DeleteVideo(c *gin.Context) {
	service := service.DeleteVideoService{}
	res := service.Delete(c.Param("id"))
	c.JSON(200, res)
}

func AddVideoLikes(c *gin.Context) {
	service := service.VideoInfoService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.AddLikes()
		c.JSON(200, res)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse(err))
	}
}

// QueryPlayInfo 查询用户视频点赞、收藏等操作信息
func QueryPlayInfo(c *gin.Context) {
	service := service.VideoInfoService{}
	//ShouldBind
	//如果是 GET 请求，只使用 Form 绑定引擎（query）。
	if err := c.ShouldBind(&service); err == nil {
		res := service.QueryInfo()
		c.JSON(200, res)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse(err))
	}
}
