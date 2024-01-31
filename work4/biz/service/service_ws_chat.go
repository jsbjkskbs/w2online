package service

import (
	"context"
	"time"
	"work/biz/dal/db"
	"work/biz/mw/jwt"
	"work/biz/mw/redis"
	"work/pkg/errmsg"
	"work/pkg/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/websocket"
)

type ChatService struct {
	ctx  context.Context
	c    *app.RequestContext
	conn *websocket.Conn
}

type _user struct {
	username string
	conn     *websocket.Conn
}

var userMap = make(map[string]*_user)

func NewChatService(ctx context.Context, c *app.RequestContext, conn *websocket.Conn) *ChatService {
	return &ChatService{
		ctx:  ctx,
		c:    c,
		conn: conn,
	}
}

func (service ChatService) Login() error {
	uid, err := jwt.CovertJWTPayloadToString(service.ctx, service.c)
	if err != nil {
		return err
	}
	user, err := db.QueryUserByUid(uid)
	if err != nil {
		return err
	}
	userMap[uid] = &_user{conn: service.conn, username: user.Username}
	return nil
}

func (service ChatService) Logout() error {
	uid, err := jwt.CovertJWTPayloadToString(service.ctx, service.c)
	if err != nil {
		return err
	}
	userMap[uid] = nil
	return nil
}

func (service ChatService) SendMessage(content string) error {
	from, err := jwt.CovertJWTPayloadToString(service.ctx, service.c)
	if err != nil {
		return errmsg.AuthenticatorError
	}
	to := service.c.Query(`to_user_id`)
	exist, err := db.UserIsExistByUid(to)
	if err != nil {
		return errmsg.ServiceError
	}
	if !exist {
		return errmsg.UserDoesNotExistError
	}
	toConn := userMap[to]
	switch toConn {
	case nil: //离线
		{
			if err := redis.PushMessage(from, to, content, time.Now().Unix()); err != nil {
				return errmsg.RedisError
			}
		}
	default: // 在线
		{
			err := toConn.conn.WriteMessage(websocket.TextMessage,
				newWebsocketMessage(from,
					content,
					utils.ConvertTimestampToStringDefault(time.Now().Unix())),
			)
			if err != nil {
				return errmsg.ServiceError
			}

		}
	}
	return nil
}

func (service ChatService) ReadOfflineMessage() error {
	uid, err := jwt.CovertJWTPayloadToString(service.ctx, service.c)
	if err != nil {
		return errmsg.AuthenticatorError
	}
	for !redis.IsMessageQueueEmpty(uid) {
		from, content, when, err := redis.PopMessage(uid)
		if err != nil {
			return errmsg.RedisError
		}
		if err := service.conn.WriteMessage(websocket.TextMessage, newWebsocketMessage(from, content, when)); err != nil {
			return errmsg.ServiceError
		}
	}
	return nil
}

func newWebsocketMessage(uid, content, when string) []byte {
	return []byte(when + ` [` + uid + `]: ` + content)
}
