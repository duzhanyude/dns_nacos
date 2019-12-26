package conf

import (
	"math/rand"
	"net"
)

var global = make(map[string][4]byte)
var up_DNS_Server = []net.IP{net.IP{114, 114, 114, 114}, net.IP{61, 128, 128, 68}, net.IP{8, 8, 8, 8}}

func init() {
	//初始化本地域名映射
	//SaveIPWithName("www.baidu.cmm.",[4]byte{172, 16, 254, 104})
	SaveIPWithName("www.baidu.comq.", [4]byte{172, 16, 254, 40})
}

func GetUPDNS() net.IP {
	size := len(up_DNS_Server)
	i := rand.Intn(size)
	return up_DNS_Server[i]
}

func GetIPWithName(name string) [4]byte {
	return global[name]
}
func DelIPWithName(name string) {
	delete(global, name)
}
func SaveIPWithName(name string, buff [4]byte) {
	global[name] = buff
}
func UpdateIPWithName(name string, buff [4]byte) {
	global[name] = buff
}
