package redis

import (
	"strconv"
	"sync"
	"work/biz/dal/db"
	"work/pkg/errmsg"

	"github.com/go-redis/redis"
)

func PutVideoLikeInfo(vid string, uidList *[]string) error {
	exist, err := redisDBVideoInfo.Exists(`l:` + vid).Result()
	if err != nil {
		return err
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
		return err
	}
	return nil
}

func RemoveVideoLikeInfo(vid, uid string) error {
	if !IsVideoExist(vid) {
		return errmsg.NoSuchVideoError
	}
	lExist, err := redisDBVideoInfo.HExists(`l:`+vid, uid).Result()
	if err != nil {
		return err
	}
	if !lExist {
		return nil
	}
	_, err = redisDBVideoInfo.HSet(`l:`+vid, uid, 2).Result()
	if err != nil {
		return err
	}
	return nil
}

func IncrVideoVisitInfo(vid string) error {
	_, err := redisDBVideoInfo.ZIncrBy(`visit`, 1, vid).Result()
	if err != nil {
		return err
	}
	return nil
}

func IsVideoLikedByUser(vid, uid string) (bool, error) {
	exist, err := redisDBVideoInfo.HExists(`l:`+vid, uid).Result()
	if err != nil {
		return true, err
	}
	return exist, nil
}

func GetVideoLikeList(vid string) (*map[string]string, error) {
	list, err := redisDBVideoInfo.HGetAll(`l:` + vid).Result()
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func GetVideoVisitCount(vid string) (int64, error) {
	_, err := redisDBVideoInfo.ZRank(`visit`, vid).Result()
	if err != nil {
		return -1, err
	}
	s, err := redisDBVideoInfo.ZScore(`visit`, vid).Result()
	if err != nil {
		return -1, err
	}
	return int64(s), nil
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

	commentList, err := db.GetVideoCommentList(vid)
	if err != nil {
		return err
	}

	videoPipe.Del(`l:` + vid)
	videoPipe.ZRem(`visit`, vid)

	for _, item := range *commentList {
		commentPipe.Del(`l:` + item)
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
