package sentinel

import (
	"github.com/alibaba/sentinel-golang/core/flow"
)

var Rules map[string]interface{}

func update(flowrules *[]*flow.Rule) {
	for resource := range Rules {
		*flowrules = append(*flowrules, buildRule(resource, Rules[resource].(map[string]interface{})))
	}
}

func buildRule(resource string, rule map[string]interface{}) *flow.Rule {
	return &flow.Rule{
		Resource:               resource,
		Threshold:              rule["threshold"].(float64),
		StatIntervalInMs:       uint32(rule["statintervalinms"].(float64)),
		TokenCalculateStrategy: flow.Direct,
		ControlBehavior:        flow.Reject,
	}
}
