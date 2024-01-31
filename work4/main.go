// Code generated by hertz generator.

package main

import (
	"work/biz/dal"
	"work/biz/mw/elasticsearch"
	"work/biz/mw/jwt"
	webs "work/biz/router/websocket"
	qiniuyunoss "work/pkg/qiniuyun_oss"
	"work/pkg/utils/dustman"
	"work/pkg/utils/syncman"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	elasticsearch.Init()
	dal.Init()
	jwt.AccessTokenJwtInit()
	jwt.RefreshTokenJwtInit()
	qiniuyunoss.OssInit()
	dustman.NewRedisDustman().Run()
	dustman.NewFileDustman().Run()
	syncman.NewVideoSyncman().Run()
	syncman.NewCommentSyncman().Run()
	syncman.NewRelationSyncman().Run()
	h := server.Default(server.WithHostPorts(`:10001`))
	ws := server.Default(server.WithHostPorts(`:10000`))
	ws.NoHijackConnPool = true
	register(h)
	webs.WebsocketRegister(ws)

	go ws.Spin()
	h.Spin()
}
