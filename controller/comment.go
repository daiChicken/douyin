package controller

import (
	"BytesDanceProject/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CommentListResponse 评论列表的返回结构体
type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

// CommentResponse 发表评论的返回结构体
type CommentResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction 发表评论
func CommentAction(c *gin.Context) {

	//用户鉴权【获取不到token参数！！！！！！！！！！！！！！客户端的问题】
	//token := c.PostForm("token")
	//claim, err := jwt.ParseToken(token)
	//if err != nil {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	//	return
	//} else if claim.Valid() != nil {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: claim.Valid().Error()})
	//	return
	//}

	//userId, err := strconv.Atoi(c.Query("user_id")) //评论者id【userid获取不到！！！！！！！！！！！！！！！！客户端的问题】
	//if err != nil {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论发布失败"})
	//	fmt.Println("user_id有问题：" + err.Error())
	//	return
	//}
	userId := 22 //假数据！！！！！！！！！！！！！！！！！

	videoId, err := strconv.Atoi(c.Query("video_id")) //被评论的视频的id
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论发布失败"})
		fmt.Println("video_id有问题：" + err.Error())
		return
	}

	actionType := c.Query("action_type") //1：发布评论；2：删除评论

	if actionType == "1" { //发布评论
		commentText := c.Query("comment_text") //评论内容（type==1时使用）

		originalComment, err := service.CreateComment(userId, videoId, commentText)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论发布失败"})
			fmt.Println("service.CreateComment有问题：" + err.Error())
			return
		}

		//获取用户对象
		originalUser, exist := service.GetUserByID(userId)
		if !exist {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论发布失败"})
			fmt.Println("CommentAction获取用户对象失败")
			return
		}

		user := User{
			Id:   originalUser.Id,
			Name: originalUser.UserName,

			FollowCount:   0,     //假数据【！！！！！！！！！！！！！！！！！！！】
			FollowerCount: 0,     //假数据【！！！！！！！！！！！！！！！！！！！】
			IsFollow:      false, //假数据【！！！！！！！！！！！！！！！！！！！】
		}

		comment := Comment{
			Id:         int64(originalComment.ID),
			User:       user,
			Content:    originalComment.Content,
			CreateDate: originalComment.CreateDate.Format("01-02"),
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
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论删除失败"})
			fmt.Println("comment_id有问题：" + err.Error())
			return
		}
		err = service.DeleteComment(commentId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论删除失败"})
			fmt.Println("service.DeleteComment有问题：" + err.Error())
			return
		}

		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "成功删除评论"})
		return
	}

	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "actionType出错"})
}

// CommentList 评论列表
func CommentList(c *gin.Context) {

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

	videoId, err := strconv.Atoi(c.Query("video_id")) //获取该视频的评论列表
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取评论列表失败"})
		fmt.Println("video_id有问题：" + err.Error())
		return
	}

	originalCommentList, err := service.ListComment(videoId)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取评论列表失败"})
		fmt.Println("service.ListCommentByVideoId有问题：" + err.Error())
		return
	}

	var commentList = make([]Comment, len(*originalCommentList))
	point := 0 //commentList的指针
	for _, originalComment := range *originalCommentList {

		userId := originalComment.UserID //发布这条评论的用户的id

		//获取用户对象
		originalUser, exist := service.GetUserByID(userId)
		if !exist {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取评论列表失败"})
			fmt.Println("service.GetUserByID有问题：" + err.Error())
			return
		}

		user := User{
			Id:   originalUser.Id,
			Name: originalUser.UserName,

			FollowCount:   0,     //假数据【！！！！！！！！！！！！！！！！！！！】
			FollowerCount: 0,     //假数据【！！！！！！！！！！！！！！！！！！！】
			IsFollow:      false, //假数据【！！！！！！！！！！！！！！！！！！！】
		}

		comment := Comment{
			Id:         int64(originalComment.ID),
			User:       user,
			Content:    originalComment.Content,
			CreateDate: originalComment.CreateDate.Format("01-02"),
		}
		commentList[point] = comment
		point++
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "成功获取评论列表",
		},
		CommentList: commentList,
	})
}
