package webs

import (
	"context"
	"work/biz/mw/jwt"

	"github.com/cloudwego/hertz/pkg/app"
)

func _wsAuth() []app.HandlerFunc {
	return append(make([]app.HandlerFunc, 0),
		tokenAuthFunc(),
		//jwt.AccessTokenJwtMiddleware.MiddlewareFunc(),
	)
}

func tokenAuthFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if !jwt.IsAccessTokenAvailable(ctx, c) {
			c.Abort()
			return
		}
		c.Next(ctx)
	}
}
