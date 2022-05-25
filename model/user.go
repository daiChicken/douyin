package model

/**
 * @author  Simon5ei
 * @date  2022/5/16 21:00
 * @version  1.0
 * @description
 */

type User struct {
	Id        int64  `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	UserName  string `gorm:"column:username"`
	Password  string `gorm:"column:password"`
	AvatarUrl string `gorm:"column:avatar_url"`
	Nickname  string `gorm:"column:nickname"`
}

func (User) TableName() string {
	return "User"
}
