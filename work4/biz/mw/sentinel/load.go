package sentinel

import (
	aliSentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Load() {
	if err := aliSentinel.InitDefault(); err != nil {
		hlog.Warn("sentinel did not start successfully")
		return
	}

	var newRule []*flow.Rule
	update(&newRule)

	if _, err := flow.LoadRules(newRule); err != nil {
		hlog.Warn("sentinel did not start successfully")
		return
	}
	hlog.Info("sentinel prepared successfully")
}
