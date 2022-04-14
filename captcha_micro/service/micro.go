package service

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
)

// CreateService 服务注册
func CreateService(port uint64) bool {
	var clientConfig = constant.ClientConfig{
		NamespaceId:         "f7f6dce8-6264-46ef-8561-beeb8026f213", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		LogLevel:            "warn",
		NotLoadCacheAtStart: true,
	}
	var serverConfigs = []constant.ServerConfig{
		{
			IpAddr:      "localhost",
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
	}
	// 创建服务发现客户端的另一种方式 (推荐)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "localhost",
		Port:        port,
		ServiceName: "captcha",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		GroupName:   "base", // 默认值DEFAULT_GROUP
	})
	if err != nil {
		log.Fatal(err)
	}
	return success
}
