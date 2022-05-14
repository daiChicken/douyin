package tool

/**
 * @author  daijizai Congregalis
 * @date  2022/5/10 20:23
 * @version  1.0
 * @description 用于存放工具函数
 */

var videoExtensions = []string{"mp4", "flv"}                //符合要求的视频扩展名
var imageExtensions = []string{"png", "bmp", "jpg", "jpeg"} //符合要求的图片扩展名

// IsVideoExtension 判断提供的字符串是否是符合要求的视频扩展名
func IsVideoExtension(suffix string) bool {
	for _, videoExt := range videoExtensions {
		if suffix == videoExt {
			return true
		}
	}
	return false
}

// IsImageExtension 判断提供的字符串是否是符合要求的图片扩展名
func IsImageExtension(suffix string) bool {
	for _, imageExt := range imageExtensions {
		if suffix == imageExt {
			return true
		}
	}
	return false
}
