// Code generated by hertz generator.

package test

import (
	"context"

	test "work/biz/model/test"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Test .
// @router /test [POST]
func Test(ctx context.Context, c *app.RequestContext) {
	var err error
	var req test.TestRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	hlog.Info(string(req.Data))

	c.JSON(consts.StatusOK, test.TestResponse{
		Msg: "receive",
	})
}
