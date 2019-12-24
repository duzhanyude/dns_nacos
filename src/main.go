package main

import (
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

var reqMap = make(map[string]request)
var respMap = make(map[string]response)

var global = make(map[string][4]byte)

func main() {
	//处理本地域名映射
	global["www.baidu.cmm."] = [4]byte{192, 168, 1, 8}
	global["www.baidu.comq."] = [4]byte{192, 168, 1, 1}

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
			bytes := global[question.Name.String()]
			if bytes[0] != 0 {
				m.Answers = createDefinePackage(question)
				pack, _ := m.Pack()
				conn.WriteTo(pack, addr)
			}

			reqMap[question.Name.String()] = request{addr, m}
			//向上游获取域名
			packed, _ := m.Pack()
			resolver := net.UDPAddr{IP: net.IP{114, 114, 114, 114}, Port: 53}
			_, _ = conn.WriteToUDP(packed, &resolver)
		}

	}
}
func createDefinePackage(question dnsmessage.Question) []dnsmessage.Resource {
	resource := dnsmessage.Resource{}
	resource.Header.Name = question.Name
	resource.Header.Type = question.Type
	resource.Header.Class = question.Class
	resource.Body = &dnsmessage.AResource{global[question.Name.String()]}
	resources := []dnsmessage.Resource{resource}
	return resources
}
