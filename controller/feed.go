package controller

import (
	"BytesDanceProject/model"
	"BytesDanceProject/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	//VideoList []Video `json:"video_list,omitempty"`//【正确的】
	VideoList []model.Video `json:"video_list,omitempty"` //【这是错误的，仅供测试使用】
	NextTime  int64         `json:"next_time,omitempty"`
}

const VideoCount = 30

// Feed same demo video list for every request
func Feed(c *gin.Context) {

	//获取参数
	//latest_time为可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
	latestTime := c.Query("latest_time")
	fmt.Println("latestTime", latestTime)

	//业务处理 获取视频列表
	videoList, err := service.GetVideoList(VideoCount)
	if err != nil {
		fmt.Println(err)
		return
	}

	//返回响应
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: *videoList,
		NextTime:  time.Now().Unix(),
	})
}
