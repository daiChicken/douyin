package redis

import (
	"BytesDanceProject/model"
	"errors"
	"github.com/go-redis/redis"
	"math"
	"strconv"
)

/**
 * @author  daijizai
 * @date  2022/6/2 23:10
 * @version  1.0
 * @description
 */

func AddCommentToSortedSet(key string, score int64, comment *model.Comment) error {
	z := redis.Z{
		Score:  float64(score),
		Member: comment,
	}

	_, err := rdb.ZAdd(key, z).Result()
	if err != nil {
		return err
	}
	return nil
}

// ListComment 根据分数从高到低获取所有的评论
func ListComment(key string) (*[]model.Comment, error) {

	var commentList []model.Comment
	err := rdb.ZRevRange(key, 0, -1).ScanSlice(&commentList)
	if err != nil {
		return nil, err
	}

	return &commentList, nil
}

func RemoveComment(key string, comment *model.Comment) error {
	count, err := rdb.ZRem(key, comment).Result()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("redis中删除评论失败")
	}
	return nil
}

func CountComment(key string) (int64, error) {
	count, err := rdb.ZCount(key, strconv.Itoa(0), strconv.Itoa(math.MaxInt)).Result()
	if err != nil {
		return 0, err
	}
	return count, nil
}
