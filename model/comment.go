package model

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"gorm.io/gorm"
	"mybili/cache"
	"mybili/utils"
	"net/http"
	"net/url"
	"os"
	"time"
)

// 视频评论模型
type Comment struct {
	gorm.Model
	Content       string `gorm:"size:500;not null"` //评论的内容
	UserName      string //发出该评论用户的名称
	Nickname      string //发出该评论用户的昵称
	UserAvatarKey string //发出该评论用户的头像
	ParentId      uint   //指向父评论的id,如果不是对评论的回复,那么该值为null
	ReplyUserName string `gorm:"size:30"`  //该评论@的用户的名称
	UserID        uint   `gorm:"not null"` //发出该评论用户的id
	User          User   //belongs to 会与另一个模型建立了一对一的连接
	VideoID       uint   `gorm:"not null"` //评论所对应的视频的id
	Video         Video
}

// AvatarURL 获取带签名的用户头像地址
func (com *Comment) AvatarURL() string {
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

	name := com.UserAvatarKey
	ctx := context.Background()

	// 获取预签名URL
	presignedURL, err := client.Object.GetPresignedURL(ctx, http.MethodGet, name, ak, sk, time.Hour, nil)
	if err != nil {
		utils.Logger.Errorf("err:%v", err)
		return ""
	}
	return presignedURL.String()
}

// AddCommentsNum 视频评论时评论数+1
func (com *Comment) AddCommentsNum() error {
	// 增加视频评论数
	_, err := cache.RedisClient.Incr(cache.VideoCommentKey(com.VideoID)).Result()
	if err != nil {
		utils.Logger.Errorf("err:%v", err)
	}
	return err
}

// DecCommentsNum 删除评论时评论数减少
func (com *Comment) DecVideoCommentsNum(num int64) error {
	// DecrBy函数，可以指定每次递减多少
	_, err := cache.RedisClient.DecrBy(cache.VideoCommentKey(com.VideoID), num).Result()
	if err != nil {
		utils.Logger.Errorf("err:%v", err)
	}
	return err
}

// DeletedCommentsLikes 删除评论点赞数据集合
func (com *Comment) DeletedCommentsLikes(id uint) error {
	err := cache.RedisClient.Del(cache.CommentLikesKey(id)).Err()
	if err != nil {
		utils.Logger.Errorf("err:%v", err)
	}
	return err
}

// DeletedCommentsDislikes 删除评论点踩数据集合
func (com *Comment) DeletedCommentsDislikes(id uint) error {
	err := cache.RedisClient.Del(cache.CommentDislikesKey(id)).Err()
	if err != nil {
		utils.Logger.Errorf("err:%v", err)
	}
	return err
}

// Likes 获取点赞数
func (com *Comment) Likes() int64 {
	count, err := cache.RedisClient.SCard(cache.CommentLikesKey(com.ID)).Result()
	if err != nil {
		utils.Logger.Errorf("err:%v", err)
	}
	return count
}

// Liked 获取点赞状态
func (com *Comment) Liked(UserID uint) bool {
	ok, err := cache.RedisClient.SIsMember(cache.CommentLikesKey(com.ID), UserID).Result()
	if err != nil {
		utils.Logger.Errorf("err:%v", err)
	}
	return ok
}

// Disliked 获取点踩状态
func (com *Comment) Disliked(UserID uint) bool {
	ok, err := cache.RedisClient.SIsMember(cache.CommentDislikesKey(com.ID), UserID).Result()
	if err != nil {
		fmt.Printf("Disliked err:%v\n", err)
	}
	return ok
}
