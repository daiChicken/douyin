package model

import "time"

/**
 * @author  Simon5ei
 * @date  2022/5/20 21:00
 * @version  1.0
 * @description
 */

type UserFollowRelation struct {
	Id             int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	UserId         int64     `gorm:"column:user_id"`
	FollowedUserId int64     `gorm:"column:followed_user_id"`
	Status         int64     `gorm:"column:status"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (UserFollowRelation) TableName() string {
	return "user_follow_relation"
}
