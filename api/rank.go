package api

import (
	"github.com/gin-gonic/gin"
	"mybili/service"
	"net/http"
)

// DailyRank 每日排行
func DailyRank(c *gin.Context) {
	service := service.DailyRankService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.Get()
		c.JSON(200, res)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse(err))
	}
}
