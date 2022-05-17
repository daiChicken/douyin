package mysql

/**
 * @author  Simon5ei
 * @date  2022/5/16 21:22
 * @version  1.0
 * @description
 */

import (
	"BytesDanceProject/model"
)

// GetUser 获取用户信息
func GetUser(username string) model.User {
	var user model.User
	db.Where("username = ?", username).First(&user)
	return user
}

// IsExist 判定用户是否存在，用于能否使用该用户名进行注册操作等功能
func IsExist(username string) bool {
	var user model.User
	db.Where("username = ?", username).First(&user)
	if user.Id == 0 {
		return false
	}
	return true
}

// InsertUser 将创建的用户插入数据库，id为主键自增
func InsertUser(auser model.User) error {
	db.Create(&auser)
	return nil
}

//VerifyPwd 验证用户密码并获取用户id，用于登陆验证
func VerifyPwd(username string, pwd string) bool {
	var user model.User
	db.Where("username = ?", username).First(&user)
	if user.Password == pwd {
		return true
	}
	return false
}
