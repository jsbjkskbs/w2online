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
	redisDBCommentInfo  *redis.Client
	redisDBRelationInfo   *redis.Client
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

	redisDBCommentInfo = redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: constants.RedisPassword,
		DB:       3,
	})

	redisDBRelationInfo = redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: constants.RedisPassword,
		DB:       4,
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
	if _, err := redisDBCommentInfo.Ping().Result(); err != nil {
		panic(err)
	}
	if _, err := redisDBRelationInfo.Ping().Result(); err != nil {
		panic(err)
	}
	hlog.Info("Redis connected successfully.")
}
