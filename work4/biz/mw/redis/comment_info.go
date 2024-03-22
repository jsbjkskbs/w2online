package redis

import (
	"work/biz/dal/db"

	"github.com/go-redis/redis"
)

func PutCommentLikeInfo(cid string, uidList *[]string) error {
	pipe := redisDBCommentInfo.TxPipeline()
	pipe.Del(`l:` + cid)
	pipe.Del(`nl:` + cid)
	for _, item := range *uidList {
		pipe.SAdd(`l:`+cid, item)
	}
	if _, err := pipe.Exec(); err != nil {
		return err
	}
	return nil
}

func GetCommentLikeCount(cid string) (int64, error) {
	countOld, err := redisDBCommentInfo.SCard(`l:` + cid).Result()
	if err != nil {
		return -1, err
	}
	countNew, err := redisDBCommentInfo.ZCount(`nl:`+cid, `1`, `1`).Result()
	if err != nil {
		return -1, err
	}
	return countOld + countNew, nil
}

func GetCommentLikeList(cid string) (*[]string, error) {
	list, err := redisDBCommentInfo.SMembers(`l:` + cid).Result()
	if err != nil {
		return nil, err
	}
	nList, err := GetNewUpdateCommentLikeList(cid)
	if err != nil {
		return nil, err
	}
	list = append(list, *nList...)
	return &list, nil
}

func GetNewUpdateCommentLikeList(cid string) (*[]string, error) {
	list, err := redisDBCommentInfo.ZRangeByScore(`nl:`+cid, redis.ZRangeBy{Min: `1`, Max: `1`}).Result()
	if err != nil {
		return nil, err
	}
	return &list, err
}

func GetNewDeleteCommentLikeList(cid string) (*[]string, error) {
	list, err := redisDBCommentInfo.ZRangeByScore(`nl:`+cid, redis.ZRangeBy{Min: `2`, Max: `2`}).Result()
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func AppendCommentLikeInfo(cid, uid string) error {
	_, err := redisDBCommentInfo.ZAdd(`nl:`+cid, redis.Z{Score: 1, Member: uid}).Result()
	if err != nil {
		return err
	}
	if _, err := redisDBCommentInfo.SRem(`l:`+cid, uid).Result(); err != nil {
		return err
	}
	return nil
}

func AppendCommentLikeInfoToStaticSpace(cid, uid string) error {
	if _, err := redisDBCommentInfo.SAdd(`l:`+cid, uid).Result(); err != nil {
		return err
	}
	return nil
}

func DeleteCommentLikeInfoFromDynamicSpace(cid, uid string) error {
	if _, err := redisDBCommentInfo.ZRem(`nl:`+cid, uid).Result(); err != nil {
		return err
	}
	return nil
}

func RemoveCommentLikeInfo(cid, uid string) error {
	_, err := redisDBCommentInfo.ZAdd(`nl:`+cid, redis.Z{Score: 2, Member: uid}).Result()
	if err != nil {
		return err
	}
	if _, err := redisDBCommentInfo.SRem(`l:`+cid, uid).Result(); err != nil {
		return err
	}
	return nil
}

func DeleteCommentAndAllAbout(cid string) error {
	commentPipe := redisDBCommentInfo.TxPipeline()

	var (
		childList *[]string
		err       error
	)
	if childList, err = db.GetCommentChildList(cid); err != nil {
		return err
	}

	commentPipe.Del(`l:` + cid)
	commentPipe.Del(`nl:` + cid)
	for _, item := range *childList {
		commentPipe.Del(`l:` + item)
		commentPipe.Del(`nl:` + item)
	}

	if _, err := commentPipe.Exec(); err != nil {
		return err
	}
	return nil
}
