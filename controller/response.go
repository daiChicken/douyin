package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
{
	"code": //程序中的错误码
	"msg": //提示信息
	“data" : //数据
}

*/


// ResponseError 返回一个code 类型的错误
func ResponseError(c *gin.Context ,code int32){
	c.JSON(http.StatusOK,&Response{
		StatusCode: code,
		StatusMsg:  Msg(code),
	})
}



// ResponseSuccess 返回成功的信息
func ResponseSuccess(c *gin.Context,code int32){
	c.JSON(http.StatusOK,&Response{
		StatusCode: code,
		StatusMsg:  Msg(code),
	})
}

// ResponseSuccessWithData 返回一个自定义结构体的数据
func ResponseSuccessWithData(c *gin.Context,data interface{}){
	c.JSON(http.StatusOK,data)
}

