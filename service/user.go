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
func Register(username string, password string) bool {
	if mysql.IsExist(username) == false {
		_ = mysql.InsertUser(model.User{
			UserName: username,
			Password: password,
		})
		return true
	}
	return false
}

// VerifyLogin 验证登陆
func VerifyLogin(username string, password string) (bool) {
	Flag := mysql.VerifyPwd(username, password)
	return Flag
}

func FindUser(username string) model.User {
	return mysql.GetUser(username)
}