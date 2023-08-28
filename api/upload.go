package api

import (
	"github.com/gin-gonic/gin"
	"mybili/service"
	"net/http"
)

// UploadToken 上传授权
func UploadToken(c *gin.Context) {
	service := service.UploadTokenService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.Post()
		c.JSON(200, res)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse(err))
	}
}

// 获取临时秘钥
func TmpCredentials(c *gin.Context) {
	service := service.TmpCredentialsService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.GetCredential()
		c.JSON(200, res)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse(err))
	}
}
