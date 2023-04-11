package route

import (
	"chat/api"
	"chat/middleware"
	"chat/service/chat"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/")
	{
		//用户注册
		v1.POST("user/register", api.UserRegister)
		//用户登录
		v1.POST("user/login", api.UserLogin)

		friends := v1.Group("friends/")
		friends.Use(middleware.JWT())
		{
			friends.GET("get", api.GetFriends)
			friends.POST("/", api.MakeFriends)
			friends.GET("application", api.GetApplication)
			friends.PUT("application", api.Reply)
			friends.GET("message", api.SingleMessage)
		}

		groups := v1.Group("groups/")
		groups.Use(middleware.JWT())
		{
			groups.GET("get", api.GetGroups)
			groups.POST("/", api.MakeGroups)
			groups.POST("join", api.JoinGroups)
			groups.GET("message", api.GroupMessage)
		}
		ws := v1.Group("ws/")
		ws.Use(middleware.WsJWT())
		ws.GET("single", chat.SingleChat)
		ws.GET("group", chat.GroupChat)

	}
	return r
}
