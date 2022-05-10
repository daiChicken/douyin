package tool

/**
 * @author  daijizai
 * @date  2022/5/10 20:23
 * @version  1.0
 * @description
 */

var videoFileExt = []string{"mp4", "flv"}
var imageFileExt = []string{"png", "bmp", "jpg", "jpeg"}

// IsVideoAllowed 判断提供的后缀是否符合视频的格式
func IsVideoAllowed(suffix string) bool {
	for _, fileExt := range videoFileExt {
		if suffix == fileExt {
			return true
		}
	}
	return false
}

// IsImageAllowed 判断提供的后缀是否符合图片的格式
func IsImageAllowed(suffix string) bool {
	for _, fileExt := range imageFileExt {
		if suffix == fileExt {
			return true
		}
	}
	return false
}
