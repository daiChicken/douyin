package controller

import (
	"BytesDanceProject/pkg/jwt"
	"BytesDanceProject/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
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

func UserInfo(c *gin.Context) {
	tmp := c.Request.Header.Get("Authorization")
	token := strings.SplitN(tmp, " ", 2)[1]
	if user, exist := usersLoginInfo[token]; exist {
		//Id, err := service.FindUser(user.Name)
		//fmt.Println(err)
		//user.Id = Id.Id
		//fmt.Println(user.Name)
		//fmt.Println("!!!!!!")
		//followerTable, _ := mysql.GetFollower(user.Id)
		//user.FollowerCount = int64(len(followerTable))
		//followTable, _ := mysql.GetFollowed(user.Id)
		//user.FollowCount = int64(len(followTable))
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0, StatusMsg: "UserInfo get"},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist (userinfo)"},
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	Flag := service.VerifyLogin(username, password)
	user, _ := service.FindUser(username)
	newUser := User{
		Id:   user.Id,
		Name: username,
	}
	token, _ := jwt.GenToken(username)
	usersLoginInfo[token] = newUser
	if Flag {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Login success"},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist or Password is wrong"},
		})
	}
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //加密处理
	if err != nil {
		fmt.Println(err)
	}
	encodePWD := string(hash)
	Flag := service.Register(username, encodePWD)
	if !Flag {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		user, _ := service.FindUser(username)
		newUser := User{
			Id:   user.Id,
			Name: username,
		}
		token, _ := jwt.GenToken(username)
		usersLoginInfo[token] = newUser
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	}
}
