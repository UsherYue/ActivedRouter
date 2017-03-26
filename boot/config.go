package boot

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	. "ActivedRouter/global"
	"ActivedRouter/hook"
	"ActivedRouter/netservice"
)

//解析配置文件
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
	case ReserveProxyMode:
		{
			//server config
			loadServerJsonConfig(ServerJsonConfig)
			//proxy config
			netservice.ProxyHandler.LoadProxyConfig(HttpProxyConfig)
		}
	}
}

//加载json文件
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
			log.Fatalln("加载client.json失败")
		}
		domain, _ := ClientMap["domain"].(string)
		cluster, _ := ClientMap["cluster"].(string)
		serverList := ClientMap["router_list"].([]interface{})
		ConfigMap["domain"] = domain
		ConfigMap["cluster"] = cluster
		Cluster = cluster
		Domain = domain
		//服务器列表
		var serverArr []string
		for _, v := range serverList {
			serverArr = append(serverArr, v.(string))
		}
		ConfigMap["serverlist"] = strings.Join(serverArr, "|")
		log.Println(ConfigMap)
	}
}

//加载服务器模式配置
func loadServerModeConfig(config string) {
	loadIni(ServerConfig)
	if val, ok := ConfigMap["host"]; !ok || val == "" {
		log.Fatalln("配置文件缺少host键值......")
	}
	if val, ok := ConfigMap["port"]; !ok || val == "" {
		log.Fatalln("配置文件缺少port键值......")
	}
	//server模式下必须配置http服务器的ip端口号
	if RunMode == "server" {
		if val, ok := ConfigMap["httphost"]; !ok || val == "" {
			log.Fatalln("配置文件缺少httphost键值......")
		}
		if val, ok := ConfigMap["httpport"]; !ok || val == "" {
			log.Fatalln("配置文件缺少httpport键值......")
		}
		if val, ok := ConfigMap["srvmode"]; !ok || val == "" {
			log.Fatalln("配置文件缺少srvmode键值......")
		}
	}
}

//加载服务器json
func loadServerJsonConfig(config string) {
	if file, err := os.Open(config); err == nil {
		if bts, err := ioutil.ReadAll(file); err == nil {
			var serverConfig ServerConfigData
			if json.Unmarshal(bts, &serverConfig) != nil {
				goto Exit
			} else {
				//解析
				//log.Println(sercerConfig)
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

//load dns router
func loadDnsRouterConfig(routerFile string) {
	file, err := os.Open(routerFile)
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
		var DnsMap []map[string]interface{}
		err := json.Unmarshal(bts, &DnsMap)
		if err != nil {
			log.Fatalln(err.Error())
		} else {
			DnsScript = DnsMap
			log.Println(string(bts))
		}
	}
}

//load ini config
func loadIni(config string) {
	//打开文件
	file, err := os.Open(config)
	defer func() {
		file.Close()
	}()
	if err != nil {
		log.Fatalln(err.Error())
	}
	//读取文件内容
	reader := bufio.NewReader(file)
	itemStr := ""
	for {
		lineStr, err := reader.ReadString(byte('\n'))
		//判断文件读取结束
		if err != nil {
			break
		}
		//remove space
		itemStr = strings.TrimSpace(lineStr)
		if itemStr == "" {
			continue
		}
		//注释
		if strings.Index(itemStr, "#") == 0 {
			continue
		}
		//配置文件语法错误
		if -1 == strings.Index(itemStr, "=") {
			continue
		}
		kvs := strings.Split(itemStr, "=")
		ConfigMap[kvs[0]] = kvs[1]
	}

}
