package authfunc

import (
	"context"
	"work/biz/model/base"
	"work/biz/model/base/user"
	"work/biz/mw/jwt"
	"work/pkg/errmsg"
	"work/pkg/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Auth() []app.HandlerFunc {
	return append(make([]app.HandlerFunc, 0),
		DoubleTokenAuthFunc(),
		//jwt.AccessTokenJwtMiddleware.MiddlewareFunc(),
	)
}

func DoubleTokenAuthFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if !jwt.IsAccessTokenAvailable(ctx, c) {
			if !jwt.IsRefreshTokenAvailable(ctx, c) {
				resp := utils.CreateBaseHttpResponse(errmsg.TokenIsInavailableError)
				c.JSON(consts.StatusOK, user.UserLoginResponse{
					Base: &base.Status{
						Code: resp.StatusCode,
						Msg:  resp.StatusMsg,
					},
				})
				c.Abort()
				return
			}
			jwt.GenerateAccessToken(ctx, c)
		}

		c.Next(ctx)
	}
}
