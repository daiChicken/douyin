package mysql

import (
	"BytesDanceProject/model"
	"gorm.io/gorm"
)

// InsertComment 插入一条新的评论
func InsertComment(comment *model.Comment) (*gorm.DB, error) {
	dbWithTransaction := db.Begin() //开启事务
	if err := dbWithTransaction.Create(&comment).Error; err != nil {
		return nil, err
	}
	return dbWithTransaction, nil
}

// UpdateCommentStatus 将is_deleted字段的值改为1
func UpdateCommentStatus(commentId int) (*gorm.DB, error) {
	dbWithTransaction := db.Begin() //开启事务
	err := dbWithTransaction.Table("comment").Where("id = ? ", commentId).Update("is_deleted", 1).Error
	if err != nil {
		return nil, err
	}
	return dbWithTransaction, nil
}

// GetComment 获取指定的评论
func GetComment(commentId int) (*model.Comment, error) {
	var comment model.Comment
	err := db.Where("id = ? ", commentId).Find(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// ListCommentDESCByCreateDate 根据创建时间倒序获取该视频所有未被删除的评论
func ListCommentDESCByCreateDate(videoId int) (*[]model.Comment, error) {
	var commentList []model.Comment
	err := db.Where("video_id = ? AND is_deleted = ?", videoId, 0).Order("create_date desc").Find(&commentList).Error
	if err != nil {
		return nil, err
	}
	return &commentList, nil
}

// CountCountByVideoId 统计视频的评论数
func CountCountByVideoId(videoId int) (int64, error) {
	var count int64
	err := db.Table("comment").Where("video_id = ? AND is_deleted = ?", videoId, 0).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
