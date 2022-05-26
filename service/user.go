package service

import (
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

func GetUserInfo(userId int64) (int64, int64) {
	followerTable, _ := mysql.GetFollower(userId)
	FollowerCount := int64(len(followerTable))
	followTable, _ := mysql.GetFollowed(userId)
	FollowCount := int64(len(followTable))
	return FollowerCount, FollowCount
}

// 根据 userId 获取 user 的所有能在单表中得到的信息
func GetUserByID(userId int) (model.User, bool) {
	return mysql.GetUserByID(userId)
}
