package test

import (
	"common"
	"encoding/json"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
	"strconv"
	"testing"
)

func TestChangeQPSConfig(t *testing.T) {
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}
	ok, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:               "qps_middleware",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              50,   // QPS
			StatIntervalInMs:       1000, // 统计周期1秒
		},
	})
	fmt.Println(ok)
	if err != nil {
		fmt.Println(err)
		return
	}
	ok, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "qps_middleware",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              20,   // QPS
			StatIntervalInMs:       1000, // 统计周期1秒
		},
	})
	fmt.Println(ok)
	if err != nil {
		fmt.Println(err)
	}
}

func TestListenConfig(t *testing.T) {
	configClient := CreateConfigClient()
	err := configClient.ListenConfig(vo.ConfigParam{
		DataId: "Sentinel",
		Group:  "base",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
			var conf struct {
				QPS int
			}
			err := json.Unmarshal([]byte(data), &conf)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(conf)
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	select {}
}

func TestNacosConfig(t *testing.T) {
	configClient := CreateConfigClient()
	qps := 20
	success, err := configClient.PublishConfig(vo.ConfigParam{
		DataId:  "Sentinel",
		Group:   "base",
		Content: "{\"QPS\": " + strconv.Itoa(qps) + "}"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(success)
}

func TestGetNacosConfig(t *testing.T) {
	configClient := CreateConfigClient()
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "Sentinel",
		Group:  "base"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(content)
	var conf struct {
		QPS int
	}
	err = json.Unmarshal([]byte(content), &conf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf)
}

func CreateConfigClient() config_client.IConfigClient {
	var clientConfig = constant.ClientConfig{
		NamespaceId:         common.GetNacosConf().NamespaceId, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		LogLevel:            "warn",
		NotLoadCacheAtStart: true,
	}
	var serverConfigs = []constant.ServerConfig{
		{
			IpAddr:      common.GetNacosConf().Host,
			ContextPath: "/nacos",
			Port:        uint64(common.GetNacosConf().Port),
			Scheme:      "http",
		},
	}
	// 创建动态配置客户端的另一种方式 (推荐)
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	return configClient
}
