package controller

import (
	"BytesDanceProject/model"
	"BytesDanceProject/service"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {

	token := c.Query("token")
	if _, exist := usersLoginInfo[token]; exist {
		// 1.获取参数，参数校验
		p := new(model.Comment)
		if err := c.ShouldBind(p); err != nil {
			fmt.Println(err)
			ResponseError(c, CodeInvalidParam)
			return
		}
		// 创建评论
		if err := service.CreateComment(p); err != nil {
			zap.L().Error("service.CreatePost(p) failed", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}
		// 返回响应
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	commentservice := new(Comment)
	if err := c.ShouldBind(&commentservice); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "请求参数有误",
		})
		return
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: DemoComments,
	})
}
