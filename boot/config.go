package boot

import (
	"ActivedRouter/global"
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func parseConfigfile() {
	switch global.RunMode {
	case "server":
		{
			global.ConfigFile = "config/server.ini"
			log.Println("开始加载dns路由服务配置!")

		}
	case "client":
		{
			global.ConfigFile = "config/client.ini"
		}
	case "proxy":
		{
			global.ConfigFile = "config/proxy.json"
			return
		}
	}
	log.Printf("正在加载配置文件 %s .......\n", global.ConfigFile)
	loadIni(global.ConfigFile)
	log.Println(global.ConfigMap)
	//必备的通用配置选项配置本机或者远程的ip或者端口号
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
	file, err := os.Open(global.ConfigFile)
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
