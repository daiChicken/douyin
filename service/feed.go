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
func ListVideos(videoCount int, latestTime int64) (*[]model.Video, int64, error) {

	videoList, err := mysql.ListVideoDESCByCreateTime(videoCount, latestTime)
	if err != nil {
		return nil, 0, err
	}

	if len(*videoList) == 0 {
		return videoList, latestTime, nil
	}
	nextTime := (*videoList)[len(*videoList)-1].CreateTime
	return videoList, nextTime, err
}
