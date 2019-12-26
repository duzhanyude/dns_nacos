package conf

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"strconv"
	"strings"
)

var configClient config_client.IConfigClient
var dataId = "domain.json"
var group = "domain"

type domainObj struct {
	Domain string `json:"domain"`
	Ip     string `json:"ip"`
}
type domainObjs struct {
	Domains []domainObj `json:"domains"`
}

func InitNacos(ip string, nameSpace string) {
	// 从控制台命名空间管理的"命名空间详情"中拷贝 End Point、命名空间 ID
	var endpoint = ip
	var namespaceId = nameSpace

	// 推荐使用 RAM 用户的 accessKey、secretKey
	//var accessKey = "${accessKey}"
	//var secretKey = "${secretKey}"

	clientConfig := constant.ClientConfig{
		//
		NamespaceId: namespaceId,
		//AccessKey:      accessKey,
		//SecretKey:      secretKey,
		TimeoutMs:      5 * 1000,
		ListenInterval: 30 * 1000,
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
			Content: "{\"domains\":[{\"domain\":\"localhost.com\",\"ip\":\"127.0.0.1\"}]}",
		})
	}
	configClient.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
			var dms domainObjs
			json.Unmarshal([]byte(data), &dms)
			objs := dms.Domains
			for _, val := range objs {
				var ip [4]byte
				ips := strings.Split(val.Ip, ".")
				for i, v := range ips {
					val, _ := strconv.Atoi(v)
					ip[i] = byte(val)
				}
				SaveIPWithName(val.Domain+".", ip)
			}
		},
	})
}
