package service

import (
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/model"
	"BytesDanceProject/pkg/snowflake"
	"BytesDanceProject/tool"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/viper"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

/**
 * @author  daijizai Congregalis
 * @date  2022/5/10 20:23
 * @version  1.0
 * @description
 */

// UploadVideo 上传视频
func UploadVideo(file *multipart.FileHeader, title string, authorId int) (err error) {

	//获取文件的后缀名
	filename := file.Filename                      //获取文件名
	indexOfDot := strings.LastIndex(filename, ".") //获取文件最后一个.的位置，这个.后的就是后缀名
	if indexOfDot < 0 {
		return errors.New("没有获取到文件的后缀名")
	}
	suffix := filename[indexOfDot+1:] //后缀名
	suffix = strings.ToLower(suffix)  //后缀名统一小写处理

	//判断文件是否符合视频格式
	if !tool.IsVideoExtension(suffix) {
		return errors.New("上传的文件不符合视频格式")
	}

	//生成新的文件名
	newFilename := strconv.FormatInt(snowflake.GenID(), 10) //使用雪花算法
	videoName := newFilename + "." + suffix                 //视频名
	coverName := newFilename + "." + "jpg"                  //封面名

	//上传视频和视频封面到七牛云（两个操作耦合）
	coverFolderName := "cover"                    //七牛云中存放图片的目录名。用于与文件名拼接，组成文件路径
	photoKey := coverFolderName + "/" + coverName //封面的访问路径，我们通过此路径在七牛云空间中定位封面
	entry := viper.GetString("qiniuyun.bucket") + ":" + photoKey
	encodedEntryURI := base64.StdEncoding.EncodeToString([]byte(entry))

	putPolicy := storage.PutPolicy{ //上传策略
		Scope: viper.GetString("qiniuyun.bucket"),
	}
	putPolicy.PersistentOps = "vframe/jpg/offset/1|saveas/" + encodedEntryURI //取视频第1秒的截图
	putPolicy.Expires = 7200                                                  //上传凭证的有效时间为2小时
	mac := qbox.NewMac(viper.GetString("qiniuyun.access_key"), viper.GetString("qiniuyun.secret_key"))
	upToken := putPolicy.UploadToken(mac) //上传凭证

	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}
	putExtra := storage.PutExtra{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	data, err := file.Open() //下文中的formUploader.Put()函数需要io.Reader类型的data
	if err != nil {
		return err
	}

	videoFolderName := "video"                    //七牛云中的目录名。用于与文件名拼接，组成文件路径
	videoKey := videoFolderName + "/" + videoName //文件访问路径，我们通过此路径在七牛云空间中定位文件

	//生成时间戳
	timeStamp := time.Now().UnixNano() / int64(time.Millisecond)

	//视频url
	//playUrl := "http://" + viper.GetString("qiniuyun.domain") + "/" + videoKey
	playUrl := videoKey

	//视频封面url
	//CoverUrl := "http://" + viper.GetString("qiniuyun.domain") + "/" + photoKey
	CoverUrl := photoKey

	newVideo := model.Video{
		AuthorId: authorId,
		PlayUrl:  playUrl,
		CoverUrl: CoverUrl,
		//FavoriteCount: 0,
		//CommentCount:  0,
		CreateTime: timeStamp,
		//IsDeleted: 0,
		Title: title,
	}

	//调用dao进行存储
	dbWithTransaction, err := mysql.InsertVideo(newVideo)
	if err != nil {
		return err
	}

	//起一个协程实现上传的异步
	go func() {
		//time.Sleep(time.Duration(5) * time.Second)//可供测试事务使用
		err := formUploader.Put(context.Background(), &ret, upToken,
			videoKey, data, file.Size, &putExtra)
		if err != nil {
			//问题：如果此处出现了问题导致上传失败，前端显示的也是上传成功。err信息没办法及时返回给controller
			fmt.Println("formUploader.Put()上传失败，错误信息：", err.Error())
			dbWithTransaction.Rollback() //事务回滚
			return
		}
		dbWithTransaction.Commit()            //提交事务
		fmt.Println("formUploader.Put()上传成功") //本行供测试使用
	}()
	//fmt.Println(ret.Key, ret.Hash)
	//到此上传视频到七牛云的工作完成

	return
}

// ListVideosByUser 获取用户所有投稿过的视频
func ListVideosByUser(authorId int) (*[]model.Video, error) {
	videoList, err := mysql.ListVideoByAuthorId(authorId)
	if err != nil {
		return nil, err
	}
	return videoList, err
}
