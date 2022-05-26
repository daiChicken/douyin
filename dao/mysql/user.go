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
func InsertUser(auser model.User) (int64, error) {
	tx := db.Begin() //开启事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return 0, err
	}
	if err := tx.Create(&auser).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	return auser.Id, tx.Commit().Error
}

//VerifyPwd 验证用户密码并获取用户id，用于登陆验证
func VerifyPwd(username string, pwd string) (int64, bool) {
	var user model.User
	db.Where("username = ?", username).First(&user)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd))
	if err != nil {
		return 0, false
	}
	return user.Id, true
}

func GetUserByID(userId int) (model.User, bool) {
	var user model.User
	db.Where("id = ?", userId).First(&user)
	if user.Id == 0 {
		return user, false
	}
	return user, true
}
