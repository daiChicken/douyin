package mysql

/**
 * @author  Simon5ei
 * @date  2022/5/16 21:22
 * @version  1.0
 * @description
 */

import (
	"BytesDanceProject/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GetUser 获取用户信息(目前还有问题)
func GetUser(username string) (model.User, error) {
	var user model.User
	err := db.Where("username = ?", username).First(&user).Error
	return user, err
}

// IsExist 判定用户是否存在，用于能否使用该用户名进行注册操作等功能
func IsExist(username string) bool {
	var user model.User
	err := db.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
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
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd))
	return err == nil
}
