package conf

import (
	"com.dnsnacos/db"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var configClient config_client.IConfigClient
var dataId = "domain.json"
var group = "domain"

type domainObj struct {
	Domain string   `json:"domain"`
	Ips    []string `json:"ips"`
}
type domainObjs struct {
	Domains []domainObj `json:"domains"`
	Updns   []string    `json:"updns"`
}

func InitNacos(ip string, nameSpace string) {
	// 从控制台命名空间管理的"命名空间详情"中拷贝 End Point、命名空间 ID
	var endpoint = ip
	var namespaceId = nameSpace
	//var endpoint = "172.16.1.81"
	//var namespaceId = "8dfcdd4a-e7fb-4798-beb0-0c3cb84dd13e"

	// 推荐使用 RAM 用户的 accessKey、secretKey
	//var accessKey = "${accessKey}"
	//var secretKey = "${secretKey}"
	clientConfig := constant.ClientConfig{
		//
		NamespaceId: namespaceId,
		//AccessKey:      accessKey,
		//SecretKey:      secretKey,
		TimeoutMs:      5 * 1000,
		ListenInterval: 10 * 1000,
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      endpoint,
			ContextPath: "/nacos",
			Port:        8848,
		},
	}

	configClient, _ = clients.CreateConfigClient(map[string]interface{}{
		"clientConfig":  clientConfig,
		"serverConfigs": serverConfigs,
	})

	// 获取配置
	_, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		configClient.PublishConfig(vo.ConfigParam{
			DataId:  dataId,
			Group:   group,
			AppName: "dns_nacos",
			Content: "{\"domains\":[{\"domain\":\"localhost.com\",\"ips\":[\"127.0.0.1\"]}],\"updns\":[\"114.114.114.114\"]}",
		})
	}
	configClient.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			//fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
			SaveLocalConfig(data)
			db.Save(LOCAL_CONFIG_TABLE, GLOBAL_CONFIG, data)
		},
	})
}
