package middleware

import (
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"mybili/serializer"
	"mybili/utils"
	"net/http"
	"os"
	"time"
)

type MyClaims struct {
	Username string `json:"user_name"`
	//Password string `json:"password"`
	jwt.StandardClaims
}

// 生成token
func SetToken(username string) (string, int) {
	expireTime := time.Now().Add(10 * time.Hour)
	SetClaims := MyClaims{
		Username: username,
		//Password: password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: &jwt.Time{expireTime},   //有效时间
			Issuer:    os.Getenv("JWT_ISSUER"), //签发人
		},
	}

	//使用指定的加密方式和声明类型创建新令牌
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaims)
	//获得完整的、签名的令牌
	token, err := tokenStruct.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		utils.Logger.Errorf("err:%v", err.Error())
		return "", utils.TOKEN_CREATE_FAILED
	}
	return token, utils.SUCCESS
}

// 验证token
func CheckToken(token string) (*MyClaims, int) {
	//解析、验证并返回token。
	// func将接收解析后的token，并返回key进行验证。
	//如果一切正常，err将为nil
	tokenObj, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		utils.Logger.Errorf("err:%v", err.Error())
		return nil, utils.ERROR
	}

	if key, _ := tokenObj.Claims.(*MyClaims); tokenObj.Valid {
		return key, utils.SUCCESS
	} else {
		return nil, utils.ERROR
	}

}

// jwt中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get("Authorization")
		code := utils.SUCCESS
		if tokenHeader == "" {
			code = utils.TOKEN_NOT_EXIST
			c.JSON(http.StatusOK, serializer.CheckToken(
				code,
				utils.GetErrMsg(code)))
			c.Abort()
			return
		}

		key, tCode := CheckToken(tokenHeader)
		if tCode == utils.ERROR {
			code = utils.TOKEN_WRONG
			c.JSON(http.StatusOK, serializer.CheckToken(
				code,
				utils.GetErrMsg(code)))
			c.Abort()
			return
		}

		if time.Now().Unix() > key.ExpiresAt.Unix() {
			code = utils.TOKEN_RUNTIME
			c.JSON(http.StatusOK, serializer.CheckToken(
				code,
				utils.GetErrMsg(code)))
			c.Abort()
			return
		}

		c.Set("username", key.Username)
		c.Next()
	}
}
