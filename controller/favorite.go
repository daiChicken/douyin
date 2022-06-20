package controller

import (
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/model"
	"BytesDanceProject/service"
	"BytesDanceProject/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
)

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

// FavoriteAction 点赞操作
func FavoriteAction(c *gin.Context) {

	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败"})
		return
	}
	activeUserId := userIdInterface.(int)

	actionType := c.Query("action_type") //1-点赞，2-取消点赞

	videoId, err := strconv.Atoi(c.Query("video_id")) //被点赞的视频的id
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败"})
		fmt.Println("点赞失败" + err.Error())
		return
	}

	likeStatus, err := service.GetLikeStatus(videoId, activeUserId)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败"})
		fmt.Println("点赞失败" + err.Error())
		return
	}

	if actionType == "1" && !likeStatus { //点赞

		err := service.LikeAction(activeUserId, videoId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败"})
			fmt.Println("点赞失败" + err.Error())
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "点赞成功"})
		return

	} else if actionType == "2" || likeStatus { //取消点赞

		err := service.CancelLikeAction(activeUserId, videoId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "取消点赞失败"})
			fmt.Println("取消点赞失败" + err.Error())
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "取消点赞成功"})
		return
	}

	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "业务出错"})
	return

}

// FavoriteList 获取点赞列表
func FavoriteList(c *gin.Context) {
	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败"})
		return
	}
	activeUserId := userIdInterface.(int)

	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取点赞列表失败"})
		fmt.Println("获取点赞列表失败" + err.Error())
		return
	}

	originalVideoList, err := service.ListLikeVideo(int(userId))
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取点赞列表失败"})
		fmt.Println("获取点赞列表失败" + err.Error())
		return
	}

	//获取登录用户的所有关注
	followList, err := service.GetFollowList(&model.FollowListRE{
		UserID: int64(activeUserId),
		Token:  "",
	})
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取点赞列表失败"})
		fmt.Println("获取点赞列表失败" + err.Error())
		return
	}

	//获取到的originalVideoList（model.Video）需要进行处理，使其变成满足前端接口的要求的videoList（controller.Video）
	var videoList = make([]Video, len(*originalVideoList))
	point := 0 //videoList的指针
	for _, originalVideo := range *originalVideoList {

		//根据authorId获取author对象
		//authorId := originalVideo.AuthorId

		user, err := service.GetUserByID(originalVideo.AuthorId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取点赞列表失败"})
			fmt.Println("获取点赞列表失败" + err.Error())
			return
		}

		followCount, followerCount, err := mysql.GetCountByID(int64(user.Id))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取点赞列表失败"})
			fmt.Println("获取点赞列表失败" + err.Error())
			return
		}

		isFollow := false
		for _, val := range followList {
			if val.UserName == user.UserName { //视频作者存在于当前登录用户的关注列表中
				isFollow = true
			}
		}

		author := User{
			Id:            user.Id,
			Name:          user.UserName,
			FollowCount:   followCount,
			FollowerCount: followerCount,
			IsFollow:      isFollow,
		}

		likeCount, err := service.CountLike(originalVideo.Id)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取点赞列表失败"})
			fmt.Println("获取点赞列表失败" + err.Error())
			return
		}

		commentCount, err := service.CountCommentByVideoId(originalVideo.Id)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取点赞列表失败"})
			fmt.Println("获取点赞列表失败" + err.Error())
			return
		}

		likeStatus, err := service.GetLikeStatus(originalVideo.Id, activeUserId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取点赞列表失败"})
			fmt.Println(err.Error())
			return
		}

		playUrl := "http://" + viper.GetString("qiniuyun.domain") + "/" + originalVideo.PlayUrl
		coverUrl := "http://" + viper.GetString("qiniuyun.domain") + "/" + originalVideo.CoverUrl

		video := Video{ //注意video中omitempty！！！
			Id:            int64(originalVideo.Id),          //若为0则生成json时不包含该字段
			Author:        author,                           //待处理
			PlayUrl:       playUrl,                          //若为空则生成json时不包含该字段
			CoverUrl:      coverUrl,                         //若为空则生成json时不包含该字段
			FavoriteCount: likeCount,                        //若为0则生成json时不包含该字段
			CommentCount:  commentCount,                     //若为0则生成json时不包含该字段
			IsFavorite:    likeStatus,                       ////若为false则生成json时不包含该字段
			Title:         tool.Filter(originalVideo.Title), //若为空则生成json时不包含该字段
		}
		videoList[point] = video
		point++
	}

	//返回响应
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "成功获取点赞列表",
		},
		VideoList: videoList,
	})

}
