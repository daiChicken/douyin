package controller

import (
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/model"
	"BytesDanceProject/pkg/jwt"
	"BytesDanceProject/service"
	"BytesDanceProject/tool"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

const maxVideoCount = 30 //一次请求最多返回的视频数

// Feed 拉取feed流
func Feed(c *gin.Context) {

	token := c.Query("token")

	claim, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "拉取feed流失败"})
		fmt.Println("拉取feed流失败" + err.Error())
	}

	//获取参数
	//latest_time 为可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
	latestTime, err := strconv.ParseInt(c.Query("latest_time"), 10, 64)
	if err != nil || latestTime == 0 {
		latestTime = time.Now().UnixNano() / int64(time.Millisecond)
	}

	//获取视频列表及下一次请求的时间戳
	originalVideoList, nextTime, err := service.ListVideos(maxVideoCount, latestTime)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "拉取feed流失败"})
		fmt.Println("拉取feed流失败" + err.Error())
		return
	}

	followList, err := service.GetFollowList(&model.FollowListRE{
		UserID: int64(claim.UserID),
		Token:  "",
	})
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "拉取feed流失败"})
		fmt.Println("拉取feed流失败" + err.Error())
		return
	}

	//获取到的originalVideoList（model.Video）需要进行处理，使其变成满足前端接口的要求的videoList（controller.Video）
	var videoList = make([]Video, len(*originalVideoList))
	point := 0 //videoList的指针
	for _, originalVideo := range *originalVideoList {

		//根据authorId获取author对象
		//authorId := originalVideo.AuthorId
		author := User{}
		user, exist := service.GetUserByID(originalVideo.AuthorId)
		followCount, followerCount, err := mysql.GetCountByID(int64(user.Id))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "拉取feed流失败"})
			fmt.Println("拉取feed流失败" + err.Error())
			return
		}

		isFollow := false
		for _, val := range followList {
			if val.UserName == user.UserName {
				isFollow = true
			}
		}

		if exist {
			author.Id = user.Id
			author.Name = user.UserName
			author.FollowCount = followCount
			author.FollowerCount = followerCount
			author.IsFollow = isFollow
		}

		likeCount, err := service.CountLike(originalVideo.Id)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "拉取feed流失败"})
			fmt.Println("拉取feed流失败" + err.Error())
			return
		}

		commentCount, err := service.CountCommentByVideoId(originalVideo.Id)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "拉取feed流失败"})
			fmt.Println("拉取feed流失败" + err.Error())
			return
		}

		likeStatus, err := service.GetLikeStatus(originalVideo.Id, claim.UserID)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取发布列表失败"})
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
	c.JSON(http.StatusOK, FeedResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "成功获取视频列表",
		},
		VideoList: videoList,
		NextTime:  nextTime,
	})
}
