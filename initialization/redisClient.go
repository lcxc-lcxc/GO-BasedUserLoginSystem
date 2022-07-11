/**
 @author: 15973
 @date: 2022/07/11
 @note:
**/
package initialization

import "github.com/go-redis/redis/v9"

func InitializeRedisClient() *redis.Client {
	redisCli := redis.NewClient(&redis.Options{
		Addr:     "192.168.56.108:6379",
		Password: "",
		DB:       0,
	})
	return redisCli
}
