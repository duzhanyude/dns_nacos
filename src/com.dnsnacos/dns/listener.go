package dns

import (
	"com.dnsnacos/conf"
	"golang.org/x/net/dns/dnsmessage"
	"net"
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
	var m dnsmessage.Message
	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 53})
	defer conn.Close()
	for {
		buf := make([]byte, 512)
		_, addr, _ := conn.ReadFromUDP(buf)
		_ = m.Unpack(buf)
		if m.Header.Response {
			name := m.Questions[0].Name.String()
			//respMap[name]=response{m}
			r := reqMap[name]
			packed, _ := m.Pack()
			conn.WriteToUDP(packed, r.addr)
			delete(reqMap, name)
		}

		if len(m.Questions) > 0 {
			question := m.Questions[0]
			//判断是否存在本地缓存中
			/*r := respMap[question.Name.String()]
			if r.mesaage.ID!=0{
				pack,_ :=r.mesaage.Pack()
				conn.WriteTo(pack,addr)
			}*/
			//本地配置的域名
			bytes := conf.GetIPWithName(question.Name.String())
			if bytes[0] != 0 {
				m.Answers = createDefinePackage(question)
				pack, _ := m.Pack()
				conn.WriteTo(pack, addr)
			}

			reqMap[question.Name.String()] = request{addr, m}
			//向上游获取域名
			packed, _ := m.Pack()
			resolver := net.UDPAddr{IP: conf.GetUPDNS(), Port: 53}
			_, _ = conn.WriteToUDP(packed, &resolver)
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
