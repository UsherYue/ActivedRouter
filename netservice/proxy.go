/**
* reserve proxy
 */
package netservice

import (
	"ActivedRouter/cache"
	"ActivedRouter/global"
	"encoding/json"
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

//获取代理服务器
func (this *ReseveProxyHandler) GetHostInfo(domain string) *HostInfo {
	requestDomain := domain
	//处理非80端口
	if strings.IndexAny(domain, ":") != -1 {
		strs := strings.Split(domain, ":")
		requestDomain = strs[0]
	}
	v, _ := this.DomainHostList.Get(requestDomain)
	vArr, _ := v.([]*HostInfo)

	proxyCount := len(vArr)
	switch this.ProxyMethod {
	case "random":
		{
			index := rand.Uint32() % uint32(proxyCount)
			log.Println(vArr[index])
			return vArr[index]
		}
	}

	return nil
}

//serve http
func (this *ReseveProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hostinfo := this.GetHostInfo(r.Host)
	//log.Println(*hostinfo)
	remote, err := url.Parse("http://" + hostinfo.Host + ":" + hostinfo.Port)
	if err != nil {
		panic(err)
	}
	//不修改 http request header
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}

//load proxy
func (this *ReseveProxyHandler) loadProxyConfig() {
	file, err := os.Open(global.ConfigFile)
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
			log.Fatalln("反向代理配置加载失败!!")
		} else {
			//获取到不同域名
			this.DomainHostList = cache.Newcache("memory")
			clients := v.([]interface{})
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
func StartProxyServer() {
	//加载配置文件
	ProxyHandler.loadProxyConfig()
	//被代理的服务器host和port
	err := http.ListenAndServe(HTTPADDR, ProxyHandler)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
