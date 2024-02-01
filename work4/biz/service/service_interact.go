package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"work/biz/dal/db"
	"work/biz/model/base"
	"work/biz/model/base/interact"
	"work/biz/mw/elasticsearch"
	"work/biz/mw/jwt"
	"work/biz/mw/redis"
	"work/pkg/constants"
	"work/pkg/errmsg"
	"work/pkg/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

type InteractService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewInteractService(ctx context.Context, c *app.RequestContext) *InteractService {
	return &InteractService{ctx: ctx, c: c}
}

func (service InteractService) NewCommentPublishEvent(request *interact.CommentPublishRequest) error {
	uid, err := jwt.CovertJWTPayloadToString(service.ctx, service.c)
	if err != nil {
		return errmsg.AuthenticatorError
	}
	if request.Content == `` {
		return errmsg.ParamError
	}
	if request.CommentId == `` && request.VideoId == `` {
		return errmsg.ServiceError
	}
	if request.CommentId == `` {
		request.CommentId = `-1`
	} else {
		parentComment, err := redis.GetCommentInfo(request.CommentId)
		if err != nil {
			return errmsg.RedisError
		}
		if parentComment.ParentId != `-1` {
			request.CommentId = parentComment.ParentId
		}
	}
	if request.VideoId == `` {
		vid, err := redis.GetCommentVideoId(request.CommentId)
		if err != nil {
			return errmsg.RedisError
		}
		request.VideoId = vid
	} else {
		if !redis.IsVideoExist(request.VideoId) {
			return errmsg.ServiceError
		}
	}
	newComment := db.Comment{
		VideoId:   request.VideoId,
		ParentId:  request.CommentId,
		UserId:    uid,
		Content:   request.Content,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		DeletedAt: 0,
	}
	if err := db.CreateComment(&newComment); err != nil {
		return err
	}
	var (
		wg      sync.WaitGroup
		errChan = make(chan error, 3)
	)
	wg.Add(3)
	go func() {
		if err := redis.PutCommentInfo(&newComment); err != nil {
			errChan <- errmsg.RedisError
		}
		wg.Done()
	}()
	go func() {
		if err := redis.AppendVideoCommentInfo(newComment.VideoId, fmt.Sprint(newComment.Id)); err != nil {
			errChan <- errmsg.RedisError
		}
		wg.Done()
	}()
	go func() {
		if newComment.ParentId != `-1` {
			if err := redis.AppendChildCommentInfo(request.CommentId, fmt.Sprint(newComment.Id)); err != nil {
				errChan <- errmsg.RedisError
			}
		}
		wg.Done()
	}()
	wg.Wait()
	select {
	case err := <-errChan:
		return err
	default:
	}
	return nil
}

func (service InteractService) NewLikeActionEvent(request *interact.LikeActionRequest) error {
	uid, err := jwt.CovertJWTPayloadToString(service.ctx, service.c)
	if err != nil {
		return errmsg.AuthenticatorError
	}
	if request.VideoId != `` {
		if !redis.IsVideoExist(request.VideoId) {
			return errmsg.ServiceError
		}
		switch request.ActionType {
		case `1`:
			{
				if err := redis.AppendVideoLikeInfo(request.VideoId, uid); err != nil {
					return errmsg.RedisError
				}
			}
		case `2`:
			{
				if err := redis.RemoveVideoLikeInfo(request.VideoId, uid); err != nil {
					return errmsg.RedisError
				}
			}
		}
	} else if request.CommentId != `` {
		if !redis.IsCommentExist(request.CommentId) {
			return errmsg.ServiceError
		}
		switch request.ActionType {
		case `1`:
			{
				if err := redis.AppendCommentLikeInfo(request.CommentId, uid); err != nil {
					return errmsg.RedisError
				}
			}
		case `2`:
			{
				if err := redis.RemoveCommentLikeInfo(request.CommentId, uid); err != nil {
					return errmsg.RedisError
				}
			}
		}
	} else {
		return errmsg.ParamError
	}
	return nil
}

func (service InteractService) NewLikeListEvent(request *interact.LikeListRequest) (*interact.LikeListResponse_LikeListResponseData, error) {
	if request.PageNum <= 0 {
		request.PageNum = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = constants.DefaultPageSize
	}
	list, err := db.GetVideoLikeListByUserId(request.UserId, request.PageNum, request.PageSize)
	if err != nil {
		return nil, errmsg.ServiceError
	}
	data := make([]*base.Video, len(*list))
	for i, item := range *list {
		if data[i], err = elasticsearch.GetVideoDoc(item); err != nil {
			return nil, errmsg.ElasticError
		}
	}
	return &interact.LikeListResponse_LikeListResponseData{Items: data}, nil
}

func (service InteractService) NewCommentListEvent(request *interact.CommentListRequest) (*interact.CommentListResponse_CommentListResponseData, error) {
	if request.PageNum <= 0 {
		request.PageNum = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = constants.DefaultPageSize
	}

	var (
		data *[]*base.Comment
		err  error
	)
	if request.VideoId != `` {
		if !redis.IsVideoExist(request.VideoId) {
			return nil, errmsg.ServiceError
		}
		if data, err = getVideoComment(request); err != nil {
			return nil, err
		}
	} else if request.CommentId != `` {
		if !redis.IsCommentExist(request.CommentId) {
			return nil, errmsg.ServiceError
		}
		if data, err = getCommentComment(request); err != nil {
			return nil, err
		}
	} else {
		return nil, errmsg.ParamError
	}
	return &interact.CommentListResponse_CommentListResponseData{Items: *data}, nil
}

func (service InteractService) NewDeleteEvent(request *interact.CommentDeleteRequest) error {
	uid, err := jwt.CovertJWTPayloadToString(service.ctx, service.c)
	if err != nil {
		return errmsg.AuthenticatorError
	}
	if request.VideoId != `` {
		if !redis.IsVideoExist(request.VideoId) {
			return errmsg.ServiceError
		}
		videoInfo, _ := elasticsearch.GetVideoDoc(request.VideoId)
		if videoInfo.UserId != uid {
			return errmsg.ServiceError
		}
		if err := deleteVideo(request); err != nil {
			return err
		}
	} else if request.CommentId != `` {
		if !redis.IsCommentExist(request.CommentId) {
			return errmsg.ServiceError
		}
		commentInfo, _ := redis.GetCommentInfo(request.CommentId)
		if commentInfo.UserId != uid {
			return errmsg.ServiceError
		}
		if err := deleteComment(request); err != nil {
			return err
		}
	} else {
		return errmsg.ParamError
	}
	return nil
}

func getVideoComment(request *interact.CommentListRequest) (*[]*base.Comment, error) {
	data := make([]*base.Comment, 0)
	list, err := redis.GetVideoCommentList(request.VideoId, request.PageNum, request.PageSize)
	if err != nil {
		return nil, errmsg.RedisError
	}
	for _, item := range *list {
		d, err := redis.GetCommentInfo(item)
		if err != nil {
			return nil, errmsg.RedisError
		}
		likeCount, err := redis.GetCommentLikeCount(item)
		if err != nil {
			return nil, errmsg.RedisError
		}
		childCount, err := redis.GetCommentChildCount(item)
		if err != nil {
			return nil, errmsg.RedisError
		}
		data = append(data, &base.Comment{
			Id:         fmt.Sprint(d.Id),
			UserId:     d.UserId,
			VideoId:    d.VideoId,
			ParentId:   d.ParentId,
			LikeCount:  likeCount,
			ChildCount: childCount,
			Content:    d.Content,
			CreatedAt:  utils.ConvertTimestampToStringDefault(d.CreatedAt),
			UpdatedAt:  utils.ConvertTimestampToStringDefault(d.UpdatedAt),
			DeletedAt:  utils.ConvertTimestampToStringDefault(d.DeletedAt),
		})
	}
	return &data, nil
}

func getCommentComment(request *interact.CommentListRequest) (*[]*base.Comment, error) {
	data := make([]*base.Comment, 0)
	list, err := redis.GetCommentChildList(request.CommentId, request.PageNum, request.PageSize)
	if err != nil {
		return nil, errmsg.RedisError
	}
	for _, item := range *list {
		d, err := redis.GetCommentInfo(item)
		if err != nil {
			return nil, errmsg.RedisError
		}
		likeCount, err := redis.GetCommentLikeCount(item)
		if err != nil {
			return nil, errmsg.RedisError
		}
		childCount, err := redis.GetCommentChildCount(item)
		if err != nil {
			return nil, errmsg.RedisError
		}
		data = append(data, &base.Comment{
			Id:         fmt.Sprint(d.Id),
			UserId:     d.UserId,
			VideoId:    d.VideoId,
			ParentId:   d.ParentId,
			LikeCount:  likeCount,
			ChildCount: childCount,
			CreatedAt:  utils.ConvertTimestampToStringDefault(d.CreatedAt),
			UpdatedAt:  utils.ConvertTimestampToStringDefault(d.UpdatedAt),
			DeletedAt:  utils.ConvertTimestampToStringDefault(d.DeletedAt),
		})
	}
	return &data, nil
}

func deleteVideo(request *interact.CommentDeleteRequest) error {
	if err := db.DeleteVideo(request.VideoId); err != nil {
		return errmsg.ServiceError
	}
	var (
		wg      sync.WaitGroup
		errChan = make(chan error, 3)
	)
	wg.Add(4)
	go func() {
		if err := db.DeleteVideo(request.VideoId); err != nil {
			errChan <- errmsg.ServiceError
		}
		wg.Done()
	}()
	go func() {
		if err := db.DeleteCommentAndCommentLikeAboutVideo(request.VideoId); err != nil {
			errChan <- errmsg.ServiceError
		}
		wg.Done()
	}()
	go func() {
		if err := redis.DeleteVideoAndAllAbout(request.VideoId); err != nil {
			errChan <- errmsg.RedisError
		}
		wg.Done()
	}()
	go func() {
		if err := elasticsearch.DeleteVideoDoc(request.VideoId); err != nil {
			errChan <- errmsg.ElasticError
		}
		wg.Done()
	}()
	wg.Wait()
	select {
	case err := <-errChan:
		return err
	default:
	}
	return nil
}

func deleteComment(request *interact.CommentDeleteRequest) error {
	if err := db.DeleteComment(request.CommentId); err != nil {
		return errmsg.ServiceError
	}
	var (
		wg      sync.WaitGroup
		errChan = make(chan error, 2)
	)
	wg.Add(2)
	go func() {
		if err := db.DeleteChildAndLikesOfParentAndChild(request.CommentId); err != nil {
			errChan <- errmsg.RedisError
		}
		wg.Done()
	}()
	go func() {
		if err := redis.DeleteCommentAndAllAbout(request.CommentId); err != nil {
			errChan <- errmsg.RedisError
		}
		wg.Done()
	}()
	wg.Wait()
	select {
	case err := <-errChan:
		return err
	default:
	}
	return nil
}
