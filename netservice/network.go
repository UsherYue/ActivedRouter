package netservice

import (
	"ActivedRouter/global"
	"log"
	"strings"
)

//控制服务启停
var NetworkChan = make(chan bool, 0)

//启动网络相关服务
func StartNetworkService() {
	switch global.RunMode {
	case "server":
		{
			log.Printf("Running Server Mode Service.......")
			go NewServer(global.ConfigMap["host"], global.ConfigMap["port"]).Run()
			go NewHttp(global.ConfigMap["httphost"], global.ConfigMap["httpport"]).Run()
		}
	case "client":
		{
			log.Printf("Running Client Mode Service.......")
			serverList := global.ConfigMap["serverlist"]
			serverListArr := strings.Split(serverList, "|")
			for _, server := range serverListArr {
				hostinfo := strings.Split(server, ":")
				go NewClient(hostinfo[0], hostinfo[1]).Run()
			}
		}
	case "proxy":
		{
			log.Println("Running Reserve Proxy.......")
			//开启server服务
			if ProxyHandler.ProxyMethod == global.Alived {
				go NewServer(global.ConfigMap["host"], global.ConfigMap["port"]).Run()
			}
			//开启http服务
			go NewHttp(global.ConfigMap["httphost"], global.ConfigMap["httpport"]).Run()
			//开启反向代理服务
			go ProxyHandler.StartProxyServer()
		}
		//	case "mix":
		//		{
		//			log.Printf("Running Mix Mode Service .......")
		//			go NewServer(global.ConfigMap["host"], global.ConfigMap["port"]).Run()
		//			go NewHttp(global.ConfigMap["httphost"], global.ConfigMap["httpport"]).Run()
		//			go ProxyHandler.StartProxyServer()
		//		}
	}
	NetworkChan <- true
}

//stop service
func StopNetworkService() {
	<-NetworkChan
}
