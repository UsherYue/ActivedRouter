/**
* reserve proxy
 */
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
	HttpsCrt    string //https 证书
	HttpsKey    string //https key
)

//default
var (
	DefaultHttpAddr = "127.0.0.1:80"
	DefaultHttsAddr = "127.0.0.1:443"
	SwitchOn        = "on"
	SwitchOff       = "off"
)

const (
	HTTP_STATISTICS_INTERVAL = 60 //http统计周期 5min
)

//hander
var (
	ProxyHandler = &ReseveProxyHandler{Cfg: &ReserveProxyConfigData{}}
)

type HostInfo struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

//负载均衡节点
type LbNode struct {
	Domain  string      `json:"domain"`
	Clients []*HostInfo `json:"clients"`
}

//反向代理配置
type ReserveProxyConfigData struct {
	ProxyMethod    string    `json:"proxy_method"`
	HttpProxyAddr  string    `json:"http_proxy_addr"`
	HttpSwitch     string    `json:"http_switch"`
	HttpsSwitch    string    `json:"https_switch"`
	HttpsCrt       string    `json:"https_crt"`
	HttpsKey       string    `json:"https_key"`
	HttpsProxyAddr string    `json:"https_proxy_addr"`
	ReserveProxy   []*LbNode `json:"reserve_proxy"`
}

//reserver
type ReseveProxyHandler struct {
	DomainHostList   cache.Cacher
	Cfg              *ReserveProxyConfigData
	ProxyCongfigFile string
	ProxyMethod      string
}

//domain list
func (this *ReseveProxyHandler) DomainInfos() []string {
	data := *this.DomainHostList.GetStorage().GetData()
	keysArr := make([]string, 0)
	for k, _ := range data {
		keysArr = append(keysArr, k)
	}
	return keysArr
}

//增加域名 同步到文件
func (this *ReseveProxyHandler) AddDomainConfig(domain string) bool {
	for _, v := range this.Cfg.ReserveProxy {
		if v.Domain == domain {
			return true
		}
	}
	this.Cfg.ReserveProxy = append(this.Cfg.ReserveProxy, &LbNode{Domain: domain})
	this.SaveToFile()
	return true
}
func (this *ReseveProxyHandler) SaveToFile() bool {
	if bts, err := json.Marshal(this.Cfg); err != nil {
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

func (this *ReseveProxyHandler) DeletenLbNodeSlice(slice []*LbNode, index int) []*LbNode {
	length := len(slice)
	if slice == nil || length == 0 || length < index {
		return nil
	}
	if length-1 == index {
		return slice[:index]
	} else if length > index {
		return append(slice[:index], slice[index+1:]...)
	}
	return nil
}
func (this *ReseveProxyHandler) DeleteClientSlice(slice []*HostInfo, index int) []*HostInfo {
	length := len(slice)
	if slice == nil || length == 0 || length < index {
		return nil
	}
	if length-1 == index {
		return slice[:index]
	} else if length > index {
		return append(slice[:index], slice[index+1:]...)
	}
	return nil
}

//删除域名 同步到文件
func (this *ReseveProxyHandler) DeleteDomainConig(domain string) bool {
	for k, v := range this.Cfg.ReserveProxy {
		if v.Domain == domain {
			//this.Cfg.ReserveProxy = append(this.Cfg.ReserveProxy[:k], this.Cfg.ReserveProxy[k+1:]...)
			//this.Cfg.ReserveProxy = this.DeletenLbNodeSlice(this.Cfg.ReserveProxy, k)
			ret, _ := tools.DeleteSlice(this.Cfg.ReserveProxy, k)
			this.Cfg.ReserveProxy = ret.([]*LbNode)
		}
	}
	this.SaveToFile()
	return true
}

//删除反向代理服务器 ,并且删除配置文件中信息
func (this *ReseveProxyHandler) DeleteProxyClient(domain, hostip, port string) bool {
	for _, v := range this.Cfg.ReserveProxy {
		if v.Domain == domain {
			for index, client := range v.Clients {
				if client.Host == hostip && client.Port == port {
					//删除元素
					//v.Clients = append(v.Clients[:index], v.Clients[index+1:]...)
					//v.Clients = this.DeleteClientSlice(v.Clients[:index], index)
					ret, _ := tools.DeleteSlice(v.Clients, index)
					v.Clients = ret.([]*HostInfo)
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

//增加反向代理服务器 ,并且增加配置文件中信息
//-1 重复添加 0 失败 1成功
func (this *ReseveProxyHandler) AddProxyClient(domain, hostip, port string) (bool, int) {
	for _, v := range this.Cfg.ReserveProxy {
		if v.Domain == domain {
			for _, client := range v.Clients {
				if client.Host == hostip && client.Port == port {
					return false, -1
				}
			}
			//域名存在无重复配置添加client
			v.Clients = append(v.Clients, &HostInfo{port, hostip})
			if this.SaveToFile() {
				return true, 1
			} else {
				return false, 0
			}
		}
	}
	//域名不存在添加域名=>clients
	this.Cfg.ReserveProxy = append(this.Cfg.ReserveProxy, &LbNode{domain, []*HostInfo{&HostInfo{port, hostip}}})
	this.SaveToFile()
	return true, 1
}

//hostlist by domain
func (this *ReseveProxyHandler) GetDomainHostList(domain string) []*HostInfo {
	v, _ := this.DomainHostList.Get(domain)
	vArr, _ := v.([]*HostInfo)
	return vArr
}

//random method
func (this *ReseveProxyHandler) GetRandomHost(domain string) *HostInfo {
	v, _ := this.DomainHostList.Get(domain)
	vArr, _ := v.([]*HostInfo)
	proxyCount := len(vArr)
	index := rand.Uint32() % uint32(proxyCount)
	return vArr[index]
}

//alived method
//根据域名 或者ip获取集群最活跃的主机
func (this *ReseveProxyHandler) GetAlivedHost(domain string) *HostInfo {
	v, _ := this.DomainHostList.Get(domain)
	vArr, _ := v.([]*HostInfo)
	hostinfo := this.BestHostInfo(vArr)
	return hostinfo

}

//获取最优服务器根据服务器权限
//根据域名对指定集群进行路由
func (this *ReseveProxyHandler) BestHostInfo(hosts []*HostInfo) *HostInfo {
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

//根据负载方法进行主机筛选
//proxy_method  random 和alived
func (this *ReseveProxyHandler) GetHostInfo(host, proxyMethod string) *HostInfo {
	requestHost := host
	//处理非80端口
	if strings.IndexAny(host, ":") != -1 {
		strs := strings.Split(host, ":")
		requestHost = strs[0]
	}
	//random proxy模式即可
	//alived 需要开启mix模式
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
func (this *ReseveProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//获取服务器
	hostinfo := this.GetHostInfo(r.Host, this.ProxyMethod)
	//如果获取不到挂载主机那么使用random
	if hostinfo == nil {
		//hostinfo = this.GetHostInfo(r.Host, global.Random)
		if hostinfo == nil {
			w.Write([]byte(r.Host + "没有发现相关集群活跃服务器........."))
			return
		}
	}
	//重定向http请求
	redirect := fmt.Sprintf("http://%s:%s", hostinfo.Host, hostinfo.Port)
	remote, err := url.Parse(redirect)
	if err != nil {
		panic(err)
	}
	//fmt.Println(remote)
	//不修改 http request header
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
	//更新转发统计
	global.GProxyHttpStatistics.UpdateClusterStatistics(r.Host, 0)

}

//load proxy
func (this *ReseveProxyHandler) LoadProxyConfig(proxyConfigFile string) {
	this.ProxyCongfigFile = proxyConfigFile
	file, err := os.Open(proxyConfigFile)
	defer file.Close()
	if err != nil {
		log.Fatalln(err.Error())
	}
	//reader
	if bts, err := ioutil.ReadAll(file); err != nil {
		log.Fatalln(err.Error())
	} else {
		json.Unmarshal(bts, &this.Cfg)
		HttpSwitch = this.Cfg.HttpSwitch
		HttpsSwitch = this.Cfg.HttpsSwitch
		//http https  off
		if HttpSwitch == SwitchOff && HttpsSwitch == SwitchOff {
			log.Fatalln("请开启http或者https代理开关.....")
		}
		//获取http开关下的配置
		if HttpSwitch == SwitchOn {
			HttpAddr = this.Cfg.HttpProxyAddr
			if HttpAddr == "" {
				HttpAddr = DefaultHttpAddr
			}
			log.Println("Http Switch:" + HttpSwitch)
			log.Println("Http  Addr:" + HttpAddr)
		}
		//获取https开关下的配置
		if HttpsSwitch == SwitchOn {
			HttpsAddr = this.Cfg.HttpsProxyAddr
			if HttpsAddr == "" {
				HttpsAddr = DefaultHttsAddr
			}
			HttpsCrt = this.Cfg.HttpsCrt
			HttpsKey = this.Cfg.HttpsKey
			log.Println("Https Switch:" + HttpsSwitch)
			log.Println("Https Addr:" + HttpsAddr)
			log.Println("Https  Crt:" + HttpsCrt)
			log.Println("Https  Key:" + HttpsKey)
		}

		if this.Cfg.ProxyMethod == "" {
			this.ProxyMethod = global.Random
		} else {
			this.ProxyMethod = this.Cfg.ProxyMethod
		}
		//获取到不同域名
		this.DomainHostList = cache.Newcache("memory")
		clients := this.Cfg.ReserveProxy
		if len(clients) == 0 {
			log.Fatalln("Config file miss proxy host......")
		}
		for _, client := range clients {
			subDomain := client.Domain
			var subClientList []*HostInfo
			for _, hostInfo := range client.Clients {
				subClientList = append(subClientList, hostInfo)
			}
			this.DomainHostList.Set(subDomain, subClientList)
		}
	}
}

//开始进入定时统计
func (this *ReseveProxyHandler) BeginHttpStatistics() {
	timerStatistics := time.NewTimer(time.Second * HTTP_STATISTICS_INTERVAL)
	for {
		select {
		case <-timerStatistics.C:
			{
				//reset timer
				timerStatistics.Reset(time.Second * HTTP_STATISTICS_INTERVAL)
				//递增统计曲线
				global.GProxyHttpStatistics.UpdateClusterStatistics("", 1)
			}
		}
	}
}

//启动proxy
func (this *ReseveProxyHandler) StartProxyServer() {
	//http switch
	if HttpSwitch == SwitchOn {
		go func() {
			//被代理的服务器host和port
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
			//被代理的服务器host和port
			err := http.ListenAndServeTLS(HttpsAddr, HttpsCrt, HttpsKey, ProxyHandler)
			if err != nil {
				log.Fatalln("ListenAndServe HTTP SSL: ", err)
			} else {
				log.Println("Listen Http SSL:", HttpsAddr)
			}
		}()
	}
	//开启http反向代理统计
	//可选择是否开启 因为此选项会影响http请求速度 关闭可优化速度
	go this.BeginHttpStatistics()
}
