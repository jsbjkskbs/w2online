package syncman

import (
	"context"
	"sync"
	"time"
	"work/biz/dal/db"
	"work/biz/mw/redis"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type RelationSyncman struct {
	ctx    context.Context
	cancle context.CancelFunc
}

func NewRelationSyncman() *RelationSyncman {
	ctx, cancle := context.WithCancel(context.Background())
	return &RelationSyncman{
		ctx:    ctx,
		cancle: cancle,
	}
}

func (sm RelationSyncman) Run() {
	if err := relationSyncMwWhenInit(); err != nil {
		panic(err)
	}
	go func() {
		for {
			time.Sleep(time.Minute * 10)
			select {
			case <-sm.ctx.Done():
				hlog.Info("Ok,stop sync[follow]")
			default:
			}
		}
	}()
}

func (sm RelationSyncman) Stop() {
	sm.cancle()
}

type relationSyncData struct {
	uid          string
	followList   *[]string
	followerList *[]string
}

func relationSyncMwWhenInit() error {
	list, err := db.GetUserIdList()
	if err != nil {
		panic(err)
	}

	var (
		wg       sync.WaitGroup
		errChan  = make(chan error, 2)
		syncList = make([]relationSyncData, 0)
		data     relationSyncData
	)
	for _, uid := range *list {
		data.uid = uid
		wg.Add(2)
		go func(data *relationSyncData) {
			if data.followList, err = db.GetFollowList(data.uid); err != nil {
				errChan <- err
			}
			wg.Done()
		}(&data)
		go func(data *relationSyncData) {
			if data.followerList, err = db.GetFollowerList(data.uid); err != nil {
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
	if err := followSyncDB2Redis(&syncList); err != nil {
		return err
	}
	return nil
}

func followSyncDB2Redis(syncList *[]relationSyncData) error {
	var (
		wg      sync.WaitGroup
		errChan = make(chan error, 2)
	)
	for _, item := range *syncList {
		wg.Add(2)
		go func(uid string, followList *[]string) {
			if err := redis.PutFollowList(uid, followList); err != nil {
				errChan <- err
			}
			wg.Done()
		}(item.uid, item.followList)
		go func(uid string, followerList *[]string) {
			if err := redis.PutFollowerList(uid, followerList); err != nil {
				errChan <- err
			}
			wg.Done()
		}(item.uid, item.followerList)
		wg.Wait()
		select {
		case result := <-errChan:
			return result
		default:
		}
	}
	return nil
}
