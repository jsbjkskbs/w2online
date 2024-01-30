package redis

import (
	"fmt"
	"strconv"
	"sync"
	"work/biz/dal/db"
	"work/pkg/errmsg"
)

func PutCommentLikeInfo(cid string, uidList *[]string) error {
	exist, err := redisDBCommentInfo.Exists(`l:` + cid).Result()
	if err != nil {
		return errmsg.RedisError
	}
	pipe := redisDBCommentInfo.TxPipeline()
	if exist != 0 {
		pipe.Del(`l:` + cid)
	}
	for _, item := range *uidList {
		pipe.HSet(`l:`+cid, item, 1)
	}
	if _, err := pipe.Exec(); err != nil {
		return errmsg.RedisError
	}
	return nil
}

func PutCommentChildInfo(cid string, cidList *[]string) error {
	exist, err := redisDBCommentInfo.Exists(`c:` + cid).Result()
	if err != nil {
		return errmsg.RedisError
	}
	pipe := redisDBCommentInfo.TxPipeline()
	if exist != 0 {
		pipe.Del(`c:` + cid)
	}
	for _, item := range *cidList {
		pipe.RPush(`c:`+cid, item).Result()
	}
	if _, err := pipe.Exec(); err != nil {
		return errmsg.RedisError
	}
	return nil
}

func PutCommentInfo(comment *db.Comment) error {
	exist, err := redisDBCommentInfo.Exists(`i:` + fmt.Sprint(comment.Id)).Result()
	if err != nil {
		return err
	}
	pipe := redisDBCommentInfo.TxPipeline()
	if exist != 0 {
		pipe.Del(`i:` + fmt.Sprint(comment.Id))
	}
	pipe.RPush(`i:`+fmt.Sprint(comment.Id), comment.UserId, comment.VideoId, comment.ParentId, comment.Content, comment.CreatedAt, comment.UpdatedAt, comment.DeletedAt)
	if _, err := pipe.Exec(); err != nil {
		return errmsg.RedisError
	}
	return nil
}

func GetCommentLikeCount(cid string) (int64, error) {
	count, err := redisDBCommentInfo.HLen(`l:` + cid).Result()
	if err != nil {
		return -1, errmsg.RedisError
	}
	return count, nil
}

func GetCommentLikeList(cid string) (*map[string]string, error) {
	list, err := redisDBCommentInfo.HGetAll(`l:` + cid).Result()
	if err != nil {
		return nil, errmsg.RedisError
	}
	return &list, nil
}

func GetCommentChildCount(cid string) (int64, error) {
	count, err := redisDBCommentInfo.LLen(`c:` + cid).Result()
	if err != nil {
		return -1, errmsg.RedisError
	}
	return count, err
}

func GetCommentChildList(cid string, pageNum, pageSize int64) (*[]string, error) {
	list, err := redisDBCommentInfo.LRange(`c:`+cid, (pageNum-1)*pageSize, pageNum*pageSize-1).Result()
	if err != nil {
		return nil, errmsg.RedisError
	}
	return &list, nil
}

func GetCommentInfo(cid string) (*db.Comment, error) {
	info, err := redisDBCommentInfo.LRange(`i:`+cid, 0, 6).Result()
	if err != nil {
		return nil, errmsg.RedisError
	}
	createdAt, _ := strconv.ParseInt(info[4], 10, 64)
	updatedAt, _ := strconv.ParseInt(info[5], 10, 64)
	deletedAt, _ := strconv.ParseInt(info[6], 10, 64)
	return &db.Comment{
		UserId:    info[0],
		VideoId:   info[1],
		ParentId:  info[2],
		Content:   info[3],
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}, nil
}

func GetCommentVideoId(cid string) (string, error) {
	info, err := redisDBCommentInfo.LRange(`i:`+cid, 1, 1).Result()
	if err != nil {
		return ``, errmsg.RedisError
	}
	return info[0], nil
}

func CreateCommentInfo(comment *db.Comment) error {
	return PutCommentInfo(comment)
}

func AppendCommentLikeInfo(cid, uid string) error {
	if !IsCommentExist(cid) {
		return errmsg.NoSuchCommentError
	}
	_, err := redisDBCommentInfo.HSet(`l:`+cid, uid, 1).Result()
	if err != nil {
		return errmsg.RedisError
	}
	return nil
}

func RemoveCommentLikeInfo(cid, uid string) error {
	if !IsCommentExist(cid) {
		return errmsg.NoSuchCommentError
	}
	lExist, err := redisDBCommentInfo.HExists(`l:`+cid, uid).Result()
	if err != nil {
		return errmsg.RedisError
	}
	if !lExist {
		return nil
	}
	_, err = redisDBCommentInfo.HSet(`l:`+cid, uid, 2).Result()
	if err != nil {
		return errmsg.RedisError
	}
	return nil
}

func AppendChildCommentInfo(cid, childId string) error {
	_, err := redisDBCommentInfo.RPush(`c:`+cid, childId).Result()
	if err != nil {
		return errmsg.RedisError
	}
	return nil
}

func RemoveChildCommentInfo(cid, childId string) error {
	_, err := redisDBCommentInfo.LRem(`c:`+cid, 1, childId).Result()
	if err != nil {
		return errmsg.RedisError
	}
	return nil
}

func IsCommentExist(cid string) bool {
	exist, _ := redisDBCommentInfo.Exists(`i:` + cid).Result()
	return exist != 0
}

func DeleteCommentAndAllAbout(cid string) error {
	commentPipe := redisDBCommentInfo.TxPipeline()
	videoPipe := redisDBVideoInfo.TxPipeline()

	comment, err := GetCommentInfo(cid)
	if err != nil {
		return err
	}
	if comment.ParentId != `-1` {
		commentPipe.LRem(`c:`+comment.ParentId, 1, cid)
	}
	childList, err := redisDBCommentInfo.LRange(`c:`+cid, 0, -1).Result()
	if err != nil {
		return err
	}

	videoPipe.LRem(`c:`+comment.VideoId, 1, cid)
	commentPipe.Del(`c:`+cid, `i:`+cid, `l:`+cid)
	for _, item := range childList {
		videoPipe.LRem(`c:`+comment.VideoId, 1, item)
		commentPipe.Del(`c:`+item, `i:`+item, `l:`+item)
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
