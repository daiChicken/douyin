package service

import (
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/model"
	"BytesDanceProject/pkg/snowflake"
	"BytesDanceProject/tool"
	"context"
	"errors"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/viper"
)

/**
 * @author  daijizai Congregalis
 * @date  2022/5/10 20:23
 * @version  1.0
 * @description
 */

// UploadVideo 上传视频
func UploadVideo(file *multipart.FileHeader) (err error) {

	//获取文件的后缀名
	filename := file.Filename                      //获取文件名
	indexOfDot := strings.LastIndex(filename, ".") //获取文件最后一个.的位置，这个.后的就是后缀名
	if indexOfDot < 0 {
		return errors.New("没有获取到文件的后缀名")
	}
	suffix := filename[indexOfDot+1 : len(filename)] //后缀名
	suffix = strings.ToLower(suffix)                 //后缀名统一小写处理

	//判断文件是否符合视频格式
	if !tool.IsVideoExtension(suffix) {
		return errors.New("上传的文件不符合视频格式")
	}

	//生成新的文件名
	newFilename := strconv.FormatInt(snowflake.GenID(), 10) //使用雪花算法
	newFilename = newFilename + "." + suffix                //给新的文件名加上后缀名

	//上传视频到七牛云
	//自定义凭证有效期（示例2小时，Expires 单位为秒，为上传凭证的有效时间）
	putPolicy := storage.PutPolicy{
		Scope: viper.GetString("qiniuyun.bucket"),
	}
	putPolicy.Expires = 7200 //示例2小时有效期
	mac := qbox.NewMac(viper.GetString("qiniuyun.access_key"), viper.GetString("qiniuyun.secret_key"))
	upToken := putPolicy.UploadToken(mac)

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

	folderName := "video"                 //七牛云中的目录名。用于与文件名拼接，组成文件路径
	key := folderName + "/" + newFilename //文件访问路径，我们通过此路径在七牛云中定位文件

	err = formUploader.Put(context.Background(), &ret, upToken,
		key, data, file.Size, &putExtra)
	if err != nil {
		return err
	}
	//fmt.Println(ret.Key, ret.Hash)
	//到此上传视频到七牛云的工作完成

	//生成时间戳
	timeStamp := time.Now().Unix()

	authorId := 0 //此处应该获取当前登录用户的id！！！！！！！！！！
	newVideo := model.Video{
		AuthorId: authorId,
		PlayUrl:  viper.GetString("qiniuyun.domain") + "/" + key,
		CoverUrl: "",
		//FavoriteCount: 0,
		//CommentCount:  0,
		CreateTime: timeStamp,
		//IsDeleted: 0,
	}

	//调用dao进行存储
	err = mysql.InsertVideo(newVideo)
	if err != nil {
		return err
	}

	return
}

// ListVideosByUser 获取用户所有投稿过的视频
func ListVideosByUser() (*[]model.Video, error) { //	【！！！！！此处应该传入当前登录用户的对象，因为还没有创建user对象，故不进行此操作】

	//通过函数的参数获取user对象
	//根据user对象获取获取userid
	userId := 0 //！！！！！！！！！！！！！假数据
	videoList, err := mysql.ListVideoByAuthorId(userId)
	if err != nil {
		return nil, err
	}
	return videoList, err
}
