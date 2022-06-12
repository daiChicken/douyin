package redis

import "github.com/go-redis/redis"

/**
 * @author  daijizai
 * @date  2022/6/12 22:34
 * @version  1.0
 * @description
 */

func AddUserToVideoSet(key string, userId int, txPipeline redis.Pipeliner) error {
	_, err := txPipeline.SAdd(key, userId).Result()
	if err != nil {
		return err
	}
	return nil
}

func RemoveUserFromVideoSet(key string, userId int, txPipeline redis.Pipeliner) error {
	_, err := txPipeline.SRem(key, userId).Result()
	if err != nil {
		return err
	}
	return nil
}

func CountLike(key string) (int64, error) {
	result, err := rdb.SCard(key).Result()
	if err != nil {
		return 0, err
	}
	return result, err
}

func GetLikeStatus(key string, userId int) (bool, error) {
	result, err := rdb.SIsMember(key, userId).Result()
	if err != nil {
		return false, err
	}
	return result, err
}

func AddVideoToUserSet(key string, VideoId int, txPipeline redis.Pipeliner) error {
	_, err := txPipeline.SAdd(key, VideoId).Result()
	if err != nil {
		return err
	}
	return nil
}

func RemoveVideoFromUserSet(key string, VideoId int, txPipeline redis.Pipeliner) error {
	_, err := txPipeline.SRem(key, VideoId).Result()
	if err != nil {
		return err
	}
	return nil
}

func ListLikedVideo(key string) ([]string, error) {
	result, err := rdb.SMembers(key).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetRDB() *redis.Client {
	return rdb
}
