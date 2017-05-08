package netservice

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"ActivedRouter/global"
	"ActivedRouter/system"

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
		info = system.SysInfo(global.RunMode, "Router", "")
		global.SetRouterInfo(info)
	}
	self.WriteJsonString(w, info)
}

//反向代理信息
func (self *Http) ProxyInfos(w http.ResponseWriter, r *http.Request, prms httprouter.Params) {
	hostInfos := ProxyHandler.GetDomainHostList(prms.ByName("domain"))
	self.WriteJsonInterface(w, hostInfos)
}

func (self *Http) DomainInfos(w http.ResponseWriter, r *http.Request, prms httprouter.Params) {
	keysArray := ProxyHandler.DomainInfos()
	self.WriteJsonInterface(w, keysArray)
}

func (self *Http) AddDomain(w http.ResponseWriter, r *http.Request, prms httprouter.Params) {
	if ProxyHandler.AddDomainConfig(prms.ByName("domain")) {
		self.WriteJsonString(w, `{"status":"1"}`)
	} else {
		self.WriteJsonString(w, `{"status":"0"}`)
	}
}
func (self *Http) DelDomain(w http.ResponseWriter, r *http.Request, prms httprouter.Params) {
	if ProxyHandler.DeleteDomainConig(prms.ByName("domain")) {
		self.WriteJsonString(w, `{"status":"1"}`)
	} else {
		self.WriteJsonString(w, `{"status":"0"}`)
	}
}

func (self *Http) AddProxyClient(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	domain := r.Form.Get("domain")
	host := r.Form.Get("host")
	port := r.Form.Get("port")
	if ret := ProxyHandler.AddProxyClient(domain, host, port); ret == -1 {
		self.WriteJsonString(w, `{"status":0,"data":{"code":-1}}`)
	} else if ret == 0 {
		self.WriteJsonString(w, `{"status":0,"data":{"code":0}}`)
	} else {
		self.WriteJsonString(w, `{"status":1}`)
	}
}
func (self *Http) DeleteProxyClient(w http.ResponseWriter, r *http.Request, prms httprouter.Params) {
	r.ParseForm()
	domain := r.Form.Get("domain")
	host := r.Form.Get("host")
	port := r.Form.Get("port")
	if ret := ProxyHandler.DeleteProxyClient(domain, host, port); !ret {
		self.WriteJsonString(w, `{"status":0}`)
	} else {
		self.WriteJsonString(w, `{"status":1}`)
	}
}

func (self *Http) UpdateProxyClient(w http.ResponseWriter, r *http.Request, prms httprouter.Params) {
	w.Write([]byte("hello,world"))
}

//创建http服务
func NewHttp(host, port string) *Http {
	return &Http{Host: host, Port: port}
}

//run router server
func (self *Http) Run() {
	log.Printf("开始启动http服务,%s:%s........\n", self.Host, self.Port)
	router := httprouter.New()
	//统计信息
	router.GET("/clientinfos", self.ClientInfos)
	router.GET("/statistics", self.Statistics)
	router.GET("/routerinfo", self.RouterInfo)
	router.GET("/activeclients", self.ActiveClientInfos)
	router.GET("/bestclients", self.ActiveClientInfos)
	//反向代理配置文件读写
	router.GET("/proxyinfos/:domain", self.ProxyInfos)
	router.GET("/adddomain/:domain", self.AddDomain)
	router.GET("/deldomain/:domain", self.DelDomain)
	router.GET("/addproxyclient", self.AddProxyClient)
	router.GET("/delproxyclient", self.DeleteProxyClient)
	router.GET("/updateproxyclient", self.UpdateProxyClient)
	router.GET("/domaininfos", self.DomainInfos)
	//静态文件路由
	router.GET("/", self.Index)
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.ServeFiles("/website/*filepath", http.Dir("website"))
	log.Fatal(http.ListenAndServe(self.Host+":"+self.Port, router))
}
