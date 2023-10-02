package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"mybili/serializer"
	"mybili/utils"
	"net/http"
	"os"
	"time"
)

type MyClaims struct {
	Username string `json:"user_name"`
	//Password string `json:"password"`
	jwt.RegisteredClaims
}

// 生成token
func SetToken(username string) (string, int) {
	SetClaims := MyClaims{
		Username: username,
		//Password: password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), //有效时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     //签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                     //生效时间
			Issuer:    os.Getenv("JWT_ISSUER"),                            //签发人
			Subject:   "somebody",                                         //主题
			ID:        "1",                                                //JWT ID用于标识该JWT
			Audience:  []string{"somebody_else"},                          //用户
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

	//ParseWithClaims是NewParser().ParseWithClaims()的快捷方式
	//第一个值是token ，
	//第二个值是我们之后需要把解析的数据放入的地方，
	//第三个值是Keyfunc将被Parse方法用作回调函数，以提供用于验证的键。函数接收已解析但未验证的令牌。
	tokenObj, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		utils.Logger.Errorf("err:%v", err.Error())
		return nil, utils.ERROR
	}

	if claims, ok := tokenObj.Claims.(*MyClaims); ok && tokenObj.Valid {
		fmt.Printf("Username:%v\n RegisteredClaims:%v\n", claims.Username, claims.RegisteredClaims)
		return claims, utils.SUCCESS
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

		//判断token是否过期
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
