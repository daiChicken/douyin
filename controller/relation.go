package controller

import (
	"BytesDanceProject/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction 关注操作
func RelationAction(c *gin.Context) {
	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "关注失败"})
		return
	}
	activeUserId := userIdInterface.(int)

	toUserId, err := strconv.Atoi(c.Query("to_user_id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "关注失败"})
		fmt.Println("关注失败" + err.Error())
		return
	}

	actionType := c.Query("action_type") //1-关注，2-取消关注

	if actionType == "1" {
		//关注操作
		err := service.Follow(activeUserId, toUserId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "关注失败"})
			fmt.Println("关注失败" + err.Error())
			return
		}

		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "关注成功！"})
		return
	} else if actionType == "2" {
		//取关操作
		err := service.Unfollow(activeUserId, int(toUserId))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "取关失败"})
			fmt.Println("取关失败" + err.Error())
			return
		}

		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "取关成功！"})
		return
	}

	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "关注失败"})
	fmt.Println("actionType错误")
	return
}

// FollowList 关注列表
func FollowList(c *gin.Context) {

	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
		return
	}
	activeUserId := userIdInterface.(int)

	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
		fmt.Println("操作失败" + err.Error())
		return
	}

	followerList, err := service.ListFollowee(userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
		fmt.Println("操作失败" + err.Error())
		return
	}

	var userList = make([]User, len(*followerList))
	point := 0 //videoList的指针
	for _, user := range *followerList {

		followerCount, err := service.CountFollower(int(user.Id))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
			fmt.Println("操作失败" + err.Error())
			return
		}

		followeeCount, err := service.CountFollowee(int(user.Id))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
			fmt.Println("操作失败" + err.Error())
			return
		}

		isFollow, err := service.CheckFollowStatus(activeUserId, int(user.Id))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
			fmt.Println("操作失败" + err.Error())
			return
		}

		userList[point] = User{
			Id:            user.Id,
			Name:          user.UserName,
			FollowCount:   followeeCount,
			FollowerCount: followerCount,
			IsFollow:      isFollow,
		}
		point++
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "成功获取关注列表",
		},
		UserList: userList,
	})
	return

}

// FollowerList 粉丝列表
func FollowerList(c *gin.Context) {

	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败"})
		return
	}
	activeUserId := userIdInterface.(int)

	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
		fmt.Println(err.Error())
		return
	}

	followerList, err := service.ListFollower(userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
		fmt.Println(err.Error())
		return
	}

	var userList = make([]User, len(*followerList))
	point := 0 //videoList的指针
	for _, user := range *followerList {

		followerCount, err := service.CountFollower(int(user.Id))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
			fmt.Println(err.Error())
			return
		}

		followeeCount, err := service.CountFollowee(int(user.Id))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
			fmt.Println(err.Error())
			return
		}

		isFollow, err := service.CheckFollowStatus(activeUserId, int(user.Id))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
			fmt.Println(err.Error())
			return
		}

		userList[point] = User{
			Id:            user.Id,
			Name:          user.UserName,
			FollowCount:   followeeCount,
			FollowerCount: followerCount,
			IsFollow:      isFollow,
		}
		point++
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "成功获取关注列表",
		},
		UserList: userList,
	})
	return
}
