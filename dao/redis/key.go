package redis

// 和redis相关的key

const (
	KeyPrefix = 					"douyin:"
	KeyFollowPrefix = 				"follow:"
	KeyFollowerPrefix = 			"follower:"
	KeyInterStorePrefix = 			"flwandflwer:" //set 取follow和follwer的交集

)

// getRedisKey 给Redis Key 加前缀
func getRedisKey(key string)string{
	return KeyPrefix + key
}
