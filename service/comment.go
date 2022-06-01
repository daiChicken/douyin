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
		CreateDate: time.Now().Format("01-02"),
		IsDeleted:  0,
	}

	//向数据库中插入评论数据
	err := mysql.InsertComment(&NewComment)
	if err != nil {
		return nil, err
	}

	return &NewComment, nil
}

// DeleteComment 删除评论
func DeleteComment(commentId int) error {
	//需要删除评论时，将评论的is_deleted字段修改为1
	comment, err := mysql.GetCommentById(commentId)
	if err != nil {
		return err
	}

	err = mysql.UpdateStatusById(comment)
	if err != nil {
		return err
	}

	return nil
}

// ListCommentByVideoId 获取id为videoId的视频的所有评论
func ListCommentByVideoId(videoId int) (*[]model.Comment, error) {
	commentList, err := mysql.ListCommentByVideoId(videoId)
	if err != nil {
		return nil, err
	}

	return commentList, err
}
