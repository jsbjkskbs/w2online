package syncman

import (
	"context"
	"fmt"
	"sync"
	"time"
	"work/biz/dal/db"
	"work/biz/mw/elasticsearch"
	"work/biz/mw/redis"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type Syncman struct {
	ctx    context.Context
	cancle context.CancelFunc
}

func NewSyncman() *Syncman {
	ctx, cancle := context.WithCancel(context.Background())
	return &Syncman{
		ctx:    ctx,
		cancle: cancle,
	}
}

func (sm Syncman) Run() {
	SyncMwWhenInit()
	go func() {
		for {
			time.Sleep(time.Minute * 10)
			select {
			case <-sm.ctx.Done():
				hlog.Info("Ok,stop sync")
				return
			default:
			}
			var (
				wg                       sync.WaitGroup
				errChan                  = make(chan error, 3)
				visitCount, commentCount int64
				likeList                 *[]string
				vidList                  *[]string
			)
			var err error
			if vidList, err = db.GetVideoIdList(); err != nil {
				hlog.Warn(err)
			}
			for _, vid := range *vidList {
				wg.Add(3)
				go func() {
					var err error
					if visitCount, err = redis.GetVideoVisitCount(vid); err != nil {
						errChan <- err
					}
					wg.Done()
				}()
				go func() {
					var err error
					if commentCount, err = redis.GetVideoCommentCount(vid); err != nil {
						errChan <- err
					}
					wg.Done()
				}()
				go func() {
					var err error
					if likeList, err = redis.GetVideoLikeList(vid); err != nil {
						errChan <- err
					}
					wg.Done()
				}()
				wg.Wait()
				select {
				case result := <-errChan:
					hlog.Error(result)
				default:
				}
				wg.Add(2)
				errChan = make(chan error, 2)
				go func() {
					if err := db.CreateIfNotExistsVideoLike(vid, likeList); err != nil {
						errChan <- err
					}
					wg.Done()
				}()
				go func() {
					err := elasticsearch.UpdateVideoLikeVisitAndCommentCount(vid, fmt.Sprint(len(*likeList)), fmt.Sprint(visitCount), fmt.Sprint(commentCount))
					if err != nil {
						errChan <- err
					}
					wg.Done()
				}()
				wg.Wait()
				select {
				case result := <-errChan:
					hlog.Error(result)
				default:
				}
			}
		}
	}()
}

func (sm Syncman) Stop() {
	sm.cancle()
}

type syncData struct {
	vid          string
	likeList     *[]string
	visitCount   string
	commentCount string
}

func SyncMwWhenInit() error {
	list, err := db.GetVideoIdList()
	if err != nil {
		panic(err)
	}

	var (
		wg       sync.WaitGroup
		errChan  = make(chan error, 3)
		syncList = make([]syncData, 0)
		data     syncData
	)
	for _, vid := range *list {
		data.vid = vid
		wg.Add(3)
		go func(data *syncData) {
			if data.likeList, err = db.GetVideoLikeList(vid); err != nil {
				errChan <- err
			}
			wg.Done()
		}(&data)
		go func(data *syncData) {
			if data.visitCount, err = db.GetVideoVisitCount(vid); err != nil {
				errChan <- err
			}
			wg.Done()
		}(&data)
		go func(data *syncData) {
			if data.commentCount, err = db.GetVideoCommentCount(vid); err != nil {
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

	errChan = make(chan error, 2)
	wg.Add(2)
	go func(syncList *[]syncData) {
		if err := syncDB2Redis(syncList); err != nil {
			errChan <- err
		}
		wg.Done()
	}(&syncList)
	go func(syncList *[]syncData) {
		if err := syncDB2Elastic(syncList); err != nil {
			errChan <- err
		}
		wg.Done()
	}(&syncList)
	wg.Wait()
	select {
	case result := <-errChan:
		return result
	default:
	}
	return nil
}

func syncDB2Redis(syncList *[]syncData) error {
	var (
		wg      sync.WaitGroup
		errChan = make(chan error, 3)
	)
	for _, item := range *syncList {
		wg.Add(3)
		go func(vid, visitCount string) {
			if err := redis.PutVideoVisitInfo(vid, visitCount); err != nil {
				errChan <- err
			}
			wg.Done()
		}(item.vid, item.visitCount)
		go func(vid string, likeList *[]string) {
			if err := redis.PutVideoLikeInfo(vid, *likeList); err != nil {
				errChan <- err
			}
			wg.Done()
		}(item.vid, item.likeList)
		go func(vid, commentCount string) {
			if err := redis.PutVideoCommentInfo(vid, commentCount); err != nil {
				errChan <- err
			}
			wg.Done()
		}(item.vid, item.commentCount)
		wg.Wait()
		select {
		case result := <-errChan:
			return result
		default:
		}
	}
	return nil
}

func syncDB2Elastic(syncList *[]syncData) error {
	for _, item := range *syncList {
		if err := elasticsearch.UpdateVideoLikeVisitAndCommentCount(item.vid, fmt.Sprint(len(*item.likeList)), item.visitCount, item.commentCount); err != nil {
			return err
		}
	}
	return nil
}
