package boot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	. "ActivedRouter/global"
	"ActivedRouter/hook"
	"ActivedRouter/netservice"
)

//parse config
func parseConfigfile() {
	switch RunMode {
	case ServerMode:
		{
			//server config
			loadServerJsonConfig(ServerJsonConfig)
			//hook script
			hook.ParseHookScript(HookConfig)
		}
	case ClientMode:
		{
			//client mode
			loadClientJsonConfig(ClientConfig)
		}
	case ReverseProxyMode:
		{
			//server config
			loadServerJsonConfig(ServerJsonConfig)
			//proxy config
			netservice.ProxyHandler.LoadProxyConfig(HttpProxyConfig)
		}
	}
}

//load client json config file
func loadClientJsonConfig(config string) {
	file, err := os.Open(config)
	defer file.Close()
	if err != nil {
		log.Fatalln(err.Error())
	}
	if bts, err := ioutil.ReadAll(file); err != nil {
		log.Fatalln(err.Error())
	} else {
		var ClientMap map[string]interface{}
		err := json.Unmarshal(bts, &ClientMap)
		if err != nil {
			log.Fatalln("load client.json error")
		}
		domain, _ := ClientMap["domain"].(string)
		cluster, _ := ClientMap["cluster"].(string)
		serverList := ClientMap["router_list"].([]interface{})
		ConfigMap["domain"] = domain
		ConfigMap["cluster"] = cluster
		Cluster = cluster
		Domain = domain
		//server list
		var serverArr []string
		for _, v := range serverList {
			serverArr = append(serverArr, v.(string))
		}
		ConfigMap["serverlist"] = strings.Join(serverArr, "|")
		log.Println(ConfigMap)
	}
}

//loade serverConfig
func loadServerJsonConfig(config string) {
	if file, err := os.Open(config); err == nil {
		if bts, err := ioutil.ReadAll(file); err == nil {
			var serverConfig ServerConfigData
			if json.Unmarshal(bts, &serverConfig) != nil {
				goto Exit
			} else {
				ConfigMap["host"] = serverConfig.Host
				ConfigMap["port"] = serverConfig.Port
				ConfigMap["srvmode"] = serverConfig.ServerMode
				ConfigMap["httpport"] = serverConfig.HttpPort
				ConfigMap["httphost"] = serverConfig.HttpHost
				return
			}
		} else {
			goto Exit
		}
	} else {
		goto Exit
	}
Exit:
	log.Fatalln("server config load error!")
}
