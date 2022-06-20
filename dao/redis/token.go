package redis

/**
 * @author  daijizai
 * @date  2022/6/20 10:41
 * @version  1.0
 * @description
 */

func InsertToken(key string, value string) error {
	_, err := rdb.Set(key, value, 0).Result()
	if err != nil {
		return err
	}
	return nil
}
