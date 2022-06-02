package service

import (
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/model"
	"time"
)

// CreateComment 创建评论
func CreateComment(userId int, videoId int, commentText string) (*model.Comment, error) {

	NewComment := model.Comment{
		UserID:     userId,
		VideoID:    videoId,
		Content:    commentText,
		CreateDate: time.Now(),
		IsDeleted:  0,
		UpdateDate: time.Now(),
	}

	err := mysql.InsertComment(&NewComment)
	if err != nil {
		return nil, err
	}

	return &NewComment, nil
}

// DeleteComment 删除评论
func DeleteComment(commentId int) error {

	err := mysql.UpdateCommentStatus(commentId)
	if err != nil {
		return err
	}
	return nil
}

// ListComment 获取videoId的所有未被删除的评论
func ListComment(videoId int) (*[]model.Comment, error) {

	commentList, err := mysql.ListCommentDESCByCreateDate(videoId)
	if err != nil {
		return nil, err
	}
	return commentList, err
}
