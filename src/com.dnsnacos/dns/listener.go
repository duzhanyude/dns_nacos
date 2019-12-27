package dns

import (
	"com.dnsnacos/conf"
	"fmt"
	"golang.org/x/net/dns/dnsmessage"
	"net"
	"strconv"
	"time"
)

type request struct {
	addr    *net.UDPAddr
	mesaage dnsmessage.Message
}
type response struct {
	mesaage dnsmessage.Message
}

var respMap = make(map[string]response)
var reqMap = make(map[string]request)

func Listen() {
	//定义定时器
	go timer()
	var m dnsmessage.Message
	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 53})
	defer conn.Close()
	for {
		buf := make([]byte, 512)
		_, addr, _ := conn.ReadFromUDP(buf)
		_ = m.Unpack(buf)
		if m.Header.Response {
			for _, question := range m.Questions {
				name := question.Name.String()
				respMap[name] = response{m}
				r := reqMap[name]
				packed, _ := m.Pack()
				conn.WriteToUDP(packed, r.addr)
				delete(reqMap, name)
			}
		}
		fmt.Println(len(respMap))
		if len(reqMap) > 1000 {
			reqMap = make(map[string]request)
		}
		if len(respMap) > 1000000 {
			respMap = make(map[string]response)
		}
		if len(m.Questions) > 0 {

			for _, question := range m.Questions {
				if question.Type == dnsmessage.TypeAAAA {
					m.Response = true
					m.Answers = []dnsmessage.Resource{}
					pack, _ := m.Pack()
					conn.WriteTo(pack, addr)
					continue
				}
				//本地配置的域名
				bytes := conf.GetIPWithName(question.Name.String())
				if bytes[0] != 0 {
					m.Answers = createDefinePackage(question)
					pack, _ := m.Pack()
					conn.WriteTo(pack, addr)
				}
				//判断是否存在本地缓存中
				r := respMap[question.Name.String()]
				if r.mesaage.ID != 0 {
					pack, _ := r.mesaage.Pack()
					conn.WriteTo(pack, addr)
				}

				reqMap[question.Name.String()] = request{addr, m}
				//向上游获取域名
				dnsIP := conf.GetUPDNS()
				fmt.Println(dnsIP)
				packed, _ := m.Pack()
				resolver := net.UDPAddr{IP: dnsIP, Port: 53}
				conn.WriteToUDP(packed, &resolver)
			}
		}

	}
}

//创建静态映射的域名资源
func createDefinePackage(question dnsmessage.Question) []dnsmessage.Resource {
	resource := dnsmessage.Resource{}
	resource.Header.Name = question.Name
	resource.Header.Type = question.Type
	resource.Header.Class = question.Class
	resource.Body = &dnsmessage.AResource{conf.GetIPWithName(question.Name.String())}
	resources := []dnsmessage.Resource{resource}
	return resources
}

func timer() {
	timer := time.NewTicker(1000 * time.Second)
	for {
		select {
		case <-timer.C:
			for k, v := range respMap {
				if len(v.mesaage.Answers) > 0 {
					u := v.mesaage.Answers[0].Header.TTL
					u = u - 1
					fmt.Println("ttl--" + strconv.Itoa(int(u)))
					if u <= 0 {
						delete(respMap, k)
					} else {
						v.mesaage.Answers[0].Header.TTL = u
						respMap[k] = v
					}
				}
			}
		}
	}
}
