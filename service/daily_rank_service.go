package service

import (
	"fmt"
	"mybili/cache"
	"mybili/model"
	"mybili/serializer"
	"mybili/utils"
	"strings"
)

// DailyRankService 每日排行的服务
type DailyRankService struct {
}

// Get 获取排行
func (service *DailyRankService) Get() serializer.Response {
	var videos []model.Video

	// 从redis读取点击前十的视频
	vids, _ := cache.RedisClient.ZRevRange(cache.DailyRankKey, 0, 9).Result()
	if len(vids) >= 1 {
		order := fmt.Sprintf("FIELD(id, %s)", strings.Join(vids, ","))
		err := model.DB.Where("id in (?)", vids).Order(order).Find(&videos).Error
		if err != nil {
			return serializer.Response{
				Code:  utils.DB_CONNECT_FAILED,
				Msg:   utils.GetErrMsg(utils.DB_CONNECT_FAILED),
				Error: err.Error(),
			}
		}
	}

	return serializer.Response{
		Data: serializer.BuildVideos(videos),
	}
}
