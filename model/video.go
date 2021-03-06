package model

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"os"
	"time"
)

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

	name := video.Avatar
	ctx := context.Background()

	// 获取预签名URL
	presignedURL, err := client.Object.GetPresignedURL(ctx, http.MethodGet, name, ak, sk, time.Hour, nil)
	if err != nil {
		panic(err)
	}
	return presignedURL.String()
}

//视频模型
type Video struct {
	gorm.Model
	Title  string
	Info   string
	URL    string
	Avatar string
	View   uint64
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

	name := video.URL
	ctx := context.Background()

	// 获取预签名URL
	presignedURL, err := client.Object.GetPresignedURL(ctx, http.MethodGet, name, ak, sk, time.Hour, nil)
	if err != nil {
		panic(err)
	}
	return presignedURL.String()
}