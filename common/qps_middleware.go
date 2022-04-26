package common

import (
	"encoding/json"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
)

func DynamicQPS() {
	configClient := CreateConfigClient()
	err := configClient.ListenConfig(vo.ConfigParam{
		DataId: "Sentinel",
		Group:  "base",
		OnChange: func(namespace, group, dataId, data string) {
			var conf struct {
				QPS int
			}
			err := json.Unmarshal([]byte(data), &conf)
			if err != nil {
				log.Println("动态加载Nacos配置失败,放弃本次配置切换")
				return
			}
			ok, err := flow.LoadRules([]*flow.Rule{
				{
					Resource:               "qps_middleware",
					TokenCalculateStrategy: flow.Direct,
					ControlBehavior:        flow.Reject,
					Threshold:              float64(conf.QPS), // QPS
					StatIntervalInMs:       1000,              // 统计周期1秒
				},
			})
			if !ok {
				log.Println("规则相同,放弃本次配置切换")
				return
			}
			log.Println("切换配置为 QPS = ", conf.QPS)
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

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
