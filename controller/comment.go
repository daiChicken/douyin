package controller

import (
	"BytesDanceProject/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {

	//用户鉴权【获取不到token参数！！！！！！！！！客户端的问题】
	//token := c.PostForm("token")
	//fmt.Println("CommentAction-token:" + token)
	//
	//claim, err := jwt.ParseToken(token)
	//if err != nil {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	//	return
	//} else if claim.Valid() != nil {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: claim.Valid().Error()})
	//	return
	//}

	//获取参数，参数校验
	userId := 22 //【userid获取不到！！！！！！！！！！！！！！！！客户端的问题】
	//userId, err := strconv.Atoi(c.Query("user_id")) //评论者id
	//if err != nil {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论发布失败"})
	//	fmt.Println("user_id" + err.Error())
	//	return
	//}

	videoId, err := strconv.Atoi(c.Query("video_id")) //被评论的视频的id
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论发布失败"})
		fmt.Println("video_id" + err.Error())
		return
	}

	actionType := c.Query("action_type") //1：发布评论；2：删除评论

	if actionType == "1" { //发布评论
		commentText := c.Query("comment_text") //评论内容（type==1时使用）

		originalComment, err := service.CreateComment(userId, videoId, commentText)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论发布失败"})
			fmt.Println("service.CreateComment" + err.Error())
			return
		}

		comment := Comment{
			Id:         int64(originalComment.ID),
			User:       User{}, //获取该评论的用户【！！！！！！！！！！！！！！！！未完成】
			Content:    originalComment.Content,
			CreateDate: originalComment.CreateDate,
		}

		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "成功发布评论",
			},
			Comment: comment,
		})
		return

	} else if actionType == "2" { //删除评论
		commentId, err := strconv.Atoi(c.Query("comment_id")) //要删除的评论的id（type==2时使用）
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论发布失败"})
			fmt.Println("comment_id" + err.Error())
			return
		}
		err = service.DeleteComment(commentId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论发布失败"})
			fmt.Println("service.DeleteComment" + err.Error())
			return
		}

		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "成功删除评论"})
		return
	}

	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "actionType出错"})

	//p := new(model.Comment)
	//if err := c.ShouldBind(p); err != nil {
	//	fmt.Println(err)
	//	ResponseError(c, CodeInvalidParam)
	//	return
	//}
	//// 创建评论
	//if err := service.CreateComment(p); err != nil {
	//	zap.L().Error("service.CreatePost(p) failed", zap.Error(err))
	//	ResponseError(c, CodeServerBusy)
	//	return
	//}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	//commentservice := new(Comment)
	//if err := c.ShouldBind(&commentservice); err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "请求参数有误",
	//	})
	//	return
	//}

	//用户鉴权【获取不到token参数！！！！！！！！！客户端的问题】
	//token := c.PostForm("token")
	//fmt.Println("CommentList-token:" + token)
	//
	//claim, err := jwt.ParseToken(token)
	//if err != nil {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	//	return
	//} else if claim.Valid() != nil {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: claim.Valid().Error()})
	//	return
	//}

	videoId, err := strconv.Atoi(c.Query("video_id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	originalCommentList, err := service.ListCommentByVideoId(videoId)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	var commentList = make([]Comment, len(*originalCommentList))
	point := 0 //videoList的指针
	for _, originalComment := range *originalCommentList {

		comment := Comment{
			Id:         int64(originalComment.ID),
			User:       User{}, //获取该评论的用户【！！！！！！！！！！！！！！！！未完成】
			Content:    originalComment.Content,
			CreateDate: originalComment.CreateDate,
		}
		commentList[point] = comment
		point++
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentList,
	})
}
