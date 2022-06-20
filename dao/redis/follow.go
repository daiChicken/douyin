package redis

import "github.com/go-redis/redis"

/**
 * @author  daijizai
 * @date  2022/6/20 18:15
 * @version  1.0
 * @description
 */

func AddToSet(key string, userId int, txPipeline redis.Pipeliner) error {
	_, err := txPipeline.SAdd(key, userId).Result()
	if err != nil {
		return err
	}
	return nil
}

func RemoveFromSet(key string, userId int, txPipeline redis.Pipeliner) error {
	_, err := txPipeline.SRem(key, userId).Result()
	if err != nil {
		return err
	}
	return nil
}

func CountSet(key string) (int64, error) {
	result, err := rdb.SCard(key).Result()
	if err != nil {
		return 0, err
	}
	return result, err
}

func Smembers(key string) ([]string, error) {
	result, err := rdb.SMembers(key).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}
