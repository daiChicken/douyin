package service

import (
	"BytesDanceProject/dao/mysql"
	rds "BytesDanceProject/dao/redis"
	"BytesDanceProject/model"
	"github.com/robfig/cron"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

/*
分析：
	1、点赞直接存到redis
	2、视频 - 用户 点赞关系：存到 mysql - video_like_relation

*/

//用户点赞的操作
func FavoriteAction(p model.FavoriteRequest) error {
	// 点赞存入 redis
	zap.L().Debug("FavoriteAction",
		zap.Int64("userID", p.UserID),
		zap.Int64("videoID", p.VideoID),
		zap.Int32("action_type", p.ActionType))
	//
	if err := rds.FavoriteForVideo(cast.ToString(p.UserID), cast.ToString(p.VideoID), p.ActionType); err != nil {
		return err
	}
	favoriteActionData := model.VideoLikeRelation{
		UserId:  p.UserID,
		VideoId: p.VideoID,
		Status:  p.ActionType,
	}
	if err := mysql.CreateFavoriteAction(&favoriteActionData); err != nil {
		return err
	}

	c := cron.New()
	c.AddFunc("0 */2 * * *", FavoriteCron)

	return nil
}

// 点赞列表
func FavoriteList(p model.FavoriteListRequest) error {
	//uid := p.UserID
	// 根据 uid 找到该 uid 所有点赞过的视频  videoID/AuthorID
	// 根据 videoID 找到该视频所有的信息——从video表中拉出 play_url/ cover_url/ favorite_count/ comment_count/ is_favorite
	// 根据 authorID 找到该 anthorID 的相关信息 nickname/followCount/followerCount/IsFollow

	return nil
}

func FavoriteCron() {

}
