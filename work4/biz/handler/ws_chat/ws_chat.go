package wschat

import (
	"context"
	"work/biz/mw/jwt"
	"work/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/websocket"
)

var upgrader = websocket.HertzUpgrader{}

var (
	BadConnection = []byte(`bad connection`)
)

func Handler(ctx context.Context, c *app.RequestContext) {
	err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
		uid, err := jwt.CovertJWTPayloadToString(ctx, c)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, BadConnection)
			return
		}
		conn.WriteMessage(websocket.TextMessage, []byte(`Welcome, `+uid))

		s := service.NewChatService(ctx, c, conn)

		if err := s.Login(); err != nil {
			conn.WriteMessage(websocket.TextMessage, BadConnection)
			return
		}
		defer s.Logout()

		if err := s.ReadOfflineMessage(); err != nil {
			conn.WriteMessage(websocket.TextMessage, BadConnection)
			return
		}

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, BadConnection)
				return
			}

			if err := s.SendMessage(string(message)); err != nil {
				conn.WriteMessage(websocket.TextMessage, BadConnection)
				return
			}
		}
	})

	if err != nil {
		c.JSON(consts.StatusOK, `error`)
		return
	}
}
