package common

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
	"sync"
)

type NacosConf struct {
	Host        string
	Port        int64
	NamespaceId string
}

var nacosConf *NacosConf

func initNacosConf() {
	conf := ReadYaml()
	nacosConf = &NacosConf{
		Host:        conf.Nacos.Host,
		Port:        conf.Nacos.Port,
		NamespaceId: conf.Nacos.NamespaceId,
	}
}

var onceNacos = sync.Once{}

func GetNacosConf() *NacosConf {
	onceNacos.Do(initNacosConf)
	return nacosConf
}

var onceNacosClient = sync.Once{}
var nacosNamingClient naming_client.INamingClient

func createNNC() {
	var clientConfig = constant.ClientConfig{
		NamespaceId:         GetNacosConf().NamespaceId, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		LogLevel:            "warn",
		NotLoadCacheAtStart: true,
	}
	var serverConfigs = []constant.ServerConfig{
		{
			IpAddr:      GetNacosConf().Host,
			ContextPath: "/nacos",
			Port:        uint64(GetNacosConf().Port),
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
	nacosNamingClient = namingClient
}
func GetNNC() naming_client.INamingClient {
	onceNacosClient.Do(createNNC)
	return nacosNamingClient
}

func CreateService(serviceName, groupName, ip string, port uint64) bool {
	if groupName == "" {
		groupName = "DEFAULT_GROUP"
	}
	success, err := GetNNC().RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		GroupName:   groupName, // 默认值DEFAULT_GROUP
	})
	if err != nil {
		log.Fatal(err)
	}
	return success
}

func FindService(serviceName string, groupName string) (ip string, port uint64, err error) {
	if groupName == "" {
		groupName = "DEFAULT_GROUP"
	}
	// SelectOneHealthyInstance将会按加权随机轮询的负载均衡策略返回一个健康的实例
	// 实例必须满足的条件：health=true,enable=true and weight>0
	instance, err := GetNNC().SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
		GroupName:   groupName, // 默认值DEFAULT_GROUP
	})
	if err != nil {
		return "", 0, err
	}
	return instance.Ip, instance.Port, err
}
