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
//All servers  info include active and inactive
func (self *Http) ClientInfos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := global.GHostInfoTable.HostsInfo.GetStorage().GetData()
	self.WriteJsonInterface(w, data)
}

//statistics
func (self *Http) Statistics(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := global.GProxyHttpStatistics.GetStatisticsList()
	self.WriteJsonInterface(w, data)
}

func (self *Http) WriteJsonString(w http.ResponseWriter, str string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(str))
}

func (self *Http) WriteJsonInterface(w http.ResponseWriter, data interface{}) {
	bts, _ := json.MarshalIndent(data, "", " ")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bts))
}

//Active ClientInfos
func (self *Http) ActiveClientInfos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := global.GHostInfoTable.ActiveHostList.ActiveHostInfo.GetStorage().GetData()
	self.WriteJsonInterface(w, data)
}

//The highest weight of the server
func (self *Http) BestClients(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := global.GHostInfoTable.ActiveHostWeightList.Front().Value
	self.WriteJsonInterface(w, data)
}

//index redirect to static
func (self *Http) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, "/static", 302)
}

//router server info
func (self *Http) RouterInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	info := global.RouterInfo()
	if info == "" {
		info = system.SysInfo(global.RunMode, "ActivedRouter")
		global.SetRouterSysInfo(info)
	}
	self.WriteJsonString(w, info)
}

//Reverse Proxy infp
func (self *Http) ProxyInfos(w http.ResponseWriter, r *http.Request, prms httprouter.Params) {
	hostInfos := DefaultHttpReverseProxy.GetDomainHostList(prms.ByName("domain"))
	self.WriteJsonInterface(w, hostInfos)
}

func (self *Http) DomainInfos(w http.ResponseWriter, r *http.Request, prms httprouter.Params) {
	keysArray := DefaultHttpReverseProxy.DomainInfos()
	self.WriteJsonInterface(w, keysArray)
}

func (self *Http) AddDomain(w http.ResponseWriter, r *http.Request, prms httprouter.Params) {
	if DefaultHttpReverseProxy.AddDomainConfig(prms.ByName("domain")) {
		self.WriteJsonString(w, `{"status":"1"}`)
	} else {
		self.WriteJsonString(w, `{"status":"0"}`)
	}
}
func (self *Http) DelDomain(w http.ResponseWriter, r *http.Request, prms httprouter.Params) {
	if DefaultHttpReverseProxy.DeleteDomainConig(prms.ByName("domain")) {
		self.WriteJsonString(w, `{"status":"1"}`)
	} else {
		self.WriteJsonString(w, `{"status":"0"}`)
	}
}

//change  proxy domain
func (self *Http) UpdateDomain(w http.ResponseWriter, r *http.Request, prms httprouter.Params) {
	r.ParseForm()
	preDomain := r.Form.Get("predomain")
	updateDomain := r.Form.Get("updatedomain")
	if DefaultHttpReverseProxy.UpdateDomain(preDomain, updateDomain, "on", "on") {
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
	if ret := DefaultHttpReverseProxy.AddProxyClient(domain, host, port, "on", "on"); ret == -1 {
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
	if ret := DefaultHttpReverseProxy.DeleteProxyClient(domain, host, port); !ret {
		self.WriteJsonString(w, `{"status":0}`)
	} else {
		self.WriteJsonString(w, `{"status":1}`)
	}
}

//http://127.0.0.1:8080/updateproxyclient?domain=www.xxx.com&prehost=121&preport=21&updatehost=xxxxxxx&updateport=1223
func (self *Http) UpdateProxyClient(w http.ResponseWriter, r *http.Request, prms httprouter.Params) {
	r.ParseForm()
	domain := r.Form.Get("domain")
	updateHost := r.Form.Get("updatehost")
	updatePort := r.Form.Get("updateport")
	preHost := r.Form.Get("prehost")
	prePort := r.Form.Get("preport")
	if domain == "" || updateHost == "" || updatePort == "" || preHost == "" || prePort == "" {
		self.WriteJsonString(w, `{"status":0}`)
		return
	}
	if ret := DefaultHttpReverseProxy.UpdateProxyClient(domain, preHost, prePort, updateHost, updatePort, "on", "on"); !ret {
		self.WriteJsonString(w, `{"status":0}`)
	} else {
		self.WriteJsonString(w, `{"status":1}`)
	}
}

//http://www.abc.com/proxyctl?protocol=http&switch=on
func (self *Http) ProxyControl(w http.ResponseWriter, r *http.Request, prms httprouter.Params) {
	r.ParseForm()
	proxyProtocol := r.Form.Get("protocol")
	proxySwitch := r.Form.Get("switch")
	var bSuccess = false
	if proxyProtocol == "https" {
		if proxySwitch == global.SwitchOn {
			bSuccess = DefaultHttpReverseProxy.StartAllHttpsService()
		} else if proxySwitch == global.SwitchOff {
			bSuccess = DefaultHttpReverseProxy.StopAllHttpsService()
		}
	} else if proxyProtocol == "http" {
		if proxySwitch == global.SwitchOn {
			bSuccess = DefaultHttpReverseProxy.StartAllHttpService()
		} else if proxySwitch == global.SwitchOff {
			bSuccess = DefaultHttpReverseProxy.StoptAllHttpService()
		}
	}
	if bSuccess {
		self.WriteJsonString(w, `{"status":1}`)
	} else {
		self.WriteJsonString(w, `{"status":0}`)
	}
}

//uploadfile
func (self *Http) UploadFile(w http.ResponseWriter, r *http.Request, prms httprouter.Params) {
	r.ParseForm()
	_, header, _ := r.FormFile("file")
	filetype := r.PostFormValue("filetype")
	domain := r.PostFormValue("domain")
	if filetype == "crt" {
		self.WriteJsonString(w, `{"status":1}`)
		fmt.Println(filetype)
		fmt.Println(header.Filename)
	} else if filetype == "key" {
		self.WriteJsonString(w, `{"status":1}`)
		fmt.Println(filetype)
		fmt.Println(domain)
	} else {
		self.WriteJsonString(w, `{"status":0}`)
	}
}

//create http service
func NewHttp(host, port string) *Http {
	return &Http{Host: host, Port: port}
}

//run router server
func (self *Http) Run() {
	log.Printf("Begin  Running http api service,%s:%s........\n", self.Host, self.Port)
	router := httprouter.New()
	//statistics
	router.GET("/clientinfos", self.ClientInfos)
	router.GET("/statistics", self.Statistics)
	router.GET("/routerinfo", self.RouterInfo)
	router.GET("/activeclients", self.ActiveClientInfos)
	router.GET("/bestclients", self.ActiveClientInfos)
	//Reverse Proxy  config file setting
	router.GET("/proxyinfos/:domain", self.ProxyInfos)
	router.GET("/adddomain/:domain", self.AddDomain)
	router.GET("/deldomain/:domain", self.DelDomain)
	router.GET("/updatedomain", self.UpdateDomain)
	router.GET("/addproxyclient", self.AddProxyClient)
	router.GET("/delproxyclient", self.DeleteProxyClient)
	router.GET("/updateproxyclient", self.UpdateProxyClient)
	router.GET("/domaininfos", self.DomainInfos)
	//reverse proxy switch
	router.GET("/proxyctl", self.ProxyControl)
	//statc file server
	router.GET("/", self.Index)
	//upload file
	router.POST("/uploadfile", self.UploadFile)
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.ServeFiles("/website/*filepath", http.Dir("website"))
	log.Fatal(http.ListenAndServe(self.Host+":"+self.Port, router))
}
