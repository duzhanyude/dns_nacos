package web

import (
	"com.dnsnacos/dns"
	"fmt"
	//_ "github.com/mkevac/debugcharts"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
)

func printCacheStatus(w http.ResponseWriter, r *http.Request) {
	/*r.ParseForm() //解析参数，默认是不会解析的
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}*/
	fmt.Fprintf(w, createResponse()) //这个写入到w的是输出到客户端的
}
func createResponse() string {
	content := "cache total:" + dns.GetCacheSize() + "\r\n"
	content += "up request 1min: " + strconv.FormatFloat(dns.Remote_m.Rate1(), 'f', 2, 64) + "\r\n"
	content += "up response 1min: " + strconv.FormatFloat(dns.Remote_r_m.Rate1(), 'f', 2, 64) + "\r\n"
	content += "up request 5min: " + strconv.FormatFloat(dns.Remote_m.Rate5(), 'f', 2, 64) + "\r\n"
	content += "up response 5min: " + strconv.FormatFloat(dns.Remote_r_m.Rate5(), 'f', 2, 64) + "\r\n"
	content += "up request 15min: " + strconv.FormatFloat(dns.Remote_m.Rate15(), 'f', 2, 64) + "\r\n"
	content += "up response 15min: " + strconv.FormatFloat(dns.Remote_r_m.Rate15(), 'f', 2, 64) + "\r\n"
	return content
}
func initWeb() {

	//http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", printCacheStatus)    //设置访问的路由
	err := http.ListenAndServe(":10053", nil) //设置监听的端口
	if err != nil {
		log.Println("ListenAndServe: ", err)
	}
}

func init() {
	go initWeb()
}
