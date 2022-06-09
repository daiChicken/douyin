package controller

import (
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/pkg/jwt"
	"BytesDanceProject/service"
	"BytesDanceProject/tool"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish 视频发布
func Publish(c *gin.Context) {
	//用户鉴权
	token := c.PostForm("token")

	claim, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	} else if claim.Valid() != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: claim.Valid().Error()})
		return
	}

	//获取标题
	title := c.PostForm("title")

	//获取文件
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//上传文件到七牛云空间
	err = service.UploadVideo(file, title, claim.UserID)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		//StatusMsg:  finalName + " uploaded successfully",
		StatusMsg: "uploaded successfully",
	})
}

// PublishList 获取发布列表
func PublishList(c *gin.Context) {
	//用户鉴权
	token := c.Query("token")

	claim, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取发布列表失败"})
		fmt.Println("获取发布列表失败", err.Error())
		return
	} else if claim.Valid() != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取发布列表失败"})
		fmt.Println("获取发布列表失败", claim.Valid().Error())
		return
	}

	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)

	//获取用户发布的所有视频
	originalVideoList, err := service.ListVideosByUser(int(userId))
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	followCount, followerCount, err := mysql.GetCountByID(int64(userId))
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论发布失败"})
		fmt.Println("评论发布失败" + err.Error())
		return
	}

	//保存视频作者信息
	author := User{
		Id:            userId,
		Name:          claim.Username,
		FollowCount:   followCount,
		FollowerCount: followerCount,

		// todo: 完成以下数据的真实获取
		IsFollow: false,
	}

	//获取到的originalVideoList（model.Video）需要进行处理，使其变成满足前端接口的要求的videoList（controller.Video）
	var videoList = make([]Video, len(*originalVideoList))
	point := 0 //videoList的指针
	for _, originalVideo := range *originalVideoList {

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
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "成功获取当前登录用户所有投稿过的视频",
		},
		VideoList: videoList,
	})
}
