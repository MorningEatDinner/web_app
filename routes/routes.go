package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiaorui/web_app/controller"
	"github.com/xiaorui/web_app/logger"
	"github.com/xiaorui/web_app/middlewares"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 设置为发布模式
	}
	r := gin.New()                                                                                          // 我尝试进行新的变化
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(time.Second*2, 10)) // 2s新增1个令牌， 容量为10

	v1 := r.Group("/api/v1")
	{
		// 创建验证相关的路由组
		authGroup := v1.Group("/auth")
		// 注册业务路由
		{
			authGroup.POST("/signup", controller.SignUpHandler)
			authGroup.POST("/login", controller.LoginHandler)
		}

		// 后面的所有请求都需要使用这个中间件，即需要验证是否进行了登陆
		v1.Use(middlewares.JWTAuthMiddleware())

		// 创建用户相关的路由组
		usersGroup := v1.Group("/user")
		{
			usersGroup.GET("", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{
					"msg": "没问题",
				})
			})
		}

		commGroup := v1.Group("/community")
		{
			commGroup.GET("/", controller.CommunityHandler)
			commGroup.GET("/:id", controller.CommunityDetailHandler)
		}

		postGroup := v1.Group("/post")
		{

			postGroup.POST("/", controller.CreatePostHandler)
			postGroup.GET("/:id", controller.GetPostHandler)
			postGroup.GET("/posts", controller.GetPostListHandler)
			postGroup.GET("/posts2", controller.GetPostListHandler0)
			postGroup.POST("/vote", controller.PostVoteHandler)
			postGroup.GET("/posts3", controller.GetPostListHandler0)
		}
	}

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
