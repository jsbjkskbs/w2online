// Code generated by hertz generator.

package relation

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	relation "work/biz/model/base/relation"
)

// RelationAction .
// @router /relation/action/ [POST]
func RelationAction(ctx context.Context, c *app.RequestContext) {
	var err error
	var req relation.RelationActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(relation.RelationActionResponse)

	c.JSON(consts.StatusOK, resp)
}

// FollowingList .
// @router /following/list/ [GET]
func FollowingList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req relation.FollowingListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(relation.FollowingListResponse)

	c.JSON(consts.StatusOK, resp)
}

// FollowerList .
// @router /follower/list/ [GET]
func FollowerList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req relation.FollowerListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(relation.FollowerListResponse)

	c.JSON(consts.StatusOK, resp)
}

// FriendList .
// @router /friend/list/ [GET]
func FriendList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req relation.FriendListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(relation.FriendListResponse)

	c.JSON(consts.StatusOK, resp)
}