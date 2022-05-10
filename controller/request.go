package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var ErrorUserNotExist = errors.New("用户未登录")
const CtxUserIDKey = "Openid"

// getCurrentUser 获取当前用户唯一标识（根据上下文GET出来那个key）
func getCurrentUser(c *gin.Context)(openid string,err error){
	uid ,ok  := c.Get(CtxUserIDKey)
	if !ok{
		err = ErrorUserNotExist
		return
	}
	//类型断言，openid是一个空接口
	openid , ok = uid.(string)
	if !ok{
		err = ErrorUserNotExist
		return
	}
	return
}
