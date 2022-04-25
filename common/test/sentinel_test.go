package test

import (
	"fmt"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"log"
	"testing"
)
import (
	sentinel "github.com/alibaba/sentinel-golang/api"
)

func TestSentinel(t *testing.T) {
	initSentinel()
}

func initSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              50,   // QPS
			StatIntervalInMs:       1000, // 统计周期1秒
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		// 请求被拒绝，在此处进行处理
	} else {
		// 请求允许通过，此处编写业务逻辑
		fmt.Println("hello world")
		// 务必保证业务结束后调用 Exit
		e.Exit()
	}
}
