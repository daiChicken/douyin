package service

import (
	"BytesDanceProject/controller"
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/model"
)

/**
 * @author  Simon5ei
 * @date  2022/5/17 02:13
 * @version  1.0
 * @description
 */

// Register 尝试注册
func Register(username string, password string) (int64, bool) {
	if mysql.IsExist(username) == false {
		Id, _ := mysql.InsertUser(model.User{
			UserName: username,
			Password: password,
		})
		return Id, true
	}
	return 0, false
}

// VerifyLogin 验证登陆
func VerifyLogin(username string, password string) (int64, bool) {
	return mysql.VerifyPwd(username, password)
}

func GetUserInfo(user controller.User) controller.User {
	followerTable, _ := mysql.GetFollower(user.Id)
	user.FollowerCount = int64(len(followerTable))
	followTable, _ := mysql.GetFollowed(user.Id)
	user.FollowCount = int64(len(followTable))
	return user
}
