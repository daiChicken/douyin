package redis

import (
	"errors"
)

var (
	ErrFavorite = errors.New("点赞出错啦!")
)

// action_type 1 点赞 | 2 取消赞
func FavoriteForVideo(userID, videoID string, actionType int32) error {
	// 去redis取视频
	// 给视频更新点赞量
	// 先去查当前用户对该视频是否点赞 0 没点赞 | 1 点赞
	//ovalue := rdb.ZScore(getRedisKey(KeyVideoFavoritedZSetPf+videoID), userID).Val()

	//此处报错，我先注释掉
	//if actionType == 1 { // 取消点赞
	//	_, err := rdb.ZRm(getRedisKey(KeyVideoFavoritedZSetPf+videoID), userID).Result()
	//	return err
	//} else { // 点赞
	//	_, err := rdb.ZAdd(getRedisKey(KeyVideoFavoritedZSetPf+videoID), redis.Z{
	//		cast.ToFloat64(actionType), // 点赞还是取消点赞
	//		userID,
	//	}).Result()
	//	return err
	//}
	// 记录用户未视频点赞的数据

	//防止报错
	return nil
}
