package routes

import (
	"BytesDanceProject/controller"
	"BytesDanceProject/logger"
	"BytesDanceProject/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	//如果设置mode为release则设置gin为该模式
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Static("/static", "./public")
	//r.Use(logger.GinLogger(),logger.GinRecovery(true),middleware.RateLimitMiddleware(time.Second,1))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//设置一个路由组
	apiRouter := r.Group("/douyin")
	// apiRouter .Use(middleware.JWTAuth())
	{
		// basic apis
		apiRouter.GET("/feed/", controller.Feed)
		apiRouter.GET("/user/", middleware.JWTAuth(), controller.UserInfo)
		apiRouter.POST("/user/register/", controller.Register)
		apiRouter.POST("/user/login/", controller.Login)
		apiRouter.POST("/publish/action/", controller.Publish)
		apiRouter.GET("/publish/list/", middleware.JWTAuth(), controller.PublishList)

		// extra apis - I
		apiRouter.POST("/favorite/action/", middleware.JWTAuth(), controller.FavoriteAction)
		apiRouter.GET("/favorite/list/", middleware.JWTAuth(), controller.FavoriteList)
		apiRouter.POST("/comment/action/", middleware.JWTAuth(), controller.CommentAction)
		apiRouter.GET("/comment/list/", middleware.JWTAuth(), controller.CommentList)

		// extra apis - II
		apiRouter.POST("/relation/action/", middleware.JWTAuth(), controller.RelationAction)
		apiRouter.GET("/relation/follow/list/", middleware.JWTAuth(), controller.FollowList)
		apiRouter.GET("/relation/follower/list/", middleware.JWTAuth(), controller.FollowerList)
	}

	return r
}
