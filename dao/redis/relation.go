package redis

import (
	"BytesDanceProject/model"
	"encoding/json"
	"fmt"
	"strconv"
)


// 给go-redis 包自动调用的 序列号与反序列化 函数

func (l *ListUserMsg) MarshalBinary() ([]byte, error) {
	return json.Marshal(l)
}

func (l *ListUserMsg) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, l)
}

// ListUserMsg 关注\粉丝列表的信息
type ListUserMsg struct {
	ID 				int64 			`json:"id" db:"id"`
	UserName		string 			`json:"user_name" db:"user_name"`
	FollowCount 	int64 			`json:"follow_count"`
	FollowerCount	int64 			`json:"follower_count"`
}


// UpdateCache 更新缓存
func UpdateCache(p *model.RelationAction,idname string,toidname string){
	// A关注B 则A的关注列表多一个人，B的粉丝列表多一个人
	// A取消关注B 则A的关注列表少一人  B的粉丝列表少一人
	keyfollow := getRedisKey(KeyFollowPrefix+strconv.Itoa(int(p.UserID)))
	keyfollower := getRedisKey(KeyFollowerPrefix+strconv.Itoa(int(p.ToUserID)))
	followmsg := &ListUserMsg{
		ID:            p.ToUserID,
		UserName:      toidname,
		FollowCount: GetCountByID("follow",p.ToUserID),
		FollowerCount: GetCountByID("follower",p.ToUserID),
	}
	followermsg := &ListUserMsg{
		ID:            	p.UserID,
		UserName: 		idname,
		FollowCount: GetCountByID("follow",p.UserID),
		FollowerCount: GetCountByID("follower",p.UserID),
	}
	if p.ActionType == 1{
		// 关注
		// 先使用无序集合set 其放入顺序是在头部放入刚好可以满足需求
		followmsg.FollowerCount ++
		rdb.SAdd(keyfollow,followmsg)
		if err := rdb.SAdd(keyfollower,followermsg).Err();err != nil{
			fmt.Println("err = ", err)
		}
		return
	}
	rdb.SRem(keyfollow,followmsg)
	rdb.SRem(keyfollower,followermsg)
	return
}

// GetCountByID 根据id获取该id下的粉丝或关注数量
func GetCountByID(tkey string,userid int64)int64{
	key := getRedisKey(KeyFollowPrefix + strconv.Itoa(int(userid)))
	if tkey == "follower"{
		key = getRedisKey(KeyFollowerPrefix + strconv.Itoa(int(userid)))
	}
	return rdb.SCard(key).Val()
}



// GetMsgListByID 根据id 获取粉丝或关注 列表
func GetMsgListByID(key string,userid int64)([]string,error){
	lkey := getRedisKey(KeyFollowPrefix + strconv.Itoa(int(userid)))
	if key == "fans"{
		lkey = getRedisKey(KeyFollowerPrefix + strconv.Itoa(int(userid)))
	}
	return rdb.SMembers(lkey).Result()
}





// SetUserCache 设置缓存
func SetUserCache(key string , uid int64,datas []model.UserFocus)(err error){
	skey := getRedisKey(KeyFollowerPrefix + strconv.Itoa(int(uid)))
	if key == "follow" {
		skey = getRedisKey(KeyFollowPrefix + strconv.Itoa(int(uid)))
	}
	pipeline := rdb.Pipeline()
	for _,data := range datas{
		pipeline.SAdd(skey,&ListUserMsg{
			ID:            data.ID,
			UserName:      data.UserName,
			FollowCount:   data.FollowCount,
			FollowerCount: data.FollowerCount,
		})
	}
	_,err =pipeline.Exec()
	return
}