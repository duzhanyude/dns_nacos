package main

import (
	"com.dnsnacos/conf"
	_ "com.dnsnacos/db"
	"com.dnsnacos/dns"
	"flag"
	"os"
)

var (
	Nacos_IP        = *flag.String("N_IP", os.Getenv("N_IP"), "nacos ip address")
	Nacos_NameSpace = *flag.String("N_Name", os.Getenv("N_Name"), "nacos  namespace")
)

func main() {
	flag.Parse()
	conf.InitNacos(Nacos_IP, Nacos_NameSpace)
	dns.Listen()
}
