package webs

import (
	wschat "work/biz/handler/ws_chat"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func register(h *server.Hertz) {
	h.GET(`/`, append(_homeMW(), wschat.Handler)...)
}
