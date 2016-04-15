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
			NetworkChan <- true
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
			NetworkChan <- true
		}
	case "proxy":
		{
			log.Println("Running Reserve Proxy.......")
			if ProxyHandler.ProxyMethod == global.Alived {
				log.Fatalln("Reserve Proxy Actived method need mix runmode")
			}
			ProxyHandler.StartProxyServer()
		}
	case "mix":
		{
			log.Printf("Running Mix Mode Service .......")
			go NewServer(global.ConfigMap["host"], global.ConfigMap["port"]).Run()
			ProxyHandler.StartProxyServer()
		}
	}
}

//stop service
func StopNetworkService() {
	<-NetworkChan
}
