package redis

import (
	"strconv"
	"sync"
	"work/pkg/errmsg"

	"github.com/go-redis/redis"
)

func PutVideoLikeInfo(vid string, uidList *[]string) error {
	exist, err := redisDBVideoInfo.Exists(`l:` + vid).Result()
	if err != nil {
		return errmsg.RedisError
	}
	pipe := redisDBVideoInfo.TxPipeline()
	if exist != 0 {
		pipe.Del(`l:` + vid)
	}
	for _, item := range *uidList {
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

func PutVideoCommentInfo(vid string, cidList *[]string) error {
	exist, err := redisDBVideoInfo.Exists(`c:` + vid).Result()
	if err != nil {
		return errmsg.RedisError
	}
	pipe := redisDBVideoInfo.TxPipeline()
	if exist != 0 {
		pipe.Del(`c:` + vid)
	}
	for _, item := range *cidList {
		pipe.RPush(`c:`+vid, item)
	}
	if _, err = pipe.Exec(); err != nil {
		return err
	}
	return nil
}

func AppendVideoLikeInfo(vid, uid string) error {
	if !IsVideoExist(vid) {
		return errmsg.NoSuchVideoError
	}
	_, err := redisDBVideoInfo.HSet(`l:`+vid, uid, 1).Result()
	if err != nil {
		return errmsg.RedisError
	}
	return nil
}

func RemoveVideoLikeInfo(vid, uid string) error {
	if !IsVideoExist(vid) {
		return errmsg.NoSuchVideoError
	}
	lExist, err := redisDBVideoInfo.HExists(`l:`+vid, uid).Result()
	if err != nil {
		return errmsg.RedisError
	}
	if !lExist {
		return nil
	}
	_, err = redisDBVideoInfo.HSet(`l:`+vid, uid, 2).Result()
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

func AppendVideoCommentInfo(vid, cid string) error {
	_, err := redisDBVideoInfo.RPush(`c:`+vid, cid).Result()
	if err != nil {
		return errmsg.RedisError
	}
	return nil
}

func RemoveVideoCommentInfo(vid, cid string) error {
	_, err := redisDBCommentInfo.LRem(`c:`+vid, 1, cid).Result()
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

func GetVideoLikeList(vid string) (*map[string]string, error) {
	list, err := redisDBVideoInfo.HGetAll(`l:` + vid).Result()
	if err != nil {
		return nil, errmsg.RedisError
	}
	return &list, nil
}

func GetVideoCommentCount(vid string) (int64, error) {
	count, err := redisDBVideoInfo.LLen(`c:` + vid).Result()
	if err != nil {
		return -1, errmsg.RedisError
	}
	return count, nil
}

func GetVideoCommentList(vid string, pageNum, pageSize int64) (*[]string, error) {
	list, err := redisDBVideoInfo.LRange(`c:`+vid, (pageNum-1)*pageSize, pageNum*pageSize-1).Result()
	if err != nil {
		return nil, errmsg.RedisError
	}
	return &list, nil
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
	list, err := redisDBVideoInfo.ZRevRange(`visit`, (pageNum-1)*pageSize, pageNum*pageSize-1).Result()
	if err != nil {
		return nil, err
	}
	return &list, err
}

func IsVideoExist(vid string) bool {
	_, err := redisDBVideoInfo.ZScore(`visit`, vid).Result()
	return err == nil
}

func DeleteVideoAndAllAbout(vid string) error {
	videoPipe := redisDBVideoInfo.TxPipeline()
	commentPipe := redisDBCommentInfo.TxPipeline()

	commentList, err := redisDBVideoInfo.LRange(`c:`+vid, 0, -1).Result()
	if err != nil {
		return err
	}

	videoPipe.Del(`c:`+vid, `l:`+vid)
	videoPipe.ZRem(`visit`, vid)

	for _, item := range commentList {
		commentPipe.Del(`i:`+item, `c:`+item, `l:`+item)
	}

	var (
		wg      sync.WaitGroup
		errChan = make(chan error, 2)
	)
	wg.Add(2)
	go func() {
		if _, err := videoPipe.Exec(); err != nil {
			errChan <- err
		}
		wg.Done()
	}()
	go func() {
		if _, err := commentPipe.Exec(); err != nil {
			errChan <- err
		}
		wg.Done()
	}()
	wg.Wait()
	select {
	case result := <-errChan:
		return result
	default:
	}
	return nil
}
