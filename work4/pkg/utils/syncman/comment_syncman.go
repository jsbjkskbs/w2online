package syncman

import (
	"context"
	"sync"
	"time"
	"work/biz/dal/db"
	"work/biz/mw/redis"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type CommentSyncman struct {
	ctx    context.Context
	cancle context.CancelFunc
}

func NewCommentSyncman() *CommentSyncman {
	ctx, cancle := context.WithCancel(context.Background())
	return &CommentSyncman{
		ctx:    ctx,
		cancle: cancle,
	}
}

func (sm CommentSyncman) Run() {
	if err := commentSyncMwWhenInit(); err != nil {
		panic(err)
	}
	go func() {
		for {
			time.Sleep(time.Minute * 10)
			select {
			case <-sm.ctx.Done():
				hlog.Info("Ok,stop sync[comment]")
				return
			default:
			}
			cidList, err := db.GetCommentIdList()
			if err != nil {
				hlog.Warn(err)
			}
			for _, cid := range *cidList {
				likeList, err := redis.GetCommentLikeList(cid)
				if err != nil {
					hlog.Error(err)
					continue
				}
				for uid, value := range *likeList {
					if value == `1` {
						err := db.CreateIfNotExistsCommentLike(cid, uid)
						if err != nil {
							hlog.Error(err)
						}
					} else {
						err := db.DeleteCommentLike(cid, uid)
						if err != nil {
							hlog.Error(err)
						}
					}
				}
			}
		}
	}()
}

func (sm CommentSyncman) Stop() {
	sm.cancle()
}

type commentSyncData struct {
	cid       string
	likeList  *[]string
	childList *[]string
}

func commentSyncMwWhenInit() error {
	list, err := db.GetCommentIdList()
	if err != nil {
		panic(err)
	}

	var (
		wg       sync.WaitGroup
		errChan  = make(chan error, 2)
		syncList = make([]commentSyncData, 0)
		data     commentSyncData
	)
	for _, cid := range *list {
		data.cid = cid
		wg.Add(2)
		go func(data *commentSyncData) {
			if data.likeList, err = db.GetCommentLikeList(data.cid); err != nil {
				errChan <- err
			}
			wg.Done()
		}(&data)
		go func(data *commentSyncData) {
			if data.childList, err = db.GetCommentChildList(data.cid); err != nil {
				errChan <- err
			}
			wg.Done()
		}(&data)
		wg.Wait()
		select {
		case result := <-errChan:
			return result
		default:
		}
		syncList = append(syncList, data)
	}
	if err := commentSyncDB2Redis(&syncList); err != nil {
		return err
	}
	return nil
}

func commentSyncDB2Redis(syncList *[]commentSyncData) error {
	var (
		wg      sync.WaitGroup
		errChan = make(chan error, 3)
	)
	for _, item := range *syncList {
		wg.Add(2)
		go func(cid string, childList *[]string) {
			if err := redis.PutCommentChildInfo(cid, childList); err != nil {
				errChan <- err
			}
			wg.Done()
		}(item.cid, item.childList)
		go func(cid string, likeList *[]string) {
			if err := redis.PutCommentLikeInfo(cid, likeList); err != nil {
				errChan <- err
			}
			wg.Done()
		}(item.cid, item.likeList)
		wg.Wait()
		select {
		case result := <-errChan:
			return result
		default:
		}
	}
	return nil
}
