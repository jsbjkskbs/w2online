// Code generated by hertz generator.

package main

import (
	"work/biz/mw/jwt"
	webs "work/biz/router/websocket"
	cfgloader "work/pkg/utils/cfg_loader"
	"work/pkg/utils/dustman"
	"work/pkg/utils/syncman"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	cfgloader.Run()

	jwt.AccessTokenJwtInit()
	jwt.RefreshTokenJwtInit()

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
