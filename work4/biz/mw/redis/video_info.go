package redis

import (
	"strconv"
	"work/pkg/errmsg"

	"github.com/go-redis/redis"
)

// 浏览量 点赞关系 评论数 (MySQL->Redis,Elasticsearch(初始化) Redis->MySQL,Elasticsearch(运行时))

func PutVideoLikeInfo(vid string, uid []string) error {
	exist, err := redisDBVideoInfo.Exists(`l:` + vid).Result()
	if err != nil {
		return errmsg.RedisError
	}
	pipe := redisDBVideoInfo.TxPipeline()
	if exist != 0 {
		pipe.Del(`l:` + vid)
	}
	for _, item := range uid {
		pipe.HSet(`l:`+vid, item, 1)
	}
	if _, err = pipe.Exec(); err != nil {
		return err
	}
	return nil
}

func PutVideoVisitInfo(vid, visitCount string) error {
	score, _ := strconv.ParseFloat(visitCount, 64)
	_, err := redisDBVideoInfo.ZAdd(`visit`, redis.Z{Score: score, Member: vid}).Result()
	if err != nil {
		return errmsg.RedisError
	}
	return nil
}

func PutVideoCommentInfo(vid, commentCount string) error {
	_, err := redisDBVideoInfo.Set(`c:`+vid, commentCount, 0).Result()
	if err != nil {
		return errmsg.RedisError
	}
	return nil
}

func AppendVideoLikeInfo(vid, uid string) error {
	_, err := redisDBVideoInfo.HSet(`l:`+vid, uid, 1).Result()
	if err != nil {
		return errmsg.RedisError
	}
	return nil
}

func RemoveVideoLikeInfo(vid, uid string) error {
	exist, err := redisDBVideoInfo.HExists(`l:`+vid, uid).Result()
	if err != nil {
		return errmsg.RedisError
	}
	if !exist {
		return nil
	}
	_, err = redisDBVideoInfo.HDel(`l:`+vid, uid).Result()
	if err != nil {
		return errmsg.RedisError
	}
	return nil
}

func IncrVideoVisitInfo(vid string) error {
	_, err := redisDBVideoInfo.ZIncrBy(`visit`, 1, vid).Result()
	if err != nil {
		return errmsg.RedisError
	}
	return nil
}

func IncrVideoCommentInfo(vid string) error {
	_, err := redisDBVideoInfo.Incr(`c:` + vid).Result()
	if err != nil {
		return errmsg.RedisError
	}
	return nil
}

func DecrVideoCommentInfo(vid string) error {
	_, err := redisDBVideoInfo.Decr(`c:` + vid).Result()
	if err != nil {
		return errmsg.RedisError
	}
	return nil
}

func IsVideoLikedByUser(vid, uid string) (bool, error) {
	exist, err := redisDBVideoInfo.HExists(`l:`+vid, uid).Result()
	if err != nil {
		return true, errmsg.RedisError
	}
	return exist, nil
}

func GetVideoLikeList(vid string) (*[]string, error) {
	list, err := redisDBVideoInfo.HKeys(`l:` + vid).Result()
	if err != nil {
		return nil, errmsg.RedisError
	}
	return &list, nil
}

func GetVideoCommentCount(vid string) (int64, error) {
	s, err := redisDBVideoInfo.Get(`c:` + vid).Result()
	if err != nil {
		return -1, errmsg.RedisError
	}
	count, _ := strconv.ParseInt(s, 10, 64)
	return count, nil
}

func GetVideoVisitCount(vid string) (int64, error) {
	if exist, err := redisDBVideoInfo.HExists(`visit`, vid).Result(); err != nil {
		return -1, errmsg.RedisError
	} else if !exist {
		return -1, nil
	}
	s, err := redisDBVideoInfo.HGet(`visit`, vid).Result()
	if err != nil {
		return -1, errmsg.RedisError
	}
	count, _ := strconv.ParseInt(s, 10, 64)
	return count, nil
}

func GetVideoPopularList(pageNum, pageSize int64) (*[]string, error) {
	list, err := redisDBVideoInfo.ZRevRange(`visit`, (pageNum-1)*pageSize, pageNum*pageSize).Result()
	if err != nil {
		return nil, err
	}
	return &list, err
}
