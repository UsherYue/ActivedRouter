/**
* reserve proxy
 */
package netservice

import (
	"ActivedRouter/cache"
	"ActivedRouter/global"
	"ActivedRouter/system"
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
	DefaultHttpAddr = "127.0.0.1:8888"
	DefaultHttsAddr = "127.0.0.1:443"
	SwitchOn        = "on"
	SwitchOff       = "off"
)

const (
	HTTP_STATISTICS_INTERVAL = 60 //http统计周期 5min
)

//hander
var (
	ProxyHandler = &ReseveProxyHandler{}
)

type HostInfo struct {
	Port string
	Host string
}

//reserver
type ReseveProxyHandler struct {
	HostList       []*HostInfo
	DomainHostList cache.Cache
	ProxyMethod    string
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
	fmt.Println(remote)
	//不修改 http request header
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
	//更新转发统计
	global.GHttpStatistics.UpdateClusterStatistics(r.Host, 0)

}

//load proxy
func (this *ReseveProxyHandler) LoadProxyConfig(proxyConfigFile string) {
	file, err := os.Open(proxyConfigFile)
	defer func() {
		file.Close()
	}()
	if err != nil {
		log.Fatalln(err.Error())
	}
	//reader
	if bts, err := ioutil.ReadAll(file); err != nil {
		log.Fatalln(err.Error())
	} else {
		var proxyConfig map[string]interface{}
		json.Unmarshal(bts, &proxyConfig)
		//获取http和https开关 缺省是off
		if httpSwitch, b := proxyConfig["http_switch"]; !b {
			HttpSwitch = SwitchOff
		} else {
			HttpSwitch = httpSwitch.(string)
		}
		if httpsSwitch, b := proxyConfig["https_switch"]; !b {
			HttpsSwitch = SwitchOff
		} else {
			HttpsSwitch = httpsSwitch.(string)
		}
		//http https  off
		if HttpSwitch == SwitchOff && HttpsSwitch == SwitchOff {
			log.Fatalln("请开启http或者https代理开关.....")
		}
		//获取http开关下的配置
		if HttpSwitch == SwitchOn {
			//获取http http_proxy_addr
			if v, ok := proxyConfig["http_proxy_addr"]; !ok {
				HttpAddr = DefaultHttpAddr
			} else {
				HttpAddr, _ = v.(string)
			}
			log.Println("Http Switch:" + HttpSwitch)
			log.Println("Http  Addr:" + HttpAddr)
		}
		//获取https开关下的配置
		if HttpsSwitch == SwitchOn {
			//获取https https_proxy_addr
			if v, ok := proxyConfig["https_proxy_addr"]; !ok {
				HttpsAddr = DefaultHttsAddr
			} else {
				HttpsAddr, _ = v.(string)
			}
			//获取证书 和key
			if v, ok := proxyConfig["https_crt"]; !ok {
				log.Fatalln("https 缺少证书........")
			} else {
				HttpsCrt = v.(string)
			}
			if v, ok := proxyConfig["https_key"]; !ok {
				log.Fatalln("https 缺少key.........")
			} else {
				HttpsKey = v.(string)
			}
			log.Println("Https Switch:" + HttpsSwitch)
			log.Println("Https Addr:" + HttpsAddr)
			log.Println("Https  Crt:" + HttpsCrt)
			log.Println("Https  Key:" + HttpsKey)
		}

		//get proxy method
		if v, ok := proxyConfig["proxy_method"]; !ok {
			this.ProxyMethod = "random"
		} else {
			this.ProxyMethod, _ = v.(string)
		}
		if v, ok := proxyConfig["reserve_proxy"]; !ok {
			log.Fatalln("Reserve Proxy Config file load error !!")
		} else {
			//获取到不同域名
			this.DomainHostList = cache.Newcache("memory")
			clients := v.([]interface{})
			if len(clients) == 0 {
				log.Fatalln("Config file miss proxy host......")
			}
			for _, client := range clients {
				mapClient, _ := client.(map[string]interface{})
				subDomain := mapClient["domain"].(string)
				subClients := mapClient["clients"].([]interface{})
				var subClientList []*HostInfo
				for _, subClient := range subClients {
					subClientMap := subClient.(map[string]interface{})
					subClientHost := subClientMap["host"].(string)
					subClientPort := subClientMap["port"].(string)
					host := &HostInfo{Host: subClientHost, Port: subClientPort}
					subClientList = append(subClientList, host)
				}
				this.DomainHostList.Set(subDomain, subClientList)
			}
			//this.DomainHostList.GetMemory().GetData()
			//log.Println(data)
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
				global.GHttpStatistics.UpdateClusterStatistics("", 1)
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
