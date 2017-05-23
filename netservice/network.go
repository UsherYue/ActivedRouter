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
			go NewServer().Run(ServerConfigData.Host, ServerConfigData.Port)
			//Run Http Service
			go NewHttp(ServerConfigData.HttpHost, ServerConfigData.HttpPort).Run()
			log.Println("ActivedRouter is Running In Server Mode...")
		}
	case global.ClientMode:
		{
			for _, server := range ClientConfigData.RouterList {
				hostinfo := strings.Split(server, ":")
				//Run Client Agent
				go NewClient(hostinfo[0], hostinfo[1]).Run()
			}
			log.Println("ActivedRouter is Running  In Client Mode...")
		}
	case global.ReverseProxyMode:
		{
			//Run Routing Service
			go NewServer().Run(ServerConfigData.Host, ServerConfigData.Port)
			//Run Http Service
			go NewHttp(ServerConfigData.HttpHost, ServerConfigData.HttpPort).Run()
			//Run ReserveProxy Service
			go DefaultHttpReverseProxy.StartProxyServer()
			log.Println("ActivedRouter is Running  In ReverseProxy Mode...")
		}
	}
	ListenAndServePProf(global.HTTP_PPROF_Default_Addr, nil)
	global.NetworkSwitch <- true
}

//stop service
func StopNetworkService() {
	<-global.NetworkSwitch
}
