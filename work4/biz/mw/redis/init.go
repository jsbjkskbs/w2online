package redis

import (
	"work/pkg/constants"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redis/redis"
)

var (
	redisDBVideoUpload  *redis.Client
	redisDBVideoInfo    *redis.Client
	redisDBCommentInfo  *redis.Client
	redisDBChatInfo     *redis.Client
)

func Load() {

	redisDBVideoUpload = redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: constants.RedisPassword,
		DB:       0,
	})

	redisDBVideoInfo = redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: constants.RedisPassword,
		DB:       1,
	})

	redisDBCommentInfo = redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: constants.RedisPassword,
		DB:       2,
	})

	redisDBChatInfo = redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: constants.RedisPassword,
		DB:       3,
	})

	if _, err := redisDBVideoUpload.Ping().Result(); err != nil {
		panic(err)
	}
	if _, err := redisDBVideoInfo.Ping().Result(); err != nil {
		panic(err)
	}
	if _, err := redisDBCommentInfo.Ping().Result(); err != nil {
		panic(err)
	}
	if _, err := redisDBChatInfo.Ping().Result(); err != nil {
		panic(err)
	}
	hlog.Info("Redis connected successfully.")
}
