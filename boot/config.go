package boot

import (
	"ActivedRouter/global"
	"ActivedRouter/hook"
	"ActivedRouter/netservice"
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func parseConfigfile() {
	switch global.RunMode {
	case global.ServerMode:
		{
			//server config
			loadServerModeConfig(global.ServerConfig)
			//hook script
			hook.ParseHookScript(global.HookConfig)
		}
	case global.ClientMode:
		{
			//client mode
			loadJsonConfig(global.ClientConfig)
		}
	case global.ProxyMode:
		{
			//proxy config
			netservice.ProxyHandler.LoadProxyConfig(global.ProxyConfig)
			return
		}
	case global.MixMode:
		{
			//server config
			loadServerModeConfig(global.ServerConfig)
			//proxy config
			netservice.ProxyHandler.LoadProxyConfig(global.ProxyConfig)
		}
	}
}

//加载json文件
func loadJsonConfig(config string) {
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
		global.ConfigMap["domain"] = domain
		global.ConfigMap["cluster"] = cluster
		//服务器列表
		var serverArr []string
		for _, v := range serverList {
			serverArr = append(serverArr, v.(string))
		}
		global.ConfigMap["serverlist"] = strings.Join(serverArr, "|")
		log.Println(global.ConfigMap)
	}
}

//加载客户端模式下的配置
func loadClientModeConfig(config string) {
	loadIni(global.ClientConfig)
	log.Println(global.ConfigMap)
	if val, ok := global.ConfigMap["host"]; !ok || val == "" {
		log.Fatalln("-------配置文件缺少host键值-------")
	}
	if val, ok := global.ConfigMap["port"]; !ok || val == "" {
		log.Fatalln("-------配置文件缺少port键值-------")
	}
	//客户端模式下的配置
	if global.RunMode == "client" {
		//获取集群分组
		if val, ok := global.ConfigMap["cluster"]; ok {
			global.Cluster = val
		}
		log.Println("集群分组:", global.Cluster)
		//获取域名配置
		if val, ok := global.ConfigMap["domain"]; ok {
			global.Domain = val
		}
		log.Println("域名:", global.Cluster)
	}
}

//加载服务器模式配置
func loadServerModeConfig(config string) {
	loadIni(global.ServerConfig)
	log.Println(global.ConfigMap)
	if val, ok := global.ConfigMap["host"]; !ok || val == "" {
		log.Fatalln("-------配置文件缺少host键值-------")
	}
	if val, ok := global.ConfigMap["port"]; !ok || val == "" {
		log.Fatalln("-------配置文件缺少port键值-------")
	}
	//server模式下必须配置http服务器的ip端口号
	if global.RunMode == "server" {
		if val, ok := global.ConfigMap["httphost"]; !ok || val == "" {
			log.Fatalln("-------配置文件缺少httphost键值-------")
		}
		if val, ok := global.ConfigMap["httpport"]; !ok || val == "" {
			log.Fatalln("-------配置文件缺少httpport键值-------")
		}
		if val, ok := global.ConfigMap["srvmode"]; !ok || val == "" {
			log.Fatalln("-------配置文件缺少srvmode键值-------")
		}
	}
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
			global.DnsScript = DnsMap
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
		global.ConfigMap[kvs[0]] = kvs[1]
	}

}
