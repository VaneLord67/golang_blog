package common

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/gin-gonic/gin"
	"log"
)

func QpsMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		e, b := sentinel.Entry("qps_middleware", sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			// 请求被拒绝，在此处进行处理
			FailCode(context, TOO_MANY_REQUESTS)
			context.Abort()
			return
		} else {
			// 请求允许通过，此处编写业务逻辑
			context.Next()
			// 务必保证业务结束后调用 Exit
			e.Exit()
		}
	}
}

func InitSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "qps_middleware",
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
}
