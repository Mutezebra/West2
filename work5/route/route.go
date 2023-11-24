package route

import (
	"five/api"
	"five/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter() *gin.Engine {
	// 初始化一个实例
	r := gin.Default()
	r.GET("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	// 加载静态资源
	r.Static("/static", "./static")

	v1 := r.Group("api/")
	{
		// 用户模块
		v1.POST("user/register", api.UserRegisterHandle())
		v1.POST("user/login", api.UserLoginHandle())
	}

	auth := v1.Group("/")
	auth.Use(middleware.JWT()) // 开启jwt鉴权
	{
		auth.GET("user/websocket", api.WsHandler())
		auth.POST("user/create-group", api.CreateGroupHandle())
		// 用户模块
	}
	return r
}
