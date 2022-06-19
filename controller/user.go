package controller

import (
	"BytesDanceProject/pkg/jwt"
	"BytesDanceProject/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

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
	token := c.Query("token")
	if user, exist := usersLoginInfo[token]; exist {
		follower, follow := service.GetUserInfo(user.Id)
		user.FollowerCount = follower
		user.FollowCount = follow
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist (userinfo)"},
		})
	}
}

// Login 登录
func Login(c *gin.Context) {

	username := c.Query("username")
	password := c.Query("password")

	id, flag := service.VerifyLogin(username, password) //验证用户名和密码

	newUser := User{
		Id:   id,
		Name: username,
	}

	token, _ := jwt.GenToken(int(id), username)

	usersLoginInfo[token] = newUser

	if flag {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Login success"},
			UserId:   id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist or Password is wrong"},
		})
	}
}

// Register 注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //加密处理
	if err != nil {
		fmt.Println(err)
	}

	encodePWD := string(hash)
	id, flag := service.Register(username, encodePWD) //注册

	if !flag {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		newUser := User{
			Id:   id,
			Name: username,
		}
		token, _ := jwt.GenToken(int(id), username)
		usersLoginInfo[token] = newUser
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   id,
			Token:    token,
		})
	}
}
