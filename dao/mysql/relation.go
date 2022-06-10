package mysql

import (
	"BytesDanceProject/model"
	"errors"
	"gorm.io/gorm"
)

var (
	UserNotExist = errors.New("用户不存在")
)

// 关注逻辑
func RelationAction(p *model.RelationAction) error {
	//插入数据库
	// 1、本id的被关注者+1
	// 2、插入数据库
	var focuscount model.FocusCount
	follow := model.Follow{
		ID:       p.UserID,
		FollowID: p.ToUserID,
		IsFollow: true,
	}
	tx := db.Begin()
	rowAffected := tx.Model(&focuscount).Where("id", p.UserID).Update("follow_count", gorm.Expr("follow_count + 1")).RowsAffected
	if rowAffected == 0 {
		// 这里回滚作用不大，因为前面没成功执行什么数据库更新操作，也没什么数据需要回滚。
		return UserNotExist
	}
	if err := tx.Create(&follow).Error; err != nil {
		//插入失败
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// UnRelationAction 取消关注逻辑
func UnRelationAction(p *model.RelationAction) error {
	// 1、本id的被关注者-1
	// 2、删除数据库
	var focuscount model.FocusCount
	tx := db.Begin()
	rowAffected := tx.Model(&focuscount).Where("id", p.UserID).Update("follow_count", gorm.Expr("follow_count - 1")).RowsAffected
	if rowAffected == 0 {
		// 这里回滚作用不大，因为前面没成功执行什么数据库更新操作，也没什么数据需要回滚。
		return UserNotExist
	}
	if err := tx.Where("id = ?", p.UserID).Delete(&model.Follower{}).Error; err != nil {
		//插入失败
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// GetCountByIDs 根据id列表获取各id下的关注和粉丝数
func GetCountByIDs(ids []model.IDs) (data []model.UserFocus, err error) {
	//	sqlStr := `	select id,follow_count,follower_count
	//				from focus_count
	//				where id in (?)
	//				order by FIND_IN_SET(id,?)
	//`
	//	sqlStr := `	select username
	//				from user
	//				where id in (?)
	//				order by FIND_IN_SET(id,?)
	id := make([]int64, len(ids))
	for idx, _ := range ids {
		id[idx] = ids[idx].ID
	}
	if err = db.Where("id in (?)", id).Find(&data).Error; err != nil {
		return nil, err
	}
	usermgs := []model.User{}

	if err = db.Where("id IN (?)", id).Find(&usermgs).Error; err != nil {
		return nil, err
	}

	for idx, _ := range usermgs {
		data[idx].UserName = usermgs[idx].UserName
		data[idx].IsFollow = ids[idx].IsFollow
	}
	return data, err
}

// GetNameFolAndFoler 获取本id和操作id 的姓名字段
func GetNameFolAndFoler(p *model.RelationAction) (string, string) {
	idmsg := &model.UserName{}
	toidmsg := &model.UserName{}
	db.Where("id = ?", p.UserID).Take(idmsg)
	db.Where("id = ?", p.ToUserID).Take(toidmsg)
	return idmsg.UserName, toidmsg.UserName
}

// GetFollowList 根据id 获取该id下的关注信息
func GetFollowList(id int64) (datas []model.UserFocus, err error) {
	// 查询 关注表 该id关注的人的id
	msg := []model.Follow{}
	if err = db.Where("id = ?", id).Find(&msg).Error; err != nil {
		return nil, err
	}
	var ids = []model.IDs{}
	for _, m := range msg {
		ids = append(ids, model.IDs{
			ID:       m.FollowID,
			IsFollow: m.IsFollow,
		})
	}
	// 根据查询出来的id列表 查询 focus count 表获取对应id的关注数和粉丝数
	return GetCountByIDs(ids)
}

// GetFollowerList 根据id 获取该id下的粉丝信息
func GetFollowerList(id int64) (datas []model.UserFocus, err error) {
	// 查询 粉丝表 该id关注的人的id和  是否关注？
	msg := []model.Follower{}
	if err = db.Where("id = ?", id).Find(&msg).Error; err != nil {
		return nil, err
	}
	var ids []model.IDs
	for _, m := range msg {
		ids = append(ids, model.IDs{
			ID:       m.FollowerID,
			IsFollow: m.IsFollow,
		})
	}
	// 根据查询出来的id列表 查询 focus count 表获取对应id的关注数和粉丝数
	return GetCountByIDs(ids)
}

// GetCountByID 根据id 取得该id的 关注数
func GetCountByID(id int64) (int64, int64, error) {
	var msg []*model.FocusCount

	err := db.Where("id = ?", id).Find(&msg).Error
	if err != nil {
		return 0, 0, err
	}
	if len(msg) == 0 {
		return 0, 0, nil
	}
	return msg[0].FollowCount, msg[0].FollowerCount, nil
}
