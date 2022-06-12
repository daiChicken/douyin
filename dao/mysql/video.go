package mysql

import (
	"BytesDanceProject/model"
	"gorm.io/gorm"
)

/**
 * @author  daijizai Congregalis
 * @date  2022/5/10 20:22
 * @version  1.0
 * @description
 */

// ListVideo 供feed流使用 获取视频列表  【！！！！应该限制时间】
func ListVideo(videoCount int) (*[]model.Video, error) {
	var videoList []model.Video
	err := db.Limit(videoCount).Find(&videoList).Error
	if err != nil {
		return nil, err
	}
	return &videoList, nil
}

// ListVideoDESCByCreateTime 根据创建时间倒序获取视频列表
// videoCount 限制返回的视频数量
// latestTime 限制返回视频的最新投稿时间
func ListVideoDESCByCreateTime(videoCount int, latestTime int64) (*[]model.Video, error) {
	var videoList []model.Video
	err := db.Where("create_time < ?", latestTime).Order("create_time desc").Limit(videoCount).Find(&videoList).Error
	if err != nil {
		return nil, err
	}
	return &videoList, nil
}

// InsertVideo 插入一条video记录 id为主键自增
func InsertVideo(v model.Video) (*gorm.DB, error) {
	dbWithTransaction := db.Begin() //开启事务
	if err := dbWithTransaction.Create(&v).Error; err != nil {
		return nil, err
	}
	return dbWithTransaction, nil
}

// ListVideoByAuthorId 根据作者id获取视频列表
func ListVideoByAuthorId(authorId int) (*[]model.Video, error) {
	var videoList []model.Video
	err := db.Where("author_id = ?", authorId).Order("create_time desc").Find(&videoList).Error
	if err != nil {
		return nil, err
	}
	return &videoList, nil
}

func GetVideoById(id int) (*model.Video, error) {
	var video model.Video
	err := db.Table("video").Where("id = ?", id).Find(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}
