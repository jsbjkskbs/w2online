package redis

import (
	"work/biz/dal/db"
)

func PutCommentLikeInfo(cid string, uidList *[]string) error {
	exist, err := redisDBCommentInfo.Exists(`l:` + cid).Result()
	if err != nil {
		return err
	}
	pipe := redisDBCommentInfo.TxPipeline()
	if exist != 0 {
		pipe.Del(`l:` + cid)
	}
	for _, item := range *uidList {
		pipe.HSet(`l:`+cid, item, 1)
	}
	if _, err := pipe.Exec(); err != nil {
		return err
	}
	return nil
}

func GetCommentLikeCount(cid string) (int64, error) {
	count, err := redisDBCommentInfo.HLen(`l:` + cid).Result()
	if err != nil {
		return -1, err
	}
	return count, nil
}

func GetCommentLikeList(cid string) (*map[string]string, error) {
	list, err := redisDBCommentInfo.HGetAll(`l:` + cid).Result()
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func AppendCommentLikeInfo(cid, uid string) error {
	_, err := redisDBCommentInfo.HSet(`l:`+cid, uid, 1).Result()
	if err != nil {
		return err
	}

	go func(_cid string) {
		if exist, _ := db.IsCommentExist(_cid); !exist {
			redisDBCommentInfo.Del(`l:` + _cid)
		}
	}(cid)

	return nil
}

func RemoveCommentLikeInfo(cid, uid string) error {
	lExist, err := redisDBCommentInfo.HExists(`l:`+cid, uid).Result()
	if err != nil {
		return err
	}
	if !lExist {
		return nil
	}
	_, err = redisDBCommentInfo.HSet(`l:`+cid, uid, 2).Result()
	if err != nil {
		return err
	}

	go func(_cid string) {
		if exist, _ := db.IsCommentExist(_cid); !exist {
			redisDBCommentInfo.Del(`l:` + _cid)
		}
	}(cid)

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
	for _, item := range *childList {
		commentPipe.Del(`l:` + item)
	}

	if _, err := commentPipe.Exec(); err != nil {
		return err
	}
	return nil
}
