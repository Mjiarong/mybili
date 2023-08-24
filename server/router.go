package server

import (
	"github.com/gin-gonic/gin"
	"mybili/api"
	"mybili/middleware"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger())
	// 中间件, 注意顺序
	//r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	r.Use(middleware.Cors())
	//r.Use(middleware.CurrentUser())
	r.Use(gin.Recovery())
	// 路由
	v1 := r.Group("/api/v1")
	{
		v1.POST("ping", api.Ping)
		// 用户注册
		v1.POST("user/register", api.UserRegister)
		// 用户登录
		v1.POST("user/login", api.UserLogin)

		// 视频操作
		v1.GET("video/:id", api.ShowVideo)
		v1.GET("videos", api.ListVideo)
		v1.GET("videos/:username", api.ListVideosOfUser)
		//获取用户对于某个视频的播放记录信息
		v1.GET("video/playinfo", api.QueryPlayInfo)

		// 排行榜
		v1.GET("rank/daily", api.DailyRank)

		//获取签名
		v1.POST("upload/token", api.UploadToken)

		//获取临时秘钥相关凭证
		v1.GET("upload/credentials", api.TmpCredentials)

		//获取评论
		v1.GET("comment", api.ListComment)

		// 需要登录保护的
		authed := v1.Group("/")
		authed.Use(middleware.JwtToken())
		{
			// User Routing
			authed.GET("user/:name", api.UserInfo)
			authed.DELETE("user/logout", api.UserLogout)
			authed.POST("videos", api.CreateVideo)
			authed.PUT("video/:id", api.UpdateVideo)
			authed.DELETE("video/:id", api.DeleteVideo)
			authed.POST("comment", api.CreateComment)
			authed.POST("comment/likes/", api.AddCommentLikes)
			authed.DELETE("comment", api.DeleteComments)
			authed.POST("comment/dislikes/", api.AddCommentDislikes)
			authed.POST("video/likes/add", api.AddVideoLikes)
		}
	}
	return r
}
