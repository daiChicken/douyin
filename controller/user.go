package controller

import (
	"BytesDanceProject/model"
	"BytesDanceProject/pkg/jwt"
	"BytesDanceProject/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

// UserInfo 用户信息 获取用户的 id、昵称，如果实现社交部分的功能，还会返回关注数和粉丝数
func UserInfo(c *gin.Context) {
	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败"})
		return
	}
	activeUserId := userIdInterface.(int)

	//获取登录用户的所有关注
	followList, err := service.GetFollowList(&model.FollowListRE{
		UserID: int64(activeUserId),
		Token:  "",
	})
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取发用户信息失败"})
		fmt.Println(err.Error())
		return
	}

	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取用户信息失败"})
		fmt.Println(err.Error())
		return
	}

	follower, follow := service.GetUserInfo(userId)

	user, err := service.GetUserByID(int(userId))
	isFollow := false
	for _, val := range followList {
		if val.UserName == user.UserName { //获取信息的用户存在于当前登录用户的关注列表中
			isFollow = true
		}
	}

	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0, StatusMsg: "成功获取用户信息！"},
		User: User{
			Id:            userId,
			Name:          user.UserName,
			FollowCount:   follow,
			FollowerCount: follower,
			IsFollow:      isFollow,
		},
	})
}

// Login 登录
func Login(c *gin.Context) {

	username := c.Query("username")
	password := c.Query("password")

	userId, err := service.Login(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "登录失败！"},
		})
		fmt.Println("登录失败！" + err.Error())
	}

	token, err := jwt.GenToken(int(userId), username)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "登录失败！"},
		})
		fmt.Println("登录失败！" + err.Error())
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0, StatusMsg: "登录成功！"},
		UserId:   userId,
		Token:    token,
	})

}

// Register 注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	userId, err := service.Register(username, password) //注册

	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "注册失败！"},
		})
		fmt.Println("注册失败！" + err.Error())
	}

	token, err := jwt.GenToken(int(userId), username)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "注册失败！"},
		})
		fmt.Println("注册失败！" + err.Error())
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0, StatusMsg: "注册成功！"},
		UserId:   userId,
		Token:    token,
	})

}
