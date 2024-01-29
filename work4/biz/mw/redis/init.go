package redis

import (
	"work/pkg/constants"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redis/redis"
)

var (
	redisDBAvatarUpload *redis.Client
	redisDBVideoUpload  *redis.Client
	redisDBVideoInfo    *redis.Client
)

func Init() {

	redisDBAvatarUpload = redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: constants.RedisPassword,
		DB:       0,
	})

	redisDBVideoUpload = redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: constants.RedisPassword,
		DB:       1,
	})

	redisDBVideoInfo = redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: constants.RedisPassword,
		DB:       2,
	})

	if _, err := redisDBAvatarUpload.Ping().Result(); err != nil {
		panic(err)
	}
	if _, err := redisDBVideoUpload.Ping().Result(); err != nil {
		panic(err)
	}
	if _, err := redisDBVideoInfo.Ping().Result(); err != nil {
		panic(err)
	}
	hlog.Info("Redis connected successfully.")
}
