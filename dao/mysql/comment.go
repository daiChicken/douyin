package mysql

import "BytesDanceProject/model"

// InsertComment 向数据库中插入一条评论数据
func InsertComment(comment *model.Comment) error {
	if err := db.Create(&comment).Error; err != nil {
		return err
	}

	return nil
}

// UpdateStatusById 根据comment改变is_deleted的值
func UpdateStatusById(comment *model.Comment) error {
	err := db.Model(comment).Update("is_deleted", 1).Error
	if err != nil {
		return err
	}

	return nil
}

// GetCommentById 根据id查询comment
func GetCommentById(commentId int) (*model.Comment, error) {

	comment := new(model.Comment)
	// 查询指定的某条记录(仅当主键为整型时可用)
	err := db.First(&comment, 10).Error
	//// SELECT * FROM users WHERE id = 10;
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// ListCommentByVideoId 根据视频id获取评论列表
func ListCommentByVideoId(videoId int) (*[]model.Comment, error) {

	var commentList []model.Comment
	err := db.Where("video_id = ?", videoId).Order("id desc").Find(&commentList).Error
	if err != nil {
		return nil, err
	}
	return &commentList, nil
}
