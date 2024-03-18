package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
	"work/biz/dal/db"
	"work/biz/model/base"
	"work/biz/model/base/video"
	"work/biz/mw/elasticsearch"
	"work/biz/mw/jwt"
	"work/biz/mw/redis"
	"work/pkg/constants"
	"work/pkg/errmsg"
	qiniuyunoss "work/pkg/qiniuyun_oss"
	"work/pkg/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

type VideoService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewVideoService(ctx context.Context, c *app.RequestContext) *VideoService {
	return &VideoService{
		ctx: ctx,
		c:   c,
	}
}

func (service VideoService) NewCancleUploadEvent(request *video.VideoPublishCancleRequest) error {
	uid, err := jwt.CovertJWTPayloadToString(service.ctx, service.c)
	if err != nil {
		return errmsg.TokenIsInavailableError
	}
	if request.Uuid == `` {
		return errmsg.ParamError
	}
	if err := service.deleteTempDir(`./pkg/data/temp/video/` + uid + `_` + request.Uuid); err != nil {
		return errmsg.ServiceError
	}
	if _, err = redis.FinishVideoEvent(request.Uuid, uid); err != nil {
		return errmsg.RedisError
	}
	return nil
}

func (service VideoService) NewUploadCompleteEvent(request *video.VideoPublishCompleteRequest) error {
	uid, err := jwt.CovertJWTPayloadToString(service.ctx, service.c)
	if err != nil {
		return errmsg.TokenIsInavailableError
	}

	if request.Uuid == `` {
		return errmsg.ParamError
	}

	reallyComplete, err := redis.IsChunkAllRecorded(request.Uuid, uid)
	if err != nil {
		return errmsg.RedisError
	}
	if !reallyComplete {
		return errmsg.ParamError
	}

	m3u8name, err := redis.GetM3U8Filename(request.Uuid, uid)
	if err != nil {
		return errmsg.RedisError
	}

	err = utils.M3u8ToMp4(`./pkg/data/temp/video/`+uid+`_`+request.Uuid+`/`+m3u8name,
		`./pkg/data/temp/video/`+uid+`_`+request.Uuid+`/`+`video.mp4`)
	if err != nil {
		return errmsg.ServiceError
	}

	err = utils.GenerateMp4CoverJpg(`./pkg/data/temp/video/`+uid+`_`+request.Uuid+`/`+`video.mp4`,
		`./pkg/data/temp/video/`+uid+`_`+request.Uuid+`/`+`cover.jpg`)
	if err != nil {
		return errmsg.ServiceError
	}

	info, err := redis.FinishVideoEvent(request.Uuid, uid)
	if err != nil {
		return errmsg.RedisError
	}

	uidInt64, _ := strconv.Atoi(uid)
	d := db.Video{
		Title:       info[0],
		Description: info[1],
		UserId:      int64(uidInt64),
		VisitCount:  0,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		DeletedAt:   0,
	}
	vid, err := db.CreateVideo(&d)
	if err != nil {
		return errmsg.ServiceError
	}

	var (
		videoUrl, coverUrl string
		user               *db.User
		wg                 sync.WaitGroup
	)
	errChan := make(chan error, 3)
	wg.Add(3)
	go func() {
		videoUrl, err = qiniuyunoss.UploadVideo(`./pkg/data/temp/video/`+uid+`_`+request.Uuid+`/`+`video.mp4`, vid)
		if err != nil {
			errChan <- errmsg.OssUploadError
		}
		wg.Done()
	}()
	go func() {
		coverUrl, err = qiniuyunoss.UploadVideoCover(`./pkg/data/temp/video/`+uid+`_`+request.Uuid+`/`+`cover.jpg`, vid)
		if err != nil {
			errChan <- errmsg.OssUploadError
		}
		wg.Done()
	}()
	go func() {
		user, err = db.QueryUserByUid(uid)
		if err != nil {
			errChan <- errmsg.ServiceError
		}
		wg.Done()
	}()
	wg.Wait()
	select {
	case err := <-errChan:
		return err
	default:
	}

	err = db.UpdateVideoUrl(videoUrl, coverUrl, vid)
	if err != nil {
		return errmsg.ServiceError
	}

	err = elasticsearch.CreateVideoDoc(&elasticsearch.Video{
		Title:       d.Title,
		Description: d.Description,
		Username:    user.Username,
		UserId:      uid,
		CreatedAt:   d.CreatedAt,
		Info: elasticsearch.VideoOtherdata{
			Id:           vid,
			VideoUrl:     videoUrl,
			CoverUrl:     coverUrl,
			UpdatedAt:    d.UpdatedAt,
			DeletedAt:    d.DeletedAt,
			VisitCount:   d.VisitCount,
			LikeCount:    0,
			CommentCount: 0,
		},
	})
	if err != nil {
		return errmsg.ElasticError
	}

	err = redis.DeleteVideoEvent(request.Uuid, uid)
	if err != nil {
		return errmsg.RedisError
	}

	err = service.deleteTempDir(`./pkg/data/temp/video/` + uid + `_` + request.Uuid)
	if err != nil {
		return errmsg.ServiceError
	}

	err = redis.PutVideoVisitInfo(vid, `0`)
	if err != nil {
		return errmsg.RedisError
	}
	err = redis.PutVideoLikeInfo(vid, nil)
	if err != nil {
		return errmsg.RedisError
	}

	return nil
}

func (service VideoService) NewUploadingEvent(request *video.VideoPublishUploadingRequest) error {
	uid, err := jwt.CovertJWTPayloadToString(service.ctx, service.c)
	if err != nil {
		return errmsg.TokenIsInavailableError
	}
	if request.Filename == `` || request.Uuid == `` || request.ChunkNumber <= 0 {
		return errmsg.ParamError
	}

	rawData, err := service.c.FormFile("data")
	if err != nil {
		return errmsg.FileIsUnableToBeCatchError
	}

	file, err := rawData.Open()
	if err != nil {
		return errmsg.ServiceError
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return errmsg.ServiceError
	}

	if !service.isMD5Same(data, request.Md5) {
		return errmsg.FileMD5IsNotMatchError
	}

	if request.IsM3U8 {
		err := redis.RecordM3U8Filename(request.Uuid, uid, request.Filename)
		if err != nil {
			return errmsg.RedisError
		}
	}

	err = service.saveTempData(`./pkg/data/temp/video/`+uid+`_`+request.Uuid+`/`+request.Filename, data)
	if err != nil {
		return errmsg.ServiceError
	}

	if err := redis.DoneChunkEvent(request.Uuid, uid, request.ChunkNumber); err != nil {
		return errmsg.RedisError
	}
	return nil
}

func (service VideoService) NewUploadEvent(request *video.VideoPublishStartRequest) (string, error) {
	uuid := ``
	uid, err := jwt.CovertJWTPayloadToString(service.ctx, service.c)
	if err != nil {
		return uuid, errmsg.TokenIsInavailableError
	}
	if request.Title == `` || request.ChunkTotalNumber <= 0 {
		return ``, errmsg.ParamError
	}
	uuid, err = redis.NewVideoEvent(request.Title, request.Description, uid, fmt.Sprint(request.ChunkTotalNumber))
	if err != nil {
		return ``, errmsg.RedisError
	}
	os.Mkdir(`./pkg/data/temp/video/`+uid+`_`+uuid, os.ModePerm)
	return uuid, nil
}

func (service VideoService) NewSearchEvent(request *video.VideoSearchRequest) (*video.VideoSearchResponse_VideoSearchResponseData, error) {
	var wg sync.WaitGroup
	wg.Add(6)
	go func() {
		if _, exist := service.c.Get(`keywords`); !exist {
			request.Keywords = constants.ESNoKeywordsFlag
		}
		wg.Done()
	}()
	go func() {
		if _, exist := service.c.Get(`page_num`); !exist {
			request.PageNum = constants.ESNoPageParamFlag
		}
		wg.Done()
	}()
	go func() {
		if _, exist := service.c.Get(`page_size`); !exist {
			request.PageSize = constants.ESNoPageParamFlag
		}
		wg.Done()
	}()
	go func() {
		if _, exist := service.c.Get(`from_date`); !exist {
			request.FromDate = constants.ESNoTimeFilterFlag
		}
		wg.Done()
	}()
	go func() {
		if _, exist := service.c.Get(`to_date`); !exist {
			request.ToDate = constants.ESNoTimeFilterFlag
		}
		wg.Done()
	}()
	go func() {
		if _, exist := service.c.Get(`username`); !exist {
			request.Username = constants.ESNoUsernameFilterFlag
		}
		wg.Done()
	}()
	wg.Wait()
	items, total, err := elasticsearch.SearchVideoDoc(
		request.Keywords,
		request.Username,
		request.PageSize, request.PageNum,
		request.FromDate, request.ToDate,
	)
	if err != nil {
		return nil, errmsg.ElasticError
	}
	return &video.VideoSearchResponse_VideoSearchResponseData{Items: items, Total: total}, nil
}

func (service VideoService) NewFeedEvent(request *video.VideoFeedRequest) (*video.VideoFeedResponse_VideoFeedResponseData, error) {
	var timestamp int64
	if len(request.LatestTime) == 0 {
		timestamp = 0
	} else {
		timestamp, _ = strconv.ParseInt(request.LatestTime, 10, 64)
	}
	items, _, err := elasticsearch.RandomVideoDoc(timestamp)
	if err != nil {
		return nil, errmsg.ElasticError
	}
	return &video.VideoFeedResponse_VideoFeedResponseData{Items: items}, err
}

func (service VideoService) NewListEvent(request *video.VideoListRequest) (*video.VideoListResponse_VideoListResponseData, error) {
	var wg sync.WaitGroup
	errChan := make(chan error, 3)
	wg.Add(3)
	go func() {
		if _, exist := service.c.Get(`page_num`); !exist {
			request.PageNum = constants.ESNoPageParamFlag
		}
		wg.Done()
	}()
	go func() {
		if _, exist := service.c.Get(`page_size`); !exist {
			request.PageSize = constants.ESNoPageParamFlag
		}
		wg.Done()
	}()
	go func() {
		if _, exist := service.c.Get(`user_id`); !exist {
			errChan <- errmsg.ServiceError
		}
		wg.Done()
	}()
	wg.Wait()
	select {
	case err := <-errChan:
		return nil, err
	default:
	}
	items, total, err := elasticsearch.SearchVideoDocByUserId(request.UserId, request.PageNum, request.PageSize)
	if err != nil {
		return nil, errmsg.ElasticError
	}
	return &video.VideoListResponse_VideoListResponseData{Items: items, Total: total}, nil
}

func (service VideoService) NewVisitEvent(request *video.VideoVisitRequest) (*base.Video, error) {
	vid := service.c.Param(`id`)
	if err := redis.IncrVideoVisitInfo(vid); err != nil {
		return nil, errmsg.RedisError
	}
	info, err := elasticsearch.GetVideoDoc(vid)
	if err != nil {
		return nil, errmsg.ElasticError
	}
	return info, nil
}

func (service VideoService) NewPopularEvent(request *video.VideoPopularRequest) (*video.VideoPopularResponse_VideoPopularResponseData, error) {
	if request.PageNum <= 0 {
		request.PageNum = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = constants.DefaultPageSize
	}
	list, err := redis.GetVideoPopularList(request.PageNum, request.PageSize)
	if err != nil {
		return nil, err
	}
	items := make([]*base.Video, len(*list))
	var (
		wg      sync.WaitGroup
		errChan = make(chan error, len(*list))
	)
	wg.Add(len(*list))
	for i, item := range *list {
		go func(i int, item string) {
			items[i], err = elasticsearch.GetVideoDoc(item)
			if err != nil {
				errChan <- errmsg.ElasticError
			}
			wg.Done()
		}(i, item)
	}
	wg.Wait()
	select {
	case err := <-errChan:
		return nil, err
	default:
	}
	return &video.VideoPopularResponse_VideoPopularResponseData{Items: items}, nil
}

func (service VideoService) deleteTempDir(path string) error {
	return os.RemoveAll(path)
}

func (service VideoService) saveTempData(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0777)
}

func (service VideoService) isMD5Same(data []byte, md5 string) bool {
	return utils.GetBytesMD5(data) == md5
}
