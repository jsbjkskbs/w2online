package service

import (
	"context"
	"work/biz/dal/db"
	"work/biz/model/base"
	"work/biz/model/base/relation"
	"work/biz/mw/jwt"
	"work/pkg/constants"
	"work/pkg/errmsg"

	"github.com/cloudwego/hertz/pkg/app"
)

type RelationService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewRelationService(ctx context.Context, c *app.RequestContext) *RelationService {
	return &RelationService{
		ctx: ctx,
		c:   c,
	}
}

func (service RelationService) NewRelationActionEvent(request *relation.RelationActionRequest) error {
	uid, err := jwt.CovertJWTPayloadToString(service.ctx, service.c)
	if err != nil {
		return errmsg.AuthenticatorError
	}
	if uid == request.ToUserId {
		return errmsg.ParamError
	}
	switch request.ActionType {
	case 0:
		if err := createFollow(uid, request); err != nil {
			return err
		}
	case 1:
		if err := cancleFollow(uid, request); err != nil {
			return err
		}
	}
	return nil
}

func (service RelationService) NewFollowingListEvent(request *relation.FollowingListRequest) (*relation.FollowingListResponse_FollowingListResponseData, error) {
	exist, err := db.UserIsExistByUid(request.UserId)
	if err != nil {
		return nil, errmsg.ServiceError
	}
	if !exist {
		return nil, errmsg.UserDoesNotExistError
	}
	if request.PageNum <= 0 {
		request.PageNum = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = constants.DefaultPageSize
	}
	list, err := db.GetFollowListPaged(request.UserId, request.PageNum, request.PageSize)
	if err != nil {
		return nil, errmsg.ServiceError
	}
	data := make([]*base.UserLite, 0)
	for _, item := range *list {
		user, err := db.QueryUserByUid(item)
		if err != nil {
			return nil, errmsg.ServiceError
		}
		d := base.UserLite{
			Uid:       item,
			Username:  user.Username,
			AvatarUrl: user.AvatarUrl,
		}
		data = append(data, &d)
	}
	total, err := db.GetFollowCount(request.UserId)
	if err != nil {
		return nil, errmsg.ServiceError
	}
	return &relation.FollowingListResponse_FollowingListResponseData{Items: data, Total: total}, nil
}

func (service RelationService) NewFollowerEvent(request *relation.FollowerListRequest) (*relation.FollowerListResponse_FollowerListResponseData, error) {
	exist, err := db.UserIsExistByUid(request.UserId)
	if err != nil {
		return nil, errmsg.ServiceError
	}
	if !exist {
		return nil, errmsg.UserDoesNotExistError
	}
	if request.PageNum <= 0 {
		request.PageNum = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = constants.DefaultPageSize
	}
	list, err := db.GetFollowerListPaged(request.UserId, request.PageNum, request.PageSize)
	if err != nil {
		return nil, errmsg.ServiceError
	}
	data := make([]*base.UserLite, 0)
	for _, item := range *list {
		user, err := db.QueryUserByUid(item)
		if err != nil {
			return nil, errmsg.ServiceError
		}
		d := base.UserLite{
			Uid:       item,
			Username:  user.Username,
			AvatarUrl: user.AvatarUrl,
		}
		data = append(data, &d)
	}
	total, err := db.GetFollowerCount(request.UserId)
	if err != nil {
		return nil, errmsg.ServiceError
	}
	return &relation.FollowerListResponse_FollowerListResponseData{Items: data, Total: total}, nil
}

func (service RelationService) NewFriendListEvent(request *relation.FriendListRequest) (*relation.FriendListResponse_FriendListResponseData, error) {
	uid, err := jwt.CovertJWTPayloadToString(service.ctx, service.c)
	if err != nil {
		return nil, errmsg.AuthenticatorError
	}
	if request.PageNum <= 0 {
		request.PageNum = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = constants.DefaultPageSize
	}
	list, err := db.GetFriendListPaged(uid, request.PageNum, request.PageSize)
	if err != nil {
		return nil, errmsg.ServiceError
	}
	data := make([]*base.UserLite, 0)
	for _, item := range *list {
		user, err := db.QueryUserByUid(item)
		if err != nil {
			return nil, errmsg.ServiceError
		}
		d := base.UserLite{
			Uid:       item,
			Username:  user.Username,
			AvatarUrl: user.AvatarUrl,
		}
		data = append(data, &d)
	}
	total, err := db.GetFriendCount(uid)
	if err != nil {
		return nil, errmsg.ServiceError
	}
	return &relation.FriendListResponse_FriendListResponseData{Items: data, Total: total}, nil

}

func createFollow(uid string, request *relation.RelationActionRequest) error {
	if err := db.CreateFollow(request.ToUserId, uid); err != nil {
		return errmsg.ServiceError
	}
	return nil
}

func cancleFollow(uid string, request *relation.RelationActionRequest) error {
	if err := db.DeleteFollow(request.ToUserId, uid); err != nil {
		return errmsg.ServiceError
	}
	return nil
}
