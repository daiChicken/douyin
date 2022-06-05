package controller

import (
	"BytesDanceProject/model"
	"BytesDanceProject/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
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
	_ = c.ShouldBindJSON(favoriteRequest)
	tmp := c.Request.Header.Get("Authorization")
	token := strings.SplitN(tmp, " ", 2)[1]

	if _, exist := usersLoginInfo[token]; exist {
		if err := service.FavoriteAction(favoriteRequest); err != nil {
			zap.L().Error("service.FavoriteAction() failed", zap.Error(err))
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败哦！"})
		} else {
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "点赞成功~"})
		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	tmp := c.Request.Header.Get("Authorization")
	token := strings.SplitN(tmp, " ", 2)[1]

	var favoriteListRequest model.FavoriteListRequest
	_ = c.ShouldBindJSON(favoriteListRequest)

	if _, exist := usersLoginInfo[token]; exist {
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
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}
