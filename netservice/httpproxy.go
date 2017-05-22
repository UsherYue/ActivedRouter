package netservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"ActivedRouter/cache"
	"ActivedRouter/global"
	"ActivedRouter/system"
	"ActivedRouter/tools"
)

// var define
var (
	HttpAddr    string
	HttpsAddr   string
	HttpSwitch  string //on off
	HttpsSwitch string //on off
)

//default
var (
	DefaultHttpAddr = "127.0.0.1:80"
	DefaultHttsAddr = "127.0.0.1:443"
	SwitchOn        = "on"
	SwitchOff       = "off"
)

const (
	HTTP_STATISTICS_INTERVAL = 60 //http statistics interval 5min
)

//hander
var (
	ProxyHandler = &ReverseProxyHandler{Cfg: &ReverseProxyConfigData{}}
)

type HostInfo struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

//Load Balance Node
type LbNode struct {
	Domain      string      `json:"domain"`
	HttpsSwitch string      `json:"https_switch"`
	HttpSwitch  string      `json:"http_switch"`
	Clients     []*HostInfo `json:"clients"`
}

//ReverseProxy Config
type ReverseProxyConfigData struct {
	ProxyMethod   string `json:"proxy_method"`
	HttpProxyAddr string `json:"http_proxy_addr"`
	//global http switch
	GlobalHttpSwitch string `json:"http_switch"`
	//global https switch
	GlobalHttpsSwitch string                       `json:"https_switch"`
	HttpsProxyAddr    string                       `json:"https_proxy_addr"`
	ReverseProxy      []*LbNode                    `json:"reserve_proxy"`
	DomainProxySwitch map[string]map[string]string `json:"-"`
}

//reverse proxy handler
type ReverseProxyHandler struct {
	DomainHostList cache.Cacher
	Cfg            *ReverseProxyConfigData
	httpsServer    *HttpsServer
	//certificate config
	CertificateConfigData []*CertificateConfig
	ProxyCongfigFile      string
	ProxyMethod           string
}

//domain list
func (this *ReverseProxyHandler) DomainInfos() []string {
	data := *this.DomainHostList.GetStorage().GetData()
	keysArr := make([]string, 0)
	for k, _ := range data {
		keysArr = append(keysArr, k)
	}
	return keysArr
}

//Add the domain name to the configuration
func (this *ReverseProxyHandler) AddDomainConfig(domain string) bool {
	for _, v := range this.Cfg.ReverseProxy {
		if v.Domain == domain {
			return false
		}
	}
	this.Cfg.ReverseProxy = append(this.Cfg.ReverseProxy, &LbNode{Domain: domain})
	if this.SaveToFile() {
		//Hot update
		this.DomainHostList.Set(domain, []*HostInfo{})
		return true
	}
	return false
}

//save to file
func (this *ReverseProxyHandler) SaveToFile() bool {
	if bts, err := json.MarshalIndent(this.Cfg, "", " "); err != nil {
		return false
	} else {
		if file, err := os.OpenFile(this.ProxyCongfigFile, os.O_RDWR|os.O_TRUNC, os.ModePerm); err != nil {
			defer file.Close()
			return false
		} else {
			if _, err := file.Write(bts); err != nil {
				return false
			}
		}
	}
	return true
}

//Delete the domain name and sync to the configuration file
func (this *ReverseProxyHandler) DeleteDomainConig(domain string) bool {
	for k, v := range this.Cfg.ReverseProxy {
		if v.Domain == domain {
			//delete item
			ret, _ := tools.DeleteSlice(this.Cfg.ReverseProxy, k)
			this.Cfg.ReverseProxy = ret.([]*LbNode)
			//hot update
			this.DomainHostList.Del(domain)
			this.SaveToFile()
		}
	}
	return false
}

//delete reverse proxydomain
func (this *ReverseProxyHandler) DeleteProxyClient(domain, hostip, port string) bool {
	for _, v := range this.Cfg.ReverseProxy {
		if v.Domain == domain {
			for index, client := range v.Clients {
				if client.Host == hostip && client.Port == port {
					//delete item
					ret, _ := tools.DeleteSlice(v.Clients, index)
					v.Clients = ret.([]*HostInfo)
					//hot update
					if this.DomainHostList.Has(domain) {
						clientInfoList := this.GetDomainHostList(domain)
						for index, item := range clientInfoList {
							if item.Host == hostip && item.Port == port {
								resultSlice, _ := tools.DeleteSlice(clientInfoList, index)
								this.DomainHostList.Set(domain, resultSlice)
							}
						}
					}
					if this.SaveToFile() {
						return true
					} else {
						return false
					}
				}
			}
		}
	}
	return false
}

//Update Reverse Proxy Client Info
func (this *ReverseProxyHandler) UpdateProxyClient(domain, preHost, prePort, updateHost, updatePort, httpsSwitch, httpSwitch string) bool {
	for _, v := range this.Cfg.ReverseProxy {
		if v.Domain == domain {
			v.HttpsSwitch = httpsSwitch
			v.HttpSwitch = HttpSwitch
			for _, client := range v.Clients {
				if client.Host == preHost && client.Port == prePort {
					client.Host = updateHost
					client.Port = updatePort
					//hot update
					if this.DomainHostList.Has(domain) {
						clientInfoList := this.GetDomainHostList(domain)
						for _, item := range clientInfoList {
							if item.Host == preHost && item.Port == prePort {
								item.Host = updateHost
								item.Port = updatePort
							}
						}
					}
					if this.SaveToFile() {
						return true
					} else {
						return false
					}
				}
			}
		}
	}
	return true
}

//Add the reverse proxy client to the specified domain name
//Return Value
// -1  Repeat
//  0  Failure
//  1  Success
func (this *ReverseProxyHandler) AddProxyClient(domain, hostip, port, httsSwitch, httpSwitch string) int {
	for _, v := range this.Cfg.ReverseProxy {
		if v.Domain == domain {
			//proxy switch
			v.HttpsSwitch = httsSwitch
			v.HttpSwitch = httpSwitch
			for _, client := range v.Clients {
				if client.Host == hostip && client.Port == port {
					return -1
				}
			}
			//Add the domain name repeatedly!
			v.Clients = append(v.Clients, &HostInfo{port, hostip})
			//hot update
			if !this.DomainHostList.Has(domain) {
				this.DomainHostList.Set(domain, []*HostInfo{&HostInfo{port, hostip}})
			} else {
				clientList, _ := this.DomainHostList.Get(domain)
				clientInfoList, _ := clientList.([]*HostInfo)
				this.DomainHostList.Set(domain, append(clientInfoList, &HostInfo{port, hostip}))
			}
			if this.SaveToFile() {
				return 1
			} else {
				return 0
			}
		}
	}
	this.Cfg.ReverseProxy = append(this.Cfg.ReverseProxy, &LbNode{Domain: domain, HttpsSwitch: "off", Clients: []*HostInfo{&HostInfo{port, hostip}}})
	this.SaveToFile()
	return 1
}

//update domain
func (this *ReverseProxyHandler) UpdateDomain(preDomain, updateDomain, httpsSwitch, httpSwitch string) bool {
	for _, v := range this.Cfg.ReverseProxy {
		if v.Domain == preDomain {
			//proxy switch
			v.HttpsSwitch = httpsSwitch
			v.HttpSwitch = httpSwitch
			v.Domain = updateDomain
			//hot update
			data, _ := this.DomainHostList.Get(preDomain)
			this.DomainHostList.Del(preDomain)
			this.DomainHostList.Set(updateDomain, data)
			if this.SaveToFile() {
				return true
			} else {
				return false
			}
		}
	}
	return true
}

func (this *ReverseProxyHandler) ChangeSwitchStatus(domain, protocol, switchStatus string) {
	for _, v := range this.Cfg.ReverseProxy {
		if v.Domain == domain {
			switch protocol {
			case "http":
				{
					v.HttpSwitch = switchStatus
				}
			case "https":
				{
					v.HttpsSwitch = switchStatus
				}
			default:
				{
					return
				}
			}
			this.SaveToFile()
		}
	}
}

//hostlist by domain
func (this *ReverseProxyHandler) GetDomainHostList(domain string) []*HostInfo {
	v, _ := this.DomainHostList.Get(domain)
	vArr, _ := v.([]*HostInfo)
	return vArr
}

//random method
func (this *ReverseProxyHandler) GetRandomHost(domain string) *HostInfo {
	v, _ := this.DomainHostList.Get(domain)
	vArr, _ := v.([]*HostInfo)
	proxyCount := len(vArr)
	//fix bug :integer divide by zero
	if proxyCount == 0 {
		return nil
	}
	index := rand.Uint32() % uint32(proxyCount)
	return vArr[index]
}

//alived method
//According to the domain name or ip to obtain the most active cluster host
func (this *ReverseProxyHandler) GetAlivedHost(domain string) *HostInfo {
	v, _ := this.DomainHostList.Get(domain)
	vArr, _ := v.([]*HostInfo)
	hostinfo := this.BestHostInfo(vArr)
	return hostinfo
}

func (this *ReverseProxyHandler) BestHostInfo(hosts []*HostInfo) *HostInfo {
	hostSortedList := global.GHostInfoTable.ActiveHostWeightList
	for el := hostSortedList.Front(); el != nil; el = el.Next() {
		bestHost := el.Value.(system.HostInfo)
		for _, host := range hosts {
			//			log.Println(host.Host)
			//			log.Println(host.Port)
			//			log.Println(bestHost.Info.IP)
			//			log.Println(bestHost.Info.Domain)
			if bestHost.Info.IP == host.Host || bestHost.Info.Domain == host.Host {
				return host
			}
		}
	}
	return nil
}

//proxy_method  random 和alived
func (this *ReverseProxyHandler) GetHostInfo(host, proxyMethod string) *HostInfo {
	requestHost := host
	//Handle non-80 ports
	if strings.IndexAny(host, ":") != -1 {
		strs := strings.Split(host, ":")
		requestHost = strs[0]
	}
	//random
	//alived
	switch proxyMethod {
	case global.Random:
		{
			return this.GetRandomHost(requestHost)
		}
	case global.Alived:
		{
			return this.GetAlivedHost(requestHost)
		}
	}
	return nil
}

//serve http
//If the request protocol is https, check whether the reverse proxy is allowed to pass
func (this *ReverseProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.TLS != nil {
		if !this.httpsServer.CheckValidHttpsReq(r.Host) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(r.Host + "&nbsp;&nbsp;can't be accessed via https,please configure a digital certificate........."))
			return
		}
		if proxySwitch, ok := this.Cfg.DomainProxySwitch[r.Host]; ok {
			if httpsSwitch, ok := proxySwitch["https"]; ok {
				if httpsSwitch == SwitchOn {
					goto STARTPROXY
				}
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(r.Host + "&nbsp;&nbsp;Please open https proxy switch........."))
		}
	} else {
		if proxySwitch, ok := this.Cfg.DomainProxySwitch[r.Host]; ok {
			if httpsSwitch, ok := proxySwitch["http"]; ok {
				if httpsSwitch == SwitchOn {
					goto STARTPROXY
				}
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(r.Host + "&nbsp;&nbsp;Please open http proxy switch........."))
		}
	}

STARTPROXY:
	//Get the business server
	hostinfo := this.GetHostInfo(r.Host, this.ProxyMethod)
	if hostinfo == nil {
		//If you can't get the active host then use the random method。
		hostinfo = this.GetHostInfo(r.Host, global.Random)
		if hostinfo == nil {
			w.Write([]byte(r.Host + "Can't find active server........."))
			return
		}
	}
	//Redirect http request
	redirect := fmt.Sprintf("http://%s:%s", hostinfo.Host, hostinfo.Port)
	remote, err := url.Parse(redirect)
	if err != nil {
		panic(err)
	}
	// Not modifyed the http request header
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
	//Update reverse proxy statistics
	go global.GProxyHttpStatistics.UpdateClusterStatistics(r.Host, 0)
}

//Load Certificate Config
func (this *ReverseProxyHandler) LoadCertificateConfig(certificateConfigFile string) {
	stat, err := os.Stat(global.CertificateData)
	if os.IsNotExist(err) || !stat.IsDir() {
		return
	}
	if fileInfos, err := ioutil.ReadDir(global.CertificateData); err != nil {
		return
	} else {
		for _, fileInfo := range fileInfos {
			if fileInfo.IsDir() {
				domain := fileInfo.Name()
				certFile := path.Join(global.CertificateData, domain, global.DefaultCertificate)
				keyFile := path.Join(global.CertificateData, domain, global.DefaultKey)
				this.CertificateConfigData = append(this.CertificateConfigData, &CertificateConfig{Domain: domain, CertFile: certFile, KeyFile: keyFile})
			}
		}
	}
}

//Load proxy config
func (this *ReverseProxyHandler) LoadProxyConfig(proxyConfigFile string) {
	this.ProxyCongfigFile = proxyConfigFile
	file, err := os.Open(proxyConfigFile)
	defer file.Close()
	if err != nil {
		log.Fatalln(err.Error())
	}
	if bts, err := ioutil.ReadAll(file); err != nil {
		log.Fatalln(err.Error())
	} else {
		if err := json.Unmarshal(bts, &this.Cfg); err != nil {
			log.Fatalln("Parse proxy config file .....")
		}
		this.Cfg.DomainProxySwitch = make(map[string]map[string]string)
		HttpSwitch = this.Cfg.GlobalHttpSwitch
		HttpsSwitch = this.Cfg.GlobalHttpsSwitch
		//http https  off
		if HttpSwitch == SwitchOff && HttpsSwitch == SwitchOff {
			log.Fatalln("Please open http or https reverse proxy switch.....")
		}
		//Get the http switch
		if HttpSwitch == SwitchOn {
			if this.Cfg.HttpProxyAddr == "" {
				HttpAddr = DefaultHttpAddr
			} else {
				HttpAddr = this.Cfg.HttpProxyAddr
			}
			log.Println("Http Switch:" + HttpSwitch)
			log.Println("Http  Addr:" + HttpAddr)
		}
		//Get the https switch
		if HttpsSwitch == SwitchOn {
			if this.Cfg.HttpsProxyAddr == "" {
				HttpsAddr = DefaultHttsAddr
			} else {
				HttpsAddr = this.Cfg.HttpsProxyAddr
			}
			//			HttpsCrt = this.Cfg.HttpsCrt
			//			HttpsKey = this.Cfg.HttpsKey
			log.Println("Https Switch:" + HttpsSwitch)
			log.Println("Https Addr:" + HttpsAddr)
		}
		//Proxy method
		if this.Cfg.ProxyMethod == "" {
			this.ProxyMethod = global.Random
		} else {
			this.ProxyMethod = this.Cfg.ProxyMethod
		}
		//Create a memory cache to store the list of domain names
		this.DomainHostList = cache.Newcache("memory")
		clients := this.Cfg.ReverseProxy
		for _, client := range clients {
			subDomain := client.Domain
			//Domain proxy switch
			this.Cfg.DomainProxySwitch[subDomain] = map[string]string{"http": client.HttpSwitch, "https": client.HttpsSwitch}
			var subClientList []*HostInfo
			for _, hostInfo := range client.Clients {
				subClientList = append(subClientList, hostInfo)
			}
			this.DomainHostList.Set(subDomain, subClientList)
		}
	}
}

//Run the http statistics service
func (this *ReverseProxyHandler) BeginHttpStatistics() {
	timerStatistics := time.NewTimer(time.Second * HTTP_STATISTICS_INTERVAL)
	for {
		select {
		case <-timerStatistics.C:
			{
				//reset timer
				timerStatistics.Reset(time.Second * HTTP_STATISTICS_INTERVAL)
				//Incremental statistical curve(曲线)
				global.GProxyHttpStatistics.UpdateClusterStatistics("", 1)
			}
		}
	}
}

//Run Reverse Proxy
func (this *ReverseProxyHandler) StartProxyServer() {
	//http switch
	if HttpSwitch == SwitchOn {
		go func() {
			err := http.ListenAndServe(HttpAddr, ProxyHandler)
			if err != nil {
				log.Fatalln("ListenAndServe HTTP: ", err)
			} else {
				log.Println("Listen Http :", HttpAddr)
			}
		}()
	}
	//https switch
	if HttpsSwitch == SwitchOn {
		go func() {
			this.httpsServer = NewHttpsServer()
			this.httpsServer.AddDomainCertificateConfig(this.CertificateConfigData)
			err := this.httpsServer.RunHttpsService(HttpsAddr, "", "", ProxyHandler)
			if err != nil {
				log.Fatalln("RunHttpServer:", err)
			} else {
				log.Println("Listen Http SSL:", HttpsAddr)
			}
		}()
	}
	//Open http reverse proxy statistics
	//You can choose whether to open, because this option will affect the http request speed,
	// you can turn off.
	go this.BeginHttpStatistics()
}
