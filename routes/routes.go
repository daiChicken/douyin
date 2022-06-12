package routes

import (
	"BytesDanceProject/controller"
	"BytesDanceProject/logger"
	"BytesDanceProject/middlewares"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	//如果设置mode为release则设置gin为该模式
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Static("/static", "./public")
	//r.Use(logger.GinLogger(),logger.GinRecovery(true),middlewares.RateLimitMiddleware(time.Second,1))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//设置一个路由组
	apiRouter := r.Group("/douyin")
	// apiRouter .Use(middlewares.JWTAuthMiddleware())
	{
		// basic apis
		apiRouter.GET("/feed/", controller.Feed)
		apiRouter.GET("/user/", middlewares.JWTAuthMiddleware(), controller.UserInfo)
		apiRouter.POST("/user/register/", controller.Register)
		apiRouter.POST("/user/login/", controller.Login)
		apiRouter.POST("/publish/action/", controller.Publish)
		apiRouter.GET("/publish/list/", controller.PublishList)

		// extra apis - I
		apiRouter.POST("/favorite/action/", middlewares.JWTAuthMiddleware(), controller.FavoriteAction)
		apiRouter.GET("/favorite/list/", middlewares.JWTAuthMiddleware(), controller.FavoriteList)
		apiRouter.POST("/comment/action/", controller.CommentAction)
		apiRouter.GET("/comment/list/", controller.CommentList)

		// extra apis - II
		apiRouter.POST("/relation/action/", middlewares.JWTAuthMiddleware(), controller.RelationAction)
		apiRouter.GET("/relation/follow/list/", middlewares.JWTAuthMiddleware(), controller.FollowList)
		apiRouter.GET("/relation/follower/list/", middlewares.JWTAuthMiddleware(), controller.FollowerList)
	}

	return r
}
