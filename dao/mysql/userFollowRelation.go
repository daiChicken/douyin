package mysql

import "BytesDanceProject/model"

/**
 * @author  Simon5ei
 * @date  2022/5/20 21:04
 * @version  1.0
 * @description
 */

func GetFollower(userId int64) ([]model.UserFollowRelation, error) {
	var users []model.UserFollowRelation
	err := db.Where("user_id = ?", userId).Where("status = ?", 1).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetFollowed(userId int64) ([]model.UserFollowRelation, error) {
	var users []model.UserFollowRelation
	err := db.Where("followed_user_id = ?", userId).Where("status = ?", 1).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
