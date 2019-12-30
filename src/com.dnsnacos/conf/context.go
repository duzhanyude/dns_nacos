package conf

import (
	"com.dnsnacos/db"
	"encoding/json"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

var global = make(map[string][][4]byte)
var up_DNS_Server = []net.IP{net.IP{61, 128, 128, 68}}
var LOCAL_CONFIG_TABLE = "local_conf_table"
var GLOBAL_CONFIG = "global_conf_key"

func init() {
	//初始化本地域名映射
	buff := db.Get(LOCAL_CONFIG_TABLE, GLOBAL_CONFIG)
	if buff != "" {
		SaveLocalConfig(buff)
	}
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

func SaveLocalConfig(data string) {

	var dms domainObjs
	json.Unmarshal([]byte(data), &dms)

	//保存域名映射
	objs := dms.Domains
	for _, val := range objs {
		DelIPWithName(val.Domain)
		for _, ip := range val.Ips {
			var i_p [4]byte
			ip_s := strings.Split(ip, ".")
			for i, v := range ip_s {
				val, _ := strconv.Atoi(v)
				i_p[i] = byte(val)
			}
			SaveIPWithName(val.Domain+".", i_p)
		}
	}
	//保存dns
	dns := dms.Updns
	ClearDNS()
	for _, val := range dns {
		var ip [4]byte
		ips := strings.Split(val, ".")
		for i, v := range ips {
			val, _ := strconv.Atoi(v)
			ip[i] = byte(val)
		}
		SaveDNS(net.IP{ip[0], ip[1], ip[2], ip[3]})
	}
}
