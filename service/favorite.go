package service

import (
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/dao/redis"
	"BytesDanceProject/model"
	"BytesDanceProject/tool"
	"strconv"
)

type ListFavoriteMsg struct {
	userID        int64  `json:"userID" db:"userID"`
	videoID       string `json:"user_name" db:"user_name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
}

// LikeAction 点赞操作
func LikeAction(userId int, videoId int) error {

	rdb := redis.GetRDB()
	txPipeline := rdb.TxPipeline() //开启事务

	videoLikeKey := tool.GetVideoLikeKey(videoId)
	err := redis.AddUserToVideoSet(videoLikeKey, userId, txPipeline)
	if err != nil {
		txPipeline.Discard() //回滚事务
		return err
	}

	userLikeKey := tool.GetUserLikeKey(userId)
	err = redis.AddVideoToUserSet(userLikeKey, videoId, txPipeline)
	if err != nil {
		txPipeline.Discard() //回滚事务
		return err
	}

	txPipeline.Exec() //提交事务

	return nil
}

// CancelLikeAction 取消点赞操作
func CancelLikeAction(userId int, videoId int) error {

	rdb := redis.GetRDB()
	txPipeline := rdb.TxPipeline() //开启事务

	videoLikeKey := tool.GetVideoLikeKey(videoId)
	err := redis.RemoveUserFromVideoSet(videoLikeKey, userId, txPipeline)
	if err != nil {
		txPipeline.Discard() //回滚事务
		return err
	}

	userLikeKey := tool.GetUserLikeKey(userId)
	err = redis.RemoveVideoFromUserSet(userLikeKey, videoId, txPipeline)
	if err != nil {
		txPipeline.Discard() //回滚事务
		return err
	}

	txPipeline.Exec() //提交事务

	return nil
}

func CountLike(videoId int) (int64, error) {
	key := tool.GetVideoLikeKey(videoId)
	count, err := redis.CountLike(key)
	if err != nil {
		return 0, err
	}
	return count, err
}

func GetLikeStatus(videoId int, userId int) (bool, error) {
	key := tool.GetVideoLikeKey(videoId)
	isLike, err := redis.GetLikeStatus(key, userId)
	if err != nil {
		return false, err
	}
	return isLike, err

}

func ListLikeVideo(userId int) (*[]model.Video, error) {
	userLikeKey := tool.GetUserLikeKey(userId)
	videoIds, err := redis.ListLikedVideo(userLikeKey)
	if err != nil {
		return nil, err
	}

	var videoList = make([]model.Video, len(videoIds))
	for idx, idString := range videoIds {
		id, _ := strconv.ParseInt(idString, 10, 64)
		videoFromMysql, err := mysql.GetVideoById(int(id))

		if err != nil {
			return nil, err
		}
		videoList[idx] = *videoFromMysql
	}

	return &videoList, nil
}
