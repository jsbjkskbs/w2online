package syncman

import (
	"context"
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
				likeList, err := redis.GetNewUpdateCommentLikeList(cid)
				if err != nil {
					hlog.Error(err)
					continue
				}
				for _, uid := range *likeList {
					if err := db.CreateCommentLike(cid, uid); err != nil {
						hlog.Error(err)
					}
					if err := redis.AppendCommentLikeInfoToStaticSpace(cid, uid); err != nil {
						hlog.Error(err)
					}
					if err := redis.DeleteCommentLikeInfoFromDynamicSpace(cid, uid); err != nil {
						hlog.Error(err)
					}
				}
				dislikeList, err := redis.GetNewDeleteCommentLikeList(cid)
				if err != nil {
					hlog.Error(err)
					continue
				}
				for _, uid := range *dislikeList {
					if err := db.DeleteCommentLike(cid, uid); err != nil {
						hlog.Error(err)
					}
					if err := redis.DeleteCommentLikeInfoFromDynamicSpace(cid, uid); err != nil {
						hlog.Error(err)
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
	cid      string
	likeList *[]string
}

func commentSyncMwWhenInit() error {
	list, err := db.GetCommentIdList()
	if err != nil {
		panic(err)
	}

	var (
		syncList = make([]commentSyncData, 0)
		data     commentSyncData
	)
	for _, cid := range *list {
		data.cid = cid
		if data.likeList, err = db.GetCommentLikeList(data.cid); err != nil {
			return err
		}
		syncList = append(syncList, data)
	}
	if err := commentSyncDB2Redis(&syncList); err != nil {
		return err
	}
	return nil
}

func commentSyncDB2Redis(syncList *[]commentSyncData) error {
	for _, item := range *syncList {
		if err := redis.PutCommentLikeInfo(item.cid, item.likeList); err != nil {
			return err
		}
	}
	return nil
}
