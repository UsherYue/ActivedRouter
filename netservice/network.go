package netservice

import (
	"ActivedRouter/global"
	"log"
)

//控制服务启停
var NetworkChan = make(chan bool, 0)

//启动网络相关服务
func StartNetworkService() {
	switch global.RunMode {
	case "server":
		{
			log.Printf("正在启动服务器模式下的网络服务.......")
			go NewServer(global.ConfigMap["host"], global.ConfigMap["port"]).Run()
			go NewHttp(global.ConfigMap["httphost"], global.ConfigMap["httpport"]).Run()
			NetworkChan <- true
		}
	case "client":
		{
			log.Printf("正在启动客户端模式下的网络服务.......")
			go NewClient(global.ConfigMap["host"], global.ConfigMap["port"]).Run()
			NetworkChan <- true
		}
	case "proxy":
		{
			log.Println("正在启动reserve proxy.......")
			StartProxyServer()
		}
	}

}

//stop service
func StopNetworkService() {
	<-NetworkChan
}

func init() {

}
