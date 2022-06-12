package redis

/**
 * @author  daijizai
 * @date  2022/6/12 22:34
 * @version  1.0
 * @description
 */

func AddUserIdToSet(key string, userId int) error {
	_, err := rdb.SAdd(key, userId).Result()
	if err != nil {
		return err
	}
	return nil
}

func RemoveUserIdFromSet(key string, userId int) error {
	_, err := rdb.SRem(key, userId).Result()
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
