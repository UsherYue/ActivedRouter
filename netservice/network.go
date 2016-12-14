package netservice

import (
	"log"
	"strings"

	"ActivedRouter/global"
)

//start networkservice
func StartNetworkService() {
	switch global.RunMode {
	case global.ServerMode:
		{
			//启动路由服务
			go NewServer(global.ConfigMap["host"], global.ConfigMap["port"]).Run()
			//运行http服务
			go NewHttp(global.ConfigMap["httphost"], global.ConfigMap["httpport"]).Run()
			log.Println("ActivedRouter is Running In Server Mode...")
		}
	case global.ClientMode:
		{
			serverList := global.ConfigMap["serverlist"]
			serverListArr := strings.Split(serverList, "|")
			for _, server := range serverListArr {
				hostinfo := strings.Split(server, ":")
				//运行客户端代理
				go NewClient(hostinfo[0], hostinfo[1]).Run()
			}
			log.Println("ActivedRouter is Running  In Client Mode...")
		}
	case global.ReserveProxyMode:
		{
			//启动路由服务
			go NewServer(global.ConfigMap["host"], global.ConfigMap["port"]).Run()
			//开启http服务
			go NewHttp(global.ConfigMap["httphost"], global.ConfigMap["httpport"]).Run()
			//开启反向代理服务
			go ProxyHandler.StartProxyServer()
			log.Println("ActivedRouter is Running  In ReserveProxy Mode...")
		}
	case global.MixMode:
		{
			//启动路由服务
			go NewServer(global.ConfigMap["host"], global.ConfigMap["port"]).Run()
			//启动http服务
			go NewHttp(global.ConfigMap["httphost"], global.ConfigMap["httpport"]).Run()
			//启动反向代理
			go ProxyHandler.StartProxyServer()
			log.Println("ActivedRouter is Running  In Mix Mode...")
		}
	}
	global.NetworkSwitch <- true
}

//stop service
func StopNetworkService() {
	<-global.NetworkSwitch
}
