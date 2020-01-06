package main

import (
	"com.dnsnacos/conf"
	_ "com.dnsnacos/db"
	"com.dnsnacos/dns"
	_ "com.dnsnacos/web"
	"flag"
	"os"
)

var (
	Nacos_IP        = *flag.String("nacos_ip", os.Getenv("nacos_ip"), "nacos ip address")
	Nacos_NameSpace = *flag.String("nacos_name", os.Getenv("nacos_name"), "nacos  namespace")
)

func main() {
	flag.Parse()
	conf.InitNacos(Nacos_IP, Nacos_NameSpace)
	dns.Listen()
}
