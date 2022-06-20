package service

import (
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/dao/redis"
	"BytesDanceProject/model"
	"BytesDanceProject/tool"
	"strconv"
)

/**
 * @author  daijizai
 * @date  2022/6/20 18:21
 * @version  1.0
 * @description
 */

// Follow 关注操作
func Follow(activeUserId int, targetUserId int) error {

	rdb := redis.GetRDB()
	txPipeline := rdb.TxPipeline() //开启事务

	followerKey := tool.GetFollowerKey(targetUserId)
	err := redis.AddToSet(followerKey, activeUserId, txPipeline)
	if err != nil {
		txPipeline.Discard() //回滚事务
		return err
	}

	followeeKey := tool.GetFolloweeKey(activeUserId)
	err = redis.AddToSet(followeeKey, targetUserId, txPipeline)
	if err != nil {
		txPipeline.Discard() //回滚事务
		return err
	}

	txPipeline.Exec() //提交事务

	return nil
}

// Unfollow 取消关注
func Unfollow(activeUserId int, targetUserId int) error {

	rdb := redis.GetRDB()
	txPipeline := rdb.TxPipeline() //开启事务

	followerKey := tool.GetFollowerKey(targetUserId)
	err := redis.RemoveFromSet(followerKey, activeUserId, txPipeline)
	if err != nil {
		txPipeline.Discard() //回滚事务
		return err
	}

	followeeKey := tool.GetFolloweeKey(activeUserId)
	err = redis.RemoveFromSet(followeeKey, targetUserId, txPipeline)
	if err != nil {
		txPipeline.Discard() //回滚事务
		return err
	}

	txPipeline.Exec() //提交事务

	return nil
}

// ListFollower 获取粉丝列表
func ListFollower(userId int) (*[]model.User, error) {
	followerKey := tool.GetFollowerKey(userId)
	userIds, err := redis.Smembers(followerKey)
	if err != nil {
		return nil, err
	}

	var userList = make([]model.User, len(userIds))
	for idx, idString := range userIds {
		id, _ := strconv.ParseInt(idString, 10, 64)
		userFromMysql, err := mysql.GetUser(int(id))
		if err != nil {
			return nil, err
		}

		userList[idx] = *userFromMysql
	}

	return &userList, nil
}

// CountFollowee 查询关注总数
func CountFollowee(userId int) (int64, error) {
	followeeKey := tool.GetFolloweeKey(userId)
	count, err := redis.CountSet(followeeKey)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CountFollower 查询粉丝总数
func CountFollower(userId int) (int64, error) {
	followerKey := tool.GetFollowerKey(userId)
	count, err := redis.CountSet(followerKey)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CheckFollowStatus 检查关注状态
func CheckFollowStatus(activeUserId int, targetUserId int) (bool, error) {
	followerKey := tool.GetFollowerKey(targetUserId)
	userIds, err := redis.Smembers(followerKey)
	if err != nil {
		return false, err
	}
	for _, idString := range userIds {
		id, _ := strconv.Atoi(idString)
		if id == activeUserId {
			//目标用户的粉丝列表里有activeUser
			return true, nil
		}
	}
	return false, nil
}

// ListFollowee 获取关注列表
func ListFollowee(userId int) (*[]model.User, error) {
	followeeKey := tool.GetFolloweeKey(userId)
	userIds, err := redis.Smembers(followeeKey)
	if err != nil {
		return nil, err
	}

	var userList = make([]model.User, len(userIds))
	for idx, idString := range userIds {
		id, _ := strconv.ParseInt(idString, 10, 64)
		userFromMysql, err := mysql.GetUser(int(id))
		if err != nil {
			return nil, err
		}

		userList[idx] = *userFromMysql
	}

	return &userList, nil
}
