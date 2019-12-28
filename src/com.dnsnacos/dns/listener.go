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
	conn    *net.UDPConn
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
	reqMap      = make(map[string]request)
	respMap     = make(map[string]response)
	reqLock     sync.Mutex
	responsLock sync.Mutex
	message     = make(chan dnsmessage.Message)
)

func Listen() {
	//定义定时器
	go timer()
	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 53})
	defer conn.Close()
	fmt.Println("开启代理成功")
	req := make(chan request)

	go handler(req)
	go sendRemote(message)
	for {
		buf := make([]byte, 512)
		_, addr, _ := conn.ReadFromUDP(buf)
		var m dnsmessage.Message
		_ = m.Unpack(buf)
		request_r := request{addr: addr, mesaage: m, conn: conn}
		//fmt.Println(request_r.mesaage.GoString())
		req <- request_r
	}
}

func handler(req chan request) {

	for {
		request_t := <-req
		/*if m.Header.Response {
			handlerResponse(conn, m)
			return
		}*/
		m := request_t.mesaage
		addr := request_t.addr

		if len(m.Questions) > 0 {
			if m.Questions[0].Type == dnsmessage.TypeAAAA {
				m.Answers = []dnsmessage.Resource{}
				m.Response = true
				//sendUDPP(addr,m)
				packed, _ := m.Pack()
				request_t.conn.WriteTo(packed, addr)
				continue
			}

			if routeLocalDNS(request_t) {
				routeRemoteDNS(request_t)
			}
		}

		//
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
func handlerResponse(message dnsmessage.Message) {
	if !message.Response {
		return
	}
	fmt.Println(message.GoString())
	msg := message
	reqLock.Lock()
	request := reqMap[strconv.Itoa(int(msg.ID))]
	if request.addr != nil {
		//sendUDPP(request.addr,request.mesaage)
		packed, _ := message.Pack()
		request.conn.WriteTo(packed, request.addr)
		delete(reqMap, strconv.Itoa(int(msg.ID)))
	}
	reqLock.Unlock()
	//缓存问题回答
	question := msg.Questions[0]

	//respMap := getRespMap()
	/*if respMap.RespMap == nil {
		respMap = responseMap{}
		respMap.RespMap = make(map[string]response)
	}*/
	responsLock.Lock()
	respMap[question.Name.String()+question.Type.String()] = response{msg}
	responsLock.Unlock()
	//saveRespMap(respMap)
}

func clearCache() {
	//fmt.Println(len(reqMap))
	if len(reqMap) > 500 {
		reqLock.Lock()
		reqMap = make(map[string]request)
		reqLock.Unlock()
	}

	//respMap := getRespMap()
	fmt.Println(len(respMap))
	if len(respMap) > 1000000 {
		responsLock.Lock()
		respMap = make(map[string]response)
		responsLock.Unlock()
	}
	//saveRespMap(respMap)

}

//路由本地
func routeLocalDNS(req request) bool {
	var m = req.mesaage
	addr := req.addr
	question := m.Questions[0]
	//本地配置的域名
	bytes := conf.GetIPWithName(question.Name.String())
	if len(bytes) > 0 {
		m.Answers = createDefinePackage(question)
		m.Response = true
		m.RecursionDesired = true
		packed, _ := m.Pack()
		req.conn.WriteTo(packed, addr)
		//fmt.Println(m.GoString())
		//sendUDPP(addr,m)
		return false
	}
	return true
}

//向上游获取域名
func routeRemoteDNS(req request) {
	var m = req.mesaage
	var addr = req.addr
	question := m.Questions[0]
	//获取本地缓存
	//respMap := getRespMap()
	//if respMap != nil {
	res := respMap[question.Name.String()+question.Type.String()]
	if len(res.Mesaage.Questions) > 0 && len(res.Mesaage.Answers) > 0 {
		m.Answers = res.Mesaage.Answers
		m.Response = true
		//sendUDPP(addr,m)
		packed, _ := m.Pack()
		req.conn.WriteTo(packed, addr)
		return
	} else {
		message <- m
	}
	reqLock.Lock()
	reqMap[strconv.Itoa(int(m.ID))] = request{addr, m, req.conn}
	reqLock.Unlock()
}
func sendRemote(msg chan dnsmessage.Message) {
	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 10053})
	defer conn.Close()
	go receiveMessage(conn)
	for {
		message := <-msg
		dnsIP := conf.GetUPDNS()
		fmt.Println(dnsIP)
		//fmt.Println(message.GoString())
		packed, _ := message.Pack()
		resolver := net.UDPAddr{IP: dnsIP, Port: 53}
		conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
		conn.WriteTo(packed, &resolver)

		//connect, _ := net.DialUDP("udp", nil, &resolver)
		//connect.Write(packed)
		//buf := make([]byte, 512)
		/*connect.Read(buf)
		var m  dnsmessage.Message
		_ = m.Unpack(buf)
			c*/

	}
}
func sendUDPP(addr *net.UDPAddr, mes dnsmessage.Message) {
	packed, _ := mes.Pack()
	connect, _ := net.DialUDP("udp", nil, addr)
	connect.Write(packed)
}
func receiveMessage(conn *net.UDPConn) {
	for {
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		buf := make([]byte, 512)
		conn.Read(buf)
		var m dnsmessage.Message
		_ = m.Unpack(buf)
		handlerResponse(m)
	}
}

//设置定时器
func timer() {
	timer := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-timer.C:
			go handlerLocalRemoteCache()
			clearCache()
		}
	}
}

//处理远程本地缓存
func handlerLocalRemoteCache() {

	//respMap := getRespMap()
	/*if respMap.RespMap == nil {
		return
	}*/
	for k, v := range respMap {
		if len(v.Mesaage.Answers) > 0 {
			u := v.Mesaage.Answers[0].Header.TTL
			u = u - 1
			//fmt.Println("ttl-"+strconv.Itoa(int(u)))
			responsLock.Lock()
			if u <= 0 {
				delete(respMap, k)
			} else {
				v.Mesaage.Answers[0].Header.TTL = u
				respMap[k] = v
			}
			responsLock.Unlock()
		}
	}
	//saveRespMap(respMap)
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
