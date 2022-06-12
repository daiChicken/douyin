package controller

import (
	"BytesDanceProject/model"
	"BytesDanceProject/pkg/jwt"
	"BytesDanceProject/service"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/http"
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
	var favoriteRequest model.FavoriteRequest
	//_ = c.ShouldBindJSON(favoriteRequest)
	//tmp := c.Request.Header.Get("Authorization")
	token := c.Query("token")

	claim, err := jwt.ParseToken(token)

	if err != nil {
		zap.L().Error("service.FavoriteAction() failed", zap.Error(err))
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败哦！"})
	}
	userId := claim.UserID
	actionType := c.Query("action_type")
	favoriteRequest.ActionType = cast.ToInt32(actionType)
	favoriteRequest.UserID = int64(userId)
	favoriteRequest.Token = token
	spew.Dump("=======================controller:favoriteRequest", favoriteRequest)

	if err := service.FavoriteAction(favoriteRequest); err != nil {
		zap.L().Error("service.FavoriteAction() failed", zap.Error(err))
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败哦！"})
	} else {
		ResponseSuccess(c, CodeSuccess, nil)
	}

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
