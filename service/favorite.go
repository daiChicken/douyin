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

//用户点赞的操作
//func FavoriteAction(p model.FavoriteRequest) error {
//	// 点赞存入 redis
//	zap.L().Debug("FavoriteAction",
//		zap.Int64("userID", p.UserID),
//		zap.Int64("videoID", p.VideoID),
//		zap.Int32("action_type", p.ActionType))
//	//
//	favoriteActionData := model.VideoLikeRelation{
//		UserId:  p.UserID,
//		VideoId: p.VideoID,
//		Status:  p.ActionType,
//	}
//	spew.Dump("==============================", favoriteActionData)
//	// 存到 mysql
//	dbWithTransaction, err := mysql.CreateFavoriteAction(&favoriteActionData)
//	if err != nil {
//		return err
//	}
//
//	//存到 redis:  用户：视频：点赞状态    视频：点赞数量
//	keyUserToVideo := rds.GetFavoriteKey(p.UserID, p.VideoID)
//	//keyVideoNum := rds.GetUserFavoriteKey(p.VideoID)
//
//	// 点赞
//	if p.ActionType == 1 {
//		rdb.SAdd(keyUserToVideo, p.ActionType)
//		if err := rdb.SAdd(keyUserToVideo, p.ActionType).Err(); err != nil {
//			fmt.Println("err = ", err)
//		}
//		dbWithTransaction.Commit()
//		return nil
//	}
//	rdb.SRem(keyUserToVideo, p.ActionType)
//
//	dbWithTransaction.Commit() //提交事务
//
//	return nil
//}

// LikeAction 点赞操作
func LikeAction(userId int, videoId int) error {

	videoLikeKey := tool.GetVideoLikeKey(videoId)
	err := redis.AddUserToVideoSet(videoLikeKey, userId)
	if err != nil {
		return err
	}

	userLikeKey := tool.GetUserLikeKey(userId)
	err = redis.AddVideoToUserSet(userLikeKey, videoId)
	if err != nil {
		return err
	}

	return nil
}

// CancelLikeAction 取消点赞操作
func CancelLikeAction(userId int, videoId int) error {

	videoLikeKey := tool.GetVideoLikeKey(videoId)
	err := redis.RemoveUserFromVideoSet(videoLikeKey, userId)
	if err != nil {
		return err
	}

	userLikeKey := tool.GetUserLikeKey(userId)
	err = redis.RemoveVideoFromUserSet(userLikeKey, videoId)
	if err != nil {
		return err
	}

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
