package redis

import (
	"strconv"
	"sync"
	"work/biz/dal/db"
	"work/pkg/errmsg"

	"github.com/go-redis/redis"
)

func PutVideoLikeInfo(vid string, uidList *[]string) error {
	pipe := redisDBVideoInfo.TxPipeline()
	pipe.Del(`l:` + vid)
	pipe.Del(`nl:` + vid)
	for _, item := range *uidList {
		pipe.SAdd(`l:`+vid, item)
	}
	if _, err := pipe.Exec(); err != nil {
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
	_, err := redisDBVideoInfo.ZAdd(`nl:`+vid, redis.Z{Score: 1, Member: uid}).Result()
	if err != nil {
		return err
	}
	if _, err := redisDBVideoInfo.SRem(`l:`+vid, uid).Result(); err != nil {
		return err
	}
	return nil
}

func AppendVideoLikeInfoToStaticSpace(vid, uid string) error {
	if _, err := redisDBVideoInfo.SAdd(`l:`+vid, uid).Result(); err != nil {
		return err
	}
	return nil
}

func DeleteVideoLikeInfoFromDynamicSpace(vid, uid string) error {
	if _, err := redisDBVideoInfo.ZRem(`nl:`+vid, uid).Result(); err != nil {
		return err
	}
	return nil
}

func RemoveVideoLikeInfo(vid, uid string) error {
	if !IsVideoExist(vid) {
		return errmsg.NoSuchVideoError
	}
	_, err := redisDBVideoInfo.ZAdd(`nl:`+vid, redis.Z{Score: 2, Member: uid}).Result()
	if err != nil {
		return err
	}
	if _, err := redisDBVideoInfo.SRem(`l:`+vid, uid).Result(); err != nil {
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
	exist, err := redisDBVideoInfo.SIsMember(`l:`+vid, uid).Result()
	if err != nil {
		return false, err
	}
	if !exist {
		_, err := redisDBVideoInfo.ZRank(`nl:`+vid, uid).Result()
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		return true, nil
	}
}

func GetVideoLikeList(vid string) (*[]string, error) {
	list, err := redisDBVideoInfo.SMembers(`l:` + vid).Result()
	if err != nil {
		return nil, err
	}
	nList, err := GetNewUpdateVideoLikeList(vid)
	if err != nil {
		return nil, err
	}
	list = append(list, *nList...)
	return &list, nil
}

func GetNewUpdateVideoLikeList(vid string) (*[]string, error) {
	list, err := redisDBVideoInfo.ZRangeByScore(`nl:`+vid, redis.ZRangeBy{Min: `1`, Max: `1`}).Result()
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func GetNewDeleteVideoLikeList(vid string) (*[]string, error) {
	list, err := redisDBVideoInfo.ZRangeByScore(`nl:`+vid, redis.ZRangeBy{Min: `2`, Max: `2`}).Result()
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func GetVideoLikeCount(vid string) (int64, error) {
	var count int64
	var err error
	if count, err = redisDBVideoInfo.SCard(`l:` + vid).Result(); err != nil {
		return -1, err
	}
	if nCount, err := redisDBVideoInfo.ZCount(`nl:`+vid, `1`, `1`).Result(); err != nil {
		return -1, err
	} else {
		return count + nCount, nil
	}
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
