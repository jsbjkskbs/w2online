package webs

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

func WebsocketRegister(h *server.Hertz) {
	register(h)
}
