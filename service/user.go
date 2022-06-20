package service

import (
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/model"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

/**
 * @author  Simon5ei
 * @date  2022/5/17 02:13
 * @version  1.0
 * @description
 */

// Register 尝试注册
func Register(username string, password string) (int64, error) {

	//空值处理
	if username == "" {
		return 0, errors.New("用户名不能为空！")
	}
	if password == "" {
		return 0, errors.New("密码不能为空！")
	}

	//用户名存在校验
	_, err := mysql.GetUserByUsername(username)
	if err == nil {
		//不出错说明成功查到数据，说明用户名已经存在了
		return 0, errors.New("用户名已存在！")
	} else if err != nil && err != gorm.ErrRecordNotFound {
		//出错了，但是出的错并不是没有查到数据这个错误
		return 0, err
	}

	EncryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //密码加密处理
	if err != nil {
		return 0, err
	}

	user := model.User{
		UserName:  username,
		Password:  string(EncryptedPassword),
		AvatarUrl: "",
		Nickname:  "",
	}

	userId, err := mysql.InsertUser(user)
	if err != nil {
		return 0, err
	}

	return userId, nil

}

// Login 验证登陆
func Login(username string, password string) (int64, error) {

	//空值处理
	if username == "" {
		return 0, errors.New("用户名不能为空！")
	}
	if password == "" {
		return 0, errors.New("密码不能为空！")
	}

	//验证用户存在
	user, err := mysql.GetUserByUsername(username)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, errors.New("密码不正确！")
	}

	return user.Id, nil
}

func GetUserInfo(userId int64) (int64, int64) {
	followerTable, _ := mysql.GetFollower(userId)
	FollowerCount := int64(len(followerTable))
	followTable, _ := mysql.GetFollowed(userId)
	FollowCount := int64(len(followTable))
	return FollowerCount, FollowCount
}

// 根据 userId 获取 user 的所有能在单表中得到的信息
func GetUserByID(userId int) (*model.User, error) {

	return mysql.GetUser(userId)
}
