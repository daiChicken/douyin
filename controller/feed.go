package controller

import (
	"BytesDanceProject/service"
	"BytesDanceProject/tool"
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

	//获取参数
	//latest_time 为可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
	latestTime, err := strconv.ParseInt(c.Query("latest_time"), 10, 64)
	if err != nil || latestTime == 0 {
		latestTime = time.Now().UnixNano() / int64(time.Millisecond)
	}

	//获取视频列表及下一次请求的时间戳
	originalVideoList, nextTime, err := service.ListVideos(maxVideoCount, latestTime)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
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
		if exist {
			author.Id = user.Id
			author.Name = user.UserName

			// todo: 完成以下数据的真实获取
			author.FollowCount = 2
			author.FollowerCount = 3
			author.IsFollow = false
		}

		var favoriteCount int64 = 666 //！！！！假数据
		//查询当前视频的点赞数

		commentCount, err := service.CountCommentByVideoId(originalVideo.Id)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}

		isFavorite := false //！！！！！！假数据
		//查询当前登录用户是否喜欢该视频。如果当前用户没有登录，则为false

		playUrl := "http://" + viper.GetString("qiniuyun.domain") + "/" + originalVideo.PlayUrl
		coverUrl := "http://" + viper.GetString("qiniuyun.domain") + "/" + originalVideo.CoverUrl

		video := Video{ //注意video中omitempty！！！
			Id:            int64(originalVideo.Id),          //若为0则生成json时不包含该字段
			Author:        author,                           //待处理
			PlayUrl:       playUrl,                          //若为空则生成json时不包含该字段
			CoverUrl:      coverUrl,                         //若为空则生成json时不包含该字段
			FavoriteCount: favoriteCount,                    //若为0则生成json时不包含该字段
			CommentCount:  commentCount,                     //若为0则生成json时不包含该字段
			IsFavorite:    isFavorite,                       ////若为false则生成json时不包含该字段
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
