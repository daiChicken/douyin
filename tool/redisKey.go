package tool

import "strconv"

/**
 * @author  daijizai
 * @date  2022/6/2 23:00
 * @version  1.0
 * @description
 */

const split = ":"

func GetVideoCommentKey(videoId int) string {
	return "video" + split + "comment" + split + strconv.Itoa(videoId)
}

func GetVideoLikeKey(videoId int) string {
	return "video" + split + "like" + split + strconv.Itoa(videoId)
}

func GetUserLikeKey(userId int) string {
	return "user" + split + "like" + split + strconv.Itoa(userId)
}
