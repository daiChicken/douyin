package redis

import (
	"github.com/spf13/cast"
)

// 用户-视频-点赞状态
func GetFavoriteKey(userID int64, videroID int64) string {
	return KeyPrefix + cast.ToString(userID) + ":" + cast.ToString(videroID)
}

func GetUserFavoriteKey(vedioID int64) string {
	return KeyPrefix + cast.ToString(vedioID)
}
