package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
	"mybili/serializer"
	"net/http"
	"net/url"
	"os"
	"time"
)

// UploadTokenService 上传oss token的服务
type UploadTokenService struct {
	Filename string `form:"filename" json:"filename"`
	Type string `form:"type" json:"type"`
}


// Post 创建签名URL
func (service *UploadTokenService) Post() serializer.Response {
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

	name := service.Type+"/"+ uuid.Must(uuid.NewRandom()).String()+"-"+ service.Filename
	ctx := context.Background()

	// 获取预签名上传URL
	presignedPutURL, err := client.Object.GetPresignedURL(ctx, http.MethodPut, name, ak, sk, time.Hour, nil)
	if err != nil {
		panic(err)
	}

	// 获取预签名下载URL
	presignedGetURL, err := client.Object.GetPresignedURL(ctx, http.MethodGet, name, ak, sk, time.Hour, nil)
	if err != nil {
		panic(err)
	}
	/*
	// 2. 通过预签名方式上传对象
	data := "test upload with presignedURL"
	f:= strings.NewReader(data)
	_, err = http.NewRequest(http.MethodPut, presignedURL.String(), f)
	if err != nil {
		panic(err)
	}*/

	return serializer.Response{
		Data: map[string]string{
		"key": name,
		"signedPutURL": presignedPutURL.String(),
		"signedGetURL": presignedGetURL.String(),
		},
	}
}

