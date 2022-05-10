package service

import (
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/model"
	"BytesDanceProject/pkg/snowflake"
	"BytesDanceProject/tool"
	"context"
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

func UploadVideo(file *multipart.FileHeader) (err error) {

	//获取文件的后缀名
	filename := file.Filename //获取文件名
	//strings.LastIndex(s string, str string) int：判断str在s中最后出现的位置，如果没有出现，则返回-1
	indexOfDot := strings.LastIndex(filename, ".")
	if indexOfDot < 0 {
		//没有获取到文件的后缀名
		return errors.New("没有获取到文件的后缀名")
	}
	suffix := filename[indexOfDot+1 : len(filename)] //后缀名
	//strings.ToLower(str string)string：转为⼩写
	suffix = strings.ToLower(suffix) //后缀名统一小写处理

	//判断文件是否符合视频格式
	if !tool.IsVideoAllowed(suffix) {
		//上传的文件不符合视频格式
		return errors.New("上传的文件不符合视频格式")
	}

	//使用雪花算法生成新的文件名
	filename = strconv.FormatInt(snowflake.GenID(), 10)

	filename = filename + "." + suffix

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

	data, err := file.Open() //！！！！！！！！！！！！！！！！！！！！！
	if err != nil {
		return err
	}

	folderName := "video" //！！！！！！！！！！！！！！！！！
	key := folderName + "/" + filename

	err = formUploader.Put(context.Background(), &ret, upToken,
		key, data, file.Size, &putExtra)
	if err != nil {
		return err
	}
	//fmt.Println(ret.Key, ret.Hash)

	//生成时间戳
	timeStamp := time.Now().Unix()

	//获取当前登录用户的id
	authorId := 0 //！！！！！！！！！！！！！！！！
	newVideo := model.Video{
		AuthorId: authorId,
		PlayUrl:  viper.GetString("qiniuyun.domain") + "/" + key,
		CoverUrl: "",
		//FavoriteCount: 0,
		//CommentCount:  0,
		CreateTime: timeStamp,
		//IsDeleted: 0,
	}

	fmt.Println("newVideo.PlayUrl", newVideo.PlayUrl)

	//调用dao进行存储
	err = mysql.InsertVideo(newVideo)
	if err != nil {
		return err
	}

	return
}
