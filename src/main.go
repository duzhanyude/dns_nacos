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
	Nacos_IP        string
	Nacos_NameSpace string
)

func main() {
	flag.StringVar(&Nacos_IP, "nacos_ip", os.Getenv("nacos_ip"), "nacos ip address")
	flag.StringVar(&Nacos_NameSpace, "nacos_name", os.Getenv("nacos_name"), "nacos  namespace")
	flag.Parse()
	conf.InitNacos(Nacos_IP, Nacos_NameSpace)
	dns.Listen()
}
