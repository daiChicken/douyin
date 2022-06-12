package controller

import (
	"BytesDanceProject/model"
	"BytesDanceProject/pkg/jwt"
	"BytesDanceProject/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

//type FavoriteActionResponse struct {
//}

type FavoriteListResponse struct {
	Response
	VideoList []VideoList `json:"video_list"`
}

type VideoList struct {
	VideoID       int64  `json:"video_id"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int    `json:"favorite_count"`
	CommentCount  int    `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}

// 点赞
// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	//var favoriteRequest model.FavoriteRequest
	//_ = c.ShouldBindJSON(favoriteRequest)
	//tmp := c.Request.Header.Get("Authorization")
	token := c.Query("token")

	claim, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	} else if claim.Valid() != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: claim.Valid().Error()})
		return
	}

	userId := claim.UserID               //当前登录用户id，点赞用户id
	actionType := c.Query("action_type") //1-点赞，2-取消点赞

	videoId, err := strconv.Atoi(c.Query("video_id")) //被点赞的视频的id
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败"})
		fmt.Println("点赞失败" + err.Error())
		return
	}

	//favoriteRequest.ActionType = cast.ToInt32(actionType)
	//favoriteRequest.UserID = int64(userId)
	//favoriteRequest.Token = token
	//spew.Dump("=======================controller:favoriteRequest", favoriteRequest)
	//
	//if err := service.FavoriteAction(favoriteRequest); err != nil {
	//	zap.L().Error("service.FavoriteAction() failed", zap.Error(err))
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败哦！"})
	//} else {
	//	ResponseSuccess(c, CodeSuccess, nil)
	//}

	likeStatus, err := service.GetLikeStatus(videoId, userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败"})
		fmt.Println("点赞失败" + err.Error())
		return
	}
	fmt.Println("!!!!!", likeStatus)
	fmt.Println("!!!!!", actionType == "1" && !likeStatus)
	fmt.Println("!!!!!", actionType == "2" || likeStatus)

	if actionType == "1" && !likeStatus { //点赞

		err := service.LikeAction(userId, videoId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败"})
			fmt.Println("点赞失败" + err.Error())
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞成功"})
		return

	} else if actionType == "2" || likeStatus { //取消点赞

		err := service.CancelLikeAction(userId, videoId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "取消点赞失败"})
			fmt.Println("取消点赞失败" + err.Error())
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "取消点赞成功"})
		return
	}

	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "业务出错"})
	return

}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	var favoriteListRequest model.FavoriteListRequest

	token := c.Query("token")
	claim, err := jwt.ParseToken(token)
	if err != nil {
		zap.L().Error("service.FavoriteAction() failed", zap.Error(err))
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败哦！"})
	}
	userId := claim.UserID
	favoriteListRequest.Token = token
	favoriteListRequest.UserID = int64(userId)

	if err := service.FavoriteList(favoriteListRequest); err != nil {
		zap.L().Error("service.FavoriteList() failed", zap.Error(err))
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取点赞列表失败哦！"})
	} else {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
			},
			VideoList: DemoVideos,
		})
	}

}
