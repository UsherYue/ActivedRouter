//对于客户端 服务器以及http服务的网络封装
package netservice

import (
	"ActivedRouter/global"
	"ActivedRouter/system"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Http struct {
	Host string
	Port string
}

//ClientInfos
//所有服务器 活跃和非活跃
func (self *Http) ClientInfos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := global.GHostInfoTable.HostsInfo.GetStorage().GetData()
	self.WriteJsonInterface(w, data)
}

//statistics
func (self *Http) Statistics(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := global.GProxyHttpStatistics.GetStatisticsList()
	self.WriteJsonInterface(w, data)
}

//输出json
func (self *Http) WriteJsonString(w http.ResponseWriter, str string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(str))
}

//输出interface
func (self *Http) WriteJsonInterface(w http.ResponseWriter, data interface{}) {
	bts, _ := json.MarshalIndent(data, "", " ")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bts))
}

//Active ClientInfos
//活跃列表
func (self *Http) ActiveClientInfos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := global.GHostInfoTable.ActiveHostList.ActiveHostInfo.GetStorage().GetData()
	self.WriteJsonInterface(w, data)
}

//BEST  ClientInfos
//权重最高的服务器
func (self *Http) BestClients(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := global.GHostInfoTable.ActiveHostWeightList.Front().Value
	self.WriteJsonInterface(w, data)
}

//index redirect to static
func (self *Http) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, "/static", 302)
}

//路由服务器的信息
func (self *Http) RouterInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	info := global.RouterInfo()
	if info == "" {
		info = system.SysInfo("Router", "")
		global.SetRouterInfo(info)
	}
	self.WriteJsonString(w, info)
}

//创建http服务
func NewHttp(host, port string) *Http {
	return &Http{Host: host, Port: port}
}

//run router server
func (self *Http) Run() {
	log.Printf("开始启动http服务,%s:%s........\n", self.Host, self.Port)
	router := httprouter.New()
	router.GET("/clientinfos", self.ClientInfos)
	router.GET("/statistics", self.Statistics)
	router.GET("/routerinfo", self.RouterInfo)
	router.GET("/activeclients", self.ActiveClientInfos)
	router.GET("/bestclients", self.ActiveClientInfos)
	router.GET("/", self.Index)
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.ServeFiles("/website/*filepath", http.Dir("website"))
	log.Fatal(http.ListenAndServe(self.Host+":"+self.Port, router))
}
