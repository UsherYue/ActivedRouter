//ActivedRouter
//Author:usher.yue
//Amail:usher.yue@gmail.com
//TencentQQ:4223665
// Provide http/https, tcp and reverse proxy  services

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
			//Run Routing Service
			go NewServer(global.ConfigMap["host"], global.ConfigMap["port"]).Run()
			//Run Http Service
			go NewHttp(global.ConfigMap["httphost"], global.ConfigMap["httpport"]).Run()
			log.Println("ActivedRouter is Running In Server Mode...")
		}
	case global.ClientMode:
		{
			serverList := global.ConfigMap["serverlist"]
			serverListArr := strings.Split(serverList, "|")
			for _, server := range serverListArr {
				hostinfo := strings.Split(server, ":")
				//Run Client Agent
				go NewClient(hostinfo[0], hostinfo[1]).Run()
			}
			log.Println("ActivedRouter is Running  In Client Mode...")
		}
	case global.ReserveProxyMode:
		{
			//Run Routing Service
			go NewServer(global.ConfigMap["host"], global.ConfigMap["port"]).Run()
			//Run Http Service
			go NewHttp(global.ConfigMap["httphost"], global.ConfigMap["httpport"]).Run()
			//Run ReserveProxy Service
			go ProxyHandler.StartProxyServer()
			log.Println("ActivedRouter is Running  In ReserveProxy Mode...")
		}
	}
	ListenAndServePProf(global.HTTP_PPROF_DEFAULT_ADDR, nil)
	global.NetworkSwitch <- true
}

//stop service
func StopNetworkService() {
	<-global.NetworkSwitch
}
