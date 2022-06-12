package service

import (
	"BytesDanceProject/dao/mysql"
	rds "BytesDanceProject/dao/redis"
	"BytesDanceProject/model"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type ListFavoriteMsg struct {
	userID        int64  `json:"userID" db:"userID"`
	videoID       string `json:"user_name" db:"user_name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
}

var rdb *redis.Client

//用户点赞的操作
func FavoriteAction(p model.FavoriteRequest) error {
	// 点赞存入 redis
	zap.L().Debug("FavoriteAction",
		zap.Int64("userID", p.UserID),
		zap.Int64("videoID", p.VideoID),
		zap.Int32("action_type", p.ActionType))
	//
	favoriteActionData := model.VideoLikeRelation{
		UserId:  p.UserID,
		VideoId: p.VideoID,
		Status:  p.ActionType,
	}
	spew.Dump("==============================", favoriteActionData)
	// 存到 mysql
	dbWithTransaction, err := mysql.CreateFavoriteAction(&favoriteActionData)
	if err != nil {
		return err
	}

	//存到 redis:  用户：视频：点赞状态    视频：点赞数量
	keyUserToVideo := rds.GetFavoriteKey(p.UserID, p.VideoID)
	//keyVideoNum := rds.GetUserFavoriteKey(p.VideoID)

	// 点赞
	if p.ActionType == 1 {
		rdb.SAdd(keyUserToVideo, p.ActionType)
		if err := rdb.SAdd(keyUserToVideo, p.ActionType).Err(); err != nil {
			fmt.Println("err = ", err)
		}
		dbWithTransaction.Commit()
		return nil
	}
	rdb.SRem(keyUserToVideo, p.ActionType)

	dbWithTransaction.Commit() //提交事务

	return nil
}

// 点赞列表
func FavoriteList(p model.FavoriteListRequest) error {
	// 根据id获取喜爱列表

	return nil
}

//func FavoriteCron() {
//
//}
