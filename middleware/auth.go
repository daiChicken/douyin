package middleware

import (
	"BytesDanceProject/controller"
	"BytesDanceProject/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWTAuth 基于JWT的认证中间件
func JWTAuth() func(c *gin.Context) { //需要登录才能访问的地方加这个中间件即可，因为这个中间件就是校验Token的
	return func(c *gin.Context) {
		token := c.Query("token")

		mc, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, controller.Response{StatusCode: 1, StatusMsg: "请先登录！"})
			c.Abort()
			return
		}

		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set("userId", mc.UserId)
		c.Set("username", mc.Username)
		c.Next() // 后续的处理函数可以用过c.Get("userID")来获取当前请求的用户信息
	}
}
