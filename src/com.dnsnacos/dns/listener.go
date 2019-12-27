package dns

import (
	"com.dnsnacos/conf"
	"com.dnsnacos/db"
	"encoding/json"
	"fmt"
	"golang.org/x/net/dns/dnsmessage"
	"net"
	"strconv"
	"sync"
	"time"
)

type request struct {
	addr    *net.UDPAddr
	mesaage dnsmessage.Message
}
type response struct {
	Mesaage dnsmessage.Message
}
type responseMap struct {
	RespMap map[string]response
}

var (
	TABLE_NAME        = "dns"
	CONSTANT_REQUEST  = "reqMap"
	CONSTANT_RESPONSE = "respMap"
)
var (
	m      dnsmessage.Message
	reqMap = make(map[string]request)
	lock   sync.Mutex
)

func Listen() {
	//定义定时器
	go timer()
	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 53})
	defer conn.Close()

	for {
		buf := make([]byte, 512)
		_, addr, _ := conn.ReadFromUDP(buf)
		_ = m.Unpack(buf)
		if m.Header.Response {
			go handlerResponse(conn, m)
			continue
		}

		go clearCache()
		if len(m.Questions) > 0 {

			routeLocalDNS(conn, addr)
			routeRemoteDNS(conn, addr)
		}
	}
}

//创建静态映射的域名资源
func createDefinePackage(question dnsmessage.Question) []dnsmessage.Resource {
	resources := []dnsmessage.Resource{}
	ips := conf.GetIPWithName(question.Name.String())
	for _, ip := range ips {
		resource := dnsmessage.Resource{}
		resource.Header.Name = question.Name
		resource.Header.Type = question.Type
		resource.Header.Class = question.Class
		resource.Body = &dnsmessage.AResource{ip}
		resources = append(resources, resource)
	}
	return resources
}

//处理上游返回
func handlerResponse(conn *net.UDPConn, message dnsmessage.Message) {
	msg := message
	connection := conn
	packed, _ := msg.Pack()
	lock.Lock()
	request := reqMap[strconv.Itoa(int(msg.ID))]
	connection.WriteToUDP(packed, request.addr)
	delete(reqMap, strconv.Itoa(int(msg.ID)))
	lock.Unlock()
	//缓存问题回答
	question := msg.Questions[0]

	lock.Lock()
	respMap := getRespMap()
	if respMap.RespMap == nil {
		respMap = responseMap{}
		respMap.RespMap = make(map[string]response)
	}
	respMap.RespMap[question.GoString()] = response{msg}
	saveRespMap(respMap)
	lock.Unlock()
}

func clearCache() {
	if len(reqMap) > 500 {
		reqMap = make(map[string]request)
	}

	lock.Lock()
	respMap := getRespMap()
	fmt.Println(len(respMap.RespMap))
	if len(respMap.RespMap) > 1000000 {
		respMap.RespMap = make(map[string]response)
	}
	saveRespMap(respMap)
	lock.Unlock()

}

//路由本地
func routeLocalDNS(conn *net.UDPConn, addr *net.UDPAddr) {
	question := m.Questions[0]
	//本地配置的域名
	bytes := conf.GetIPWithName(question.Name.String())
	if len(bytes) > 0 {
		m.Answers = createDefinePackage(question)
		pack, _ := m.Pack()
		conn.WriteTo(pack, addr)
	}
}

//向上游获取域名
func routeRemoteDNS(conn *net.UDPConn, addr *net.UDPAddr) {

	question := m.Questions[0]
	//获取本地缓存
	lock.Lock()
	respMap := getRespMap()
	if respMap.RespMap != nil {
		r := respMap.RespMap[question.GoString()]
		if r.Mesaage.ID != 0 {
			r.Mesaage.ID = m.ID
			pack, _ := r.Mesaage.Pack()
			conn.WriteTo(pack, addr)
		}
	}
	reqMap[strconv.Itoa(int(m.ID))] = request{addr, m}
	lock.Unlock()
	dnsIP := conf.GetUPDNS()
	fmt.Println(dnsIP)
	packed, _ := m.Pack()
	resolver := net.UDPAddr{IP: dnsIP, Port: 53}
	conn.WriteToUDP(packed, &resolver)
}

//设置定时器
func timer() {
	timer := time.NewTicker(1000 * time.Second)
	for {
		select {
		case <-timer.C:
			go handlerLocalRemoteCache()
		}
	}
}

//处理远程本地缓存
func handlerLocalRemoteCache() {

	respMap := getRespMap()
	if respMap.RespMap == nil {
		return
	}
	lock.Lock()
	for k, v := range respMap.RespMap {
		if len(v.Mesaage.Answers) > 0 {
			u := v.Mesaage.Answers[0].Header.TTL
			u = u - 1
			//fmt.Println("ttl-"+strconv.Itoa(int(u)))
			if u <= 0 {
				delete(respMap.RespMap, k)
			} else {
				v.Mesaage.Answers[0].Header.TTL = u
				respMap.RespMap[k] = v
			}
		}
	}
	saveRespMap(respMap)
	lock.Unlock()
}
func getRespMap() responseMap {
	var respMap responseMap
	json.Unmarshal([]byte(db.Get(TABLE_NAME, CONSTANT_RESPONSE)), &respMap)
	return respMap
}
func saveRespMap(respMap responseMap) {
	b, _ := json.Marshal(respMap)
	db.Save(TABLE_NAME, CONSTANT_RESPONSE, string(b))
}
