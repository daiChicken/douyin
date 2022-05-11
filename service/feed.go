package service

import (
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/model"
)

/**
 * @author  daijizai Congregalis
 * @date  2022/5/10 20:23
 * @version  1.0
 * @description
 */

// ListVideos 获取视频列表
func ListVideos(videoCount int) (*[]model.Video, error) {

	videoList, err := mysql.ListVideo(videoCount)
	if err != nil {
		return nil, err
	}
	return videoList, err
}
