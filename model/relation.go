package model

/**
 * @author  jiangjingyu
 * @date  2022/5/18 16:53
 * @version  1.0
 * @description  存放关注相关的数据模型
 */


type RelationAction struct {
	UserID 			int64 		`json:"user_id" binding:"required"`		// 用户id
	Token 			string 		`json:"token" binding:"required"`			// 用户鉴权token
	ToUserID 		int64 		`json:"to_user_id" binding:"required"`		// 对方用户id
	ActionType 		int32 		`json:"action_type" binding:"required,oneof=1 2"`	// 1-关注，2-取消关注
}

type FocusCount struct {
	FollowCount		int64 		`json:"follow_count" gorm:"column:follow_count"`
	FollowerCount 	int64 		`json:"follower_count" gorm:"column:follower_count"`
}

func (FocusCount) TableName() string {
	return "focus_count"
}

type Follow struct {
	ID 				int64 			`json:"id" gorm:"column:id"`
	FollowID 		int64 			`json:"follow_id" gorm:"column:follow_id"`
	IsFollow 		bool 			`json:"is_follow" gorm:"column:is_follow"`
}

func (Follow) TableName() string{
	return "focus_follow"
}

type Follower struct {
	ID 				int64 			`json:"id" gorm:"column:id"`
	FollowerID 		int64 			`json:"follow_id" gorm:"column:follower_id"`
	IsFollow 		bool 			`json:"is_follow" gorm:"column:is_follow"`
}

func (Follower) TableName() string{
	return "focus_follower"
}

// UserFocus 用户关注列表
type UserFocus struct {
	ID 				int64 			`json:"id" gorm:"column:id"`
	UserName		string 			`json:"user_name" gorm:"column:username"`
	FollowCount		int64 			`json:"follow_count" gorm:"column:follow_count"`
	FollowerCount	int64 			`json:"follower_count" gorm:"column:follower_count"`
	IsFollow		bool 			`json:"is_follow" gorm:"column:is_follow"`
}

func(UserFocus) TableName() string{
	return "focus_count"
}

type ListUserMsg struct {
	ID 				int64 			`json:"id" gorm:"column:id"`
	UserName		string 			`json:"user_name" gorm:"column:username"`
	FollowCount		int64 			`json:"follow_count" gorm:"column:follow_count"`
	FollowerCount	int64 			`json:"follower_count" gorm:"column:follower_count"`
}
// FollowListRE 关注列表的请求参数表
type FollowListRE struct {
	UserID 			int64 			`form:"user_id" binding:"required"`
	Token			string 			`form:"token" binding:"required"`
}


// ListUserMsg 关注\粉丝列表的信息
//type ListUserMsg struct {
//	ID 				int64 			`json:"id" db:"id"`
//	UserName		string 			`json:"user_name" db:"user_name"`
//	FollowCount		int64 			`json:"follow_count" db:"follow_count"`
//	FollowerCount	int64 			`json:"follower_count" db:"follower_count"`
//}

// UserName  获取id 和 toid 的姓名字段
type UserName struct {
	ID 				int64 			`json:"id" db:"id" gorm:"column:id"`
	UserName  		string 			`json:"user_name" db:"username" gorm:"column:username"`
}
func (UserName) TableName() string {
	return "user"
}


type IDs struct {
	ID 				int64 			`json:"id"`
	IsFollow		bool 			`json:"is_follow"`
}


