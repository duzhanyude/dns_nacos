package main

import (
	"com.dnsnacos/conf"
	"com.dnsnacos/dns"
	"flag"
)

var (
	Nacos_IP        = flag.String("N_IP", "", "nacos ip address")
	Nacos_NameSpace = flag.String("N_Name", "", "nacos  namespace")
)

func main() {
	flag.Parse()
	conf.InitNacos(*Nacos_IP, *Nacos_NameSpace)
	dns.Listen()
}
