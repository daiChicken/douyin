package service

import (
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/dao/redis"
	"BytesDanceProject/model"
	"encoding/json"
	"go.uber.org/zap"
	"gorm.io/gorm/utils"
	"strconv"
)

// Business logic that handles follow

const (
	KeyFollow    = "follow"
	KeyFollower  = "follower"
	KeyFansList  = "fans"
	KeyFocusList = "focus"
)

func RelationAction(p *model.RelationAction) {

	// 直接更新缓存，数据库由定时任务持久化
	//取出本id和对应id下的姓名
	idname, toidname := mysql.GetNameFolAndFoler(p)
	redis.UpdateCache(p, idname, toidname)
}

// GetFollowList 获取用户关注的所用用户列表
func GetFollowList(p *model.FollowListRE) (datas []model.UserFocus, err error) {
	// 从redis 取得id切片
	ids, err := redis.GetMsgListByID(KeyFocusList, p.UserID)
	if err != nil {
		zap.L().Error("redis.GetFollowIDList(p.UserID) failed err := ", zap.Error(err))
		return nil, err
	}
	// 先查询缓存
	if len(ids) == 0 {
		//  以下是缓存没有的情况下 查询数据库的逻辑，并重新设置缓存
		datas, err = mysql.GetFollowList(p.UserID)
		if err != nil {
			zap.L().Error("mysql.GetFollowList(p.UserID) failed err := ", zap.Error(err))
			return nil, err
		}
		//重新设置缓存
		if err = redis.SetUserCache(KeyFollow, p.UserID, datas); err != nil {
			zap.L().Error("redis.SetUserFollowCache(p.UserID,datas) failed err := ", zap.Error(err))
			return nil, err
		}

	} else {
		// 查询缓存
		msg := model.ListUserMsg{}
		for _, id := range ids {
			if err = json.Unmarshal([]byte(id), &msg); err != nil {
				zap.L().Error("json.Unmarshal([]byte(id),msg) failed err : ", zap.Error(err))
				return nil, err
			}
			followc := redis.GetCountByID(KeyFollow, msg.ID)
			followerc := redis.GetCountByID(KeyFollower, msg.ID)
			if followc == 0 || followerc == 0 {
				// 查数据库
				followc, followerc, err = mysql.GetCountByID(msg.ID)
				if err != nil {
					return nil, err
				}
			}
			data := model.UserFocus{
				ID:            msg.ID,
				UserName:      msg.UserName,
				FollowCount:   followc,
				FollowerCount: followerc,
				IsFollow:      true,
			}

			datas = append(datas, data)
		}
	}

	return datas, err
}

// GetFollowerList 获取用户粉丝的所用用户列表
func GetFollowerList(p *model.FollowListRE) (datas []model.UserFocus, err error) {
	// 从redis 取得id切片
	ids, err := redis.GetMsgListByID(KeyFansList, p.UserID)
	if err != nil {
		zap.L().Error("redis.GetFollowIDList(p.UserID) failed err := ", zap.Error(err))
		return nil, err
	}
	// 根据id切片取各id的信息（关注、粉丝数）
	if len(ids) == 0 {
		// 说明没有缓存 则查询数据库并设置缓存
		datas, err = mysql.GetFollowerList(p.UserID)
		if err != nil {
			zap.L().Error("mysql.GetFollowerList(p.UserID) failed err := ", zap.Error(err))
			return nil, err
		}

		// 重新设置缓存
		if err = redis.SetUserCache(KeyFollower, p.UserID, datas); err != nil {
			zap.L().Error("redis.SetUserFocusCache(p.UserID,datas) failed err := ", zap.Error(err))
			return nil, err
		}
	} else {

		msg := model.ListUserMsg{}
		fids, err := redis.GetMsgListByID(KeyFocusList, p.UserID)
		var judge map[string]bool
		if len(fids) == 0 {
			// 去数据库取
			IDlist, _ := mysql.GetFollowList(p.UserID)
			for _, id := range IDlist {
				fids = append(fids, strconv.Itoa(int(id.ID)))
			}
			judge = make(map[string]bool, len(fids))
			for _, id := range fids {
				// 加入hash
				judge[utils.ToString(id)] = true
			}
		} else {
			judge = make(map[string]bool, len(fids))
			for _, id := range fids {
				if err = json.Unmarshal([]byte(id), &msg); err != nil {
					zap.L().Error("json.Unmarshal([]byte(id),msg) failed err : ", zap.Error(err))
					return nil, err
				}
				// 加入hash
				judge[utils.ToString(msg.ID)] = true
			}
		}

		for _, id := range ids {
			if err = json.Unmarshal([]byte(id), &msg); err != nil {
				zap.L().Error("json.Unmarshal([]byte(id),msg) failed err : ", zap.Error(err))
				return nil, err
			}
			followc := redis.GetCountByID(KeyFollow, msg.ID)
			followerc := redis.GetCountByID(KeyFollower, msg.ID)
			if followc == 0 || followerc == 0 {
				// 查数据库
				followc, followerc, err = mysql.GetCountByID(msg.ID)
				if err != nil {
					return nil, err
				}
			}
			data := model.UserFocus{
				ID:            msg.ID,
				UserName:      msg.UserName,
				FollowCount:   followc,
				FollowerCount: followerc,
				IsFollow:      judge[utils.ToString(msg.ID)],
			}
			datas = append(datas, data)
		}
	}

	return datas, err
}
