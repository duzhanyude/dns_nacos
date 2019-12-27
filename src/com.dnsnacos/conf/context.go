package conf

import (
	"math/rand"
	"net"
)

var global = make(map[string][][4]byte)
var up_DNS_Server = []net.IP{net.IP{61, 128, 128, 68}}

func init() {
	//初始化本地域名映射
}
func SaveDNS(ip net.IP) {
	up_DNS_Server = append(up_DNS_Server, ip)
}
func ClearDNS() {
	up_DNS_Server = []net.IP{net.IP{61, 128, 128, 68}}
}

func GetUPDNS() net.IP {
	size := len(up_DNS_Server)
	i := rand.Intn(size)
	return up_DNS_Server[i]
}
func GetIPWithName(name string) [][4]byte {
	return global[name]
}
func DelIPWithName(name string) {
	delete(global, name)
}
func SaveIPWithName(name string, buff [4]byte) {
	global[name] = append(global[name], buff)
}
