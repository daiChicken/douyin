package controller

import (
	"BytesDanceProject/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

const maxVideoCount = 30 //一次请求最多返回的视频数

// Feed same demo video list for every request
func Feed(c *gin.Context) {

	//获取参数
	//latest_time为可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
	latestTime := c.Query("latest_time") //【未处理！！！！！！！！！重要】
	fmt.Println("latestTime", latestTime)

	//获取视频列表
	originalVideoList, err := service.ListVideos(maxVideoCount)
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

		var favoriteCount int64 = 666 //！！！！假数据
		//查询当前视频的点赞数

		var commentCount int64 = 777 //！！！！假数据
		//查询当前视频的评论数

		isFavorite := false //！！！！！！假数据
		//查询当前登录用户是否喜欢该视频。如果当前用户没有登录，则为false

		video := Video{ //注意video中omitempty！！！
			Id:            int64(originalVideo.Id),           //若为0则生成json时不包含该字段
			Author:        User{},                            //待处理
			PlayUrl:       "http://" + originalVideo.PlayUrl, //若为空则生成json时不包含该字段
			CoverUrl:      originalVideo.CoverUrl,            //若为空则生成json时不包含该字段
			FavoriteCount: favoriteCount,                     //若为0则生成json时不包含该字段
			CommentCount:  commentCount,                      //若为0则生成json时不包含该字段
			IsFavorite:    isFavorite,                        ////若为false则生成json时不包含该字段
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
		NextTime:  time.Now().Unix(),
	})
}
