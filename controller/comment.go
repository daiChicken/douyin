package controller

import (
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/pkg/jwt"
	"BytesDanceProject/service"
	"BytesDanceProject/tool"
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

	//用户鉴权
	token := c.Query("token")

	claim, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论失败"})
		fmt.Println("评论失败", err.Error())
		return
	} else if claim.Valid() != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论失败"})
		fmt.Println("评论失败", claim.Valid().Error())
		return
	}

	userId := claim.UserID

	videoId, err := strconv.Atoi(c.Query("video_id")) //被评论的视频的id
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论发布失败"})
		fmt.Println("评论发布失败" + err.Error())
		return
	}

	actionType := c.Query("action_type") //1：发布评论；2：删除评论

	if actionType == "1" { //发布评论
		commentText := c.Query("comment_text") //评论内容（type==1时使用）

		originalComment, err := service.CreateComment(userId, videoId, commentText, claim.Username)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论发布失败"})
			fmt.Println("评论发布失败" + err.Error())
			return
		}

		followCount, followerCount, err := mysql.GetCountByID(int64(userId))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论发布失败"})
			fmt.Println("评论发布失败" + err.Error())
			return
		}

		user := User{
			Id:   int64(userId),
			Name: claim.Username,

			FollowCount:   followCount,   //关注总数
			FollowerCount: followerCount, //粉丝总数
			IsFollow:      false,         //关注关系【！！！！！！！！！！！！！！！！！！！】
		}

		comment := Comment{
			Id:         int64(originalComment.ID),
			User:       user,
			Content:    tool.Filter(originalComment.Content), //使用过滤器过滤评论内容
			CreateDate: originalComment.CreateDate.Format("01-02"),
		}

		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{StatusCode: 0, StatusMsg: "成功发布评论"},
			Comment:  comment,
		})
		return

	} else if actionType == "2" { //删除评论
		commentId, err := strconv.Atoi(c.Query("comment_id")) //要删除的评论的id（type==2时使用）
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论删除失败"})
			fmt.Println("评论删除失败" + err.Error())
			return
		}
		err = service.DeleteComment(commentId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论删除失败"})
			fmt.Println("评论删除失败" + err.Error())
			return
		}

		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "成功删除评论"})
		return
	}

	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "actionType出错"})
}

// CommentList 评论列表
func CommentList(c *gin.Context) {

	//用户鉴权
	claim, err := jwt.ParseToken(c.Query("token"))
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取评论列表失败"})
		fmt.Println("获取评论列表失败", err.Error())
		return
	} else if claim.Valid() != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取评论列表失败"})
		fmt.Println("获取评论列表失败", claim.Valid().Error())
		return
	}

	videoId, err := strconv.Atoi(c.Query("video_id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取评论列表失败"})
		fmt.Println("获取评论列表失败" + err.Error())
		return
	}

	originalCommentList, err := service.ListComment(videoId)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取评论列表失败"})
		fmt.Println("获取评论列表失败" + err.Error())
		return
	}

	var commentList = make([]Comment, len(*originalCommentList))
	point := 0 //commentList的指针
	for _, originalComment := range *originalCommentList {

		followCount, followerCount, err := mysql.GetCountByID(int64(originalComment.UserID))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取评论列表失败"})
			fmt.Println("评论列表获取失败" + err.Error())
			return
		}

		user := User{
			Id:   int64(originalComment.UserID),
			Name: originalComment.UserName,

			FollowCount:   followCount,   //关注总数
			FollowerCount: followerCount, //粉丝总数
			IsFollow:      false,         //假数据【！！！！！！！！！！！！！！！！！！！】
		}

		comment := Comment{
			Id:         int64(originalComment.ID),
			User:       user,
			Content:    tool.Filter(originalComment.Content), //使用过滤器过滤评论内容
			CreateDate: originalComment.CreateDate.Format("01-02"),
		}
		commentList[point] = comment
		point++
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "成功获取评论列表"},
		CommentList: commentList,
	})
}
