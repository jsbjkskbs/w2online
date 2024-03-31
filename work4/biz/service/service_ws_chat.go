package service

import (
	"context"
	"time"
	"work/biz/dal/db"
	"work/biz/mw/jwt"
	"work/pkg/errmsg"
	"work/pkg/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
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
	rsa      *utils.RsaService
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
	rsaClientKey := service.c.GetHeader(`rsa_public_key`)
	r := utils.NewRsaService()
	if err := r.Build(rsaClientKey); err != nil {
		hlog.Info(err)
		return errmsg.ServiceError
	}
	userMap[uid] = &_user{conn: service.conn, username: user.Username, rsa: r}
	publicKey, err := r.GetPublicKeyPemFormat()
	if err != nil {
		return errmsg.ServiceError
	}
	if err := service.conn.WriteMessage(websocket.TextMessage, []byte(publicKey)); err != nil {
		return errmsg.ServiceError
	}

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

func (service ChatService) SendMessage(content []byte) error {
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
	fromConn := userMap[from]
	switch toConn {
	case nil: // 离线
		{
			plainText, err := fromConn.rsa.Decode(content)
			if err != nil {
				return errmsg.ServiceError
			}
			if err := db.CreateMessage(from, to, string(plainText)); err != nil {
				return errmsg.RedisError
			}
		}
	default: // 在线
		{
			plainText, err := fromConn.rsa.Decode(content)
			if err != nil {
				return errmsg.ServiceError
			}
			ciphertext, err := toConn.rsa.Encode(userinfoAppend(plainText, from))
			if err != nil {
				return errmsg.ServiceError
			}
			err = toConn.conn.WriteMessage(websocket.BinaryMessage, ciphertext)
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
	list, err := db.PopMessage(uid)
	if err != nil {
		return errmsg.ServiceError
	}
	toConn := userMap[uid]
	for _, item := range *list {
		ciphertext, err := toConn.rsa.Encode(userinfoAppend([]byte(item.Content), item.FromUserId))
		if err != nil {
			return errmsg.ServiceError
		}
		err = service.conn.WriteMessage(websocket.BinaryMessage, ciphertext)
		if err != nil {
			return errmsg.ServiceError
		}
	}
	return nil
}
func userinfoAppend(rawText []byte, from string) []byte {
	return []byte(utils.ConvertTimestampToStringDefault(time.Now().Unix()) + ` [` + from + `]: ` + string(rawText))
}
