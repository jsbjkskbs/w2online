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

type VideoSyncman struct {
	ctx    context.Context
	cancle context.CancelFunc
}

func NewVideoSyncman() *VideoSyncman {
	ctx, cancle := context.WithCancel(context.Background())
	return &VideoSyncman{
		ctx:    ctx,
		cancle: cancle,
	}
}

func (sm VideoSyncman) Run() {
	if err := videoSyncMwWhenInit(); err != nil {
		panic(err)
	}
	go func() {
		for {
			time.Sleep(time.Minute * 10)
			select {
			case <-sm.ctx.Done():
				hlog.Info("Ok,stop sync[video]")
				return
			default:
			}
			var (
				wg                       sync.WaitGroup
				errChan                  = make(chan error, 3)
				visitCount, commentCount int64
				likeList                 *map[string]string
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
					continue
				default:
				}
				likeCount := 0
				for uid, value := range *likeList {
					if value == `1` {
						likeCount++
						err := db.CreateIfNotExistsVideoLike(vid, uid)
						if err != nil {
							hlog.Error(err)
						}
					} else {
						err := db.DeleteVideoLike(vid, uid)
						if err != nil {
							hlog.Error(err)
						}
					}
				}
				if err := db.UpdateVideoVisit(vid, fmt.Sprint(visitCount)); err != nil {
					hlog.Error(err)
				}

				err := elasticsearch.UpdateVideoLikeVisitAndCommentCount(vid, fmt.Sprint(likeCount), fmt.Sprint(visitCount), fmt.Sprint(commentCount))
				if err != nil {
					hlog.Error(err)
				}
			}
		}
	}()
}

func (sm VideoSyncman) Stop() {
	sm.cancle()
}

type videoSyncData struct {
	vid         string
	likeList    *[]string
	visitCount  string
	commentList *[]string
}

func videoSyncMwWhenInit() error {
	list, err := db.GetVideoIdList()
	if err != nil {
		panic(err)
	}

	var (
		wg       sync.WaitGroup
		errChan  = make(chan error, 3)
		syncList = make([]videoSyncData, 0)
		data     videoSyncData
	)
	for _, vid := range *list {
		data.vid = vid
		wg.Add(3)
		go func(data *videoSyncData) {
			if data.likeList, err = db.GetVideoLikeList(vid); err != nil {
				errChan <- err
			}
			wg.Done()
		}(&data)
		go func(data *videoSyncData) {
			if data.visitCount, err = db.GetVideoVisitCount(vid); err != nil {
				errChan <- err
			}
			wg.Done()
		}(&data)
		go func(data *videoSyncData) {
			if data.commentList, err = db.GetVideoCommentList(vid); err != nil {
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
	go func(syncList *[]videoSyncData) {
		if err := videoSyncDB2Redis(syncList); err != nil {
			errChan <- err
		}
		wg.Done()
	}(&syncList)
	go func(syncList *[]videoSyncData) {
		if err := vidoeSyncDB2Elastic(syncList); err != nil {
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

func videoSyncDB2Redis(syncList *[]videoSyncData) error {
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
			if err := redis.PutVideoLikeInfo(vid, likeList); err != nil {
				errChan <- err
			}
			wg.Done()
		}(item.vid, item.likeList)
		go func(vid string, commentList *[]string) {
			if err := redis.PutVideoCommentInfo(vid, commentList); err != nil {
				errChan <- err
			}
			wg.Done()
		}(item.vid, item.commentList)
		wg.Wait()
		select {
		case result := <-errChan:
			return result
		default:
		}
	}
	return nil
}

func vidoeSyncDB2Elastic(syncList *[]videoSyncData) error {
	for _, item := range *syncList {
		if err := elasticsearch.UpdateVideoLikeVisitAndCommentCount(item.vid, fmt.Sprint(len(*item.likeList)), item.visitCount, fmt.Sprint(len(*item.commentList))); err != nil {
			return err
		}
	}
	return nil
}
