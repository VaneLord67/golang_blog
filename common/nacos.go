package common

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func FindService() (ip string, port uint64, err error) {
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
		return "", 0, err
	}

	// SelectOneHealthyInstance将会按加权随机轮询的负载均衡策略返回一个健康的实例
	// 实例必须满足的条件：health=true,enable=true and weight>0
	instance, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: "captcha",
		GroupName:   "base", // 默认值DEFAULT_GROUP
	})
	if err != nil {
		return "", 0, err
	}
	return instance.Ip, instance.Port, err
}
