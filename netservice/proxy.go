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
)

var HTTPADDR string

var DefaultHttpAddr = "127.0.0.1:8888"

type HostInfo struct {
	Port string
	Host string
}

//hander
var ProxyHandler = &ReseveProxyHandler{}

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
	redirect := fmt.Sprintf("http://%s:%s", hostinfo.Host, hostinfo.Port)
	remote, err := url.Parse(redirect)
	if err != nil {
		panic(err)
	}
	//不修改 http request header
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
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
		if v, ok := proxyConfig["proxy_addr"]; !ok {
			HTTPADDR = DefaultHttpAddr
		} else {
			HTTPADDR, _ = v.(string)
		}
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
			data := *this.DomainHostList.GetMemory().GetData()
			log.Println(data)
		}
	}
}

//启动proxy
func (this *ReseveProxyHandler) StartProxyServer() {
	//被代理的服务器host和port
	err := http.ListenAndServe(HTTPADDR, ProxyHandler)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
