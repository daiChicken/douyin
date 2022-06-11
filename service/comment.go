package service

import (
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/dao/redis"
	"BytesDanceProject/model"
	"BytesDanceProject/tool"
	"time"
)

// CreateComment 创建评论
func CreateComment(userId int, videoId int, commentText string, userName string) (*model.Comment, error) {

	now := time.Now()                //获取当前时间
	time := time.Unix(now.Unix(), 0) //将时间的精度降低到秒级

	NewComment := model.Comment{
		UserID:     userId,
		VideoID:    videoId,
		Content:    commentText,
		CreateDate: time,
		IsDeleted:  0,
		UpdateDate: time,
		UserName:   userName,
	}

	//将评论存入MySQL中
	dbWithTransaction, err := mysql.InsertComment(&NewComment)
	if err != nil {
		return nil, err
	}

	//将评论存入Redis中
	key := tool.GetVideoCommentKey(videoId)
	err = redis.AddCommentToSortedSet(key, now.Unix(), &NewComment)
	if err != nil {
		dbWithTransaction.Rollback() //事务回滚
		return nil, err
	}

	dbWithTransaction.Commit() //提交事务

	return &NewComment, nil
}

// DeleteComment 删除评论
func DeleteComment(commentId int) error {

	//修改mysql中评论的状态
	dbWithTransaction, err := mysql.UpdateCommentStatus(commentId)
	if err != nil {
		return err
	}

	comment, err := mysql.GetComment(commentId)
	if err != nil {
		return err
	}

	//从redis中删除评论
	key := tool.GetVideoCommentKey(comment.VideoID)
	err = redis.RemoveComment(key, comment)
	if err != nil {
		dbWithTransaction.Rollback() //回滚事务
		return err
	}

	dbWithTransaction.Commit() //提交事务

	return nil
}

// ListComment 获取videoId的所有未被删除的评论
func ListComment(videoId int) (*[]model.Comment, error) {

	//commentList, err := mysql.ListCommentDESCByCreateDate(videoId)
	//if err != nil {
	//	return nil, err
	//}

	key := tool.GetVideoCommentKey(videoId)
	commentList, err := redis.ListComment(key)
	if err != nil {
		return nil, err
	}

	return commentList, err
}

// CountCommentByVideoId 获取视频未被删除的评论数
func CountCommentByVideoId(videoId int) (int64, error) {
	//count, err := mysql.CountCountByVideoId(videoId)

	key := tool.GetVideoCommentKey(videoId)
	count, err := redis.CountComment(key)
	if err != nil {
		return 0, err
	}

	return count, err
}
