/**
 @author: 15973
 @date: 2022/07/11
 @note:
**/
package initialization

import (
	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
)

func InitializeRedisClient() *redis.Client {
	redisCli := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Password: "",
		DB:       0,
	})
	return redisCli
}
