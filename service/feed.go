package service

import (
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/model"
)

/**
 * @author  daijizai
 * @date  2022/5/10 20:23
 * @version  1.0
 * @description
 */

// GetVideoList 获取视频列表
func GetVideoList(videoCount int) (*[]model.Video, error) {

	videoList, err := mysql.ListVideo(videoCount)
	if err != nil {
		return nil, err
	}
	return videoList, err
}
