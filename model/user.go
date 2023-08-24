package model

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"mybili/utils"
	"net/http"
	"net/url"
	"os"
	"time"
)

// User 用户模型
type User struct {
	gorm.Model
	UserName       string
	PasswordDigest string
	Nickname       string
	Status         string
	Avatar         string `gorm:"size:1000"`
}

const (
	// PassWordCost 密码加密难度
	PassWordCost = 12
	// Active 激活用户
	Active string = "active"
	// Inactive 未激活用户
	Inactive string = "inactive"
	// Suspend 被封禁用户
	Suspend string = "suspend"
)

// GetUser 用ID获取用户
func GetUserByID(ID interface{}) (User, error) {
	var user User
	result := DB.First(&user, ID)
	return user, result.Error
}

// GetUser 用UserName获取用户
func GetUserByName(username string) (User, error) {
	var user User
	result := DB.Where("user_name = ?", username).First(&user)
	return user, result.Error
}

// SetPassword 设置密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		if err != nil {
			utils.Logger.Errorln(err)
		}
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}

// AvatarURL 获取带签名的用户头像地址
func (user *User) AvatarURL() string {
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

	name := user.Avatar
	ctx := context.Background()

	// 获取预签名URL
	presignedURL, err := client.Object.GetPresignedURL(ctx, http.MethodGet, name, ak, sk, time.Hour, nil)
	if err != nil {
		utils.Logger.Errorln(err)
	}
	return presignedURL.String()
}
