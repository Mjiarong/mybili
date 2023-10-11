package model

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"gorm.io/gorm"
	"mybili/cache"
	"mybili/utils"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// 视频模型
type Video struct {
	gorm.Model
	Title     string  `gorm:"size:80;not null"`
	Info      string  `gorm:"size:1000"`
	VideoKey  string  `gorm:"size:256;not null"`
	AvatarKey string  `gorm:"size:256;not null"`
	Duration  float32 `gorm:"not null"`
	Creator   string  `gorm:"size:30;not null"`
}

// AvatarURL 获取带签名的封面地址
func (video *Video) AvatarURL() string {
	u, _ := url.Parse(os.Getenv("BUCKET_ADDR"))
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId
			SecretID: os.Getenv("SECRET_ID"),
			// 环境变量 SECRETKEY 表示用户的 SecretKey
			SecretKey: os.Getenv("SECRET_KEY"),
		},
	})

	ak := os.Getenv("SECRET_ID")
	sk := os.Getenv("SECRET_KEY")

	name := video.AvatarKey
	ctx := context.Background()

	// 获取预签名URL
	presignedURL, err := client.Object.GetPresignedURL(ctx, http.MethodGet, name, ak, sk, time.Hour, nil)
	if err != nil {
		utils.Logger.Errorln(err)
	}
	return presignedURL.String()
}

// VideoURL 获取带签名的视频地址
func (video *Video) VideoURL() string {
	u, _ := url.Parse(os.Getenv("BUCKET_ADDR"))
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId
			SecretID: os.Getenv("SECRET_ID"),
			// 环境变量 SECRETKEY 表示用户的 SecretKey
			SecretKey: os.Getenv("SECRET_KEY"),
		},
	})

	ak := os.Getenv("SECRET_ID")
	sk := os.Getenv("SECRET_KEY")

	name := video.VideoKey
	ctx := context.Background()

	// 获取预签名URL
	presignedURL, err := client.Object.GetPresignedURL(ctx, http.MethodGet, name, ak, sk, time.Hour, nil)
	if err != nil {
		utils.Logger.Errorln(err)
	}
	return presignedURL.String()
}

// View 获取点击数
func (video *Video) View() uint64 {
	countStr, err := cache.RedisClient.Get(cache.VideoViewKey(video.ID)).Result()
	if err != nil {
		utils.Logger.Errorln(err)
	}
	count, _ := strconv.ParseUint(countStr, 10, 64)
	return count
}

// AddView 视频浏览时浏览数+1
func (video *Video) AddView() {
	// 增加视频点击数
	cache.RedisClient.Incr(cache.VideoViewKey(video.ID))
	// 增加排行点击数
	cache.RedisClient.ZIncrBy(cache.DailyRankKey, 1, strconv.Itoa(int(video.ID)))
}

func (video *Video) GetCommentsNum() uint64 {
	countStr, err := cache.RedisClient.Get(cache.VideoCommentKey(video.ID)).Result()
	if err != nil {
		utils.Logger.Errorln(err)
	}
	count, _ := strconv.ParseUint(countStr, 10, 64)
	return count
}

func (video *Video) GetLikesNum() int64 {
	count, err := cache.RedisClient.SCard(cache.VideoLikesKey(video.ID)).Result()
	if err != nil {
		utils.Logger.Errorln(err)
	}
	return count
}

// DeletedVideoLikes 删除视频的播放数量数据
func (video *Video) DeletedVideoViewNum(id uint) error {
	err := cache.RedisClient.Del(cache.VideoViewKey(id)).Err()
	if err != nil {
		utils.Logger.Errorln(err)
	}
	return err
}

// DeletedVideoLikes 删除视频的点赞数据
func (video *Video) DeletedVideoLikesNum(id uint) error {
	err := cache.RedisClient.Del(cache.VideoLikesKey(id)).Err()
	if err != nil {
		utils.Logger.Errorln(err)
	}
	return err
}

// DeletedVideoCommentsNum 删除视频的评论数数据
func (video *Video) DeletedVideoCommentsNum(id uint) error {
	err := cache.RedisClient.Del(cache.VideoCommentKey(id)).Err()
	if err != nil {
		utils.Logger.Errorln(err)
	}
	return err
}
