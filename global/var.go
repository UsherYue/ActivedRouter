//ActivedRouter
//Author:usher.yue
//Amail:usher.yue@gmail.com
//TencentQQ:4223665
//Global Var

package global

import (
	"sync"

	"ActivedRouter/cache"
	"ActivedRouter/system"
)

//cache driver
//memory
//mysql
//file
//redis
//mongodb
var GlobalCache cache.Cacher = cache.Newcache("memory")

//Control network service start and stop
var NetworkSwitch = make(chan bool, 0)

//Store server related status
var GHostInfoTable = system.NewHostInfoTable()

var GRouterInfo string = ""

//Config  Mapping
var ConfigMap map[string]string = map[string]string{}

//Run Mode
var RunMode = ""

//Client cluster name
var Cluster = ""

// Client Domain Name Use in Client Mode .
var Domain = ""

//Global http reverse proxy statistics
var GProxyHttpStatistics = system.NewSysHttpStatistics()

//Read-Write Mutex
var rwMutexRouterInfo = &sync.RWMutex{}

func SetRouterInfo(info string) {
	rwMutexRouterInfo.Lock()
	GRouterInfo = info
	rwMutexRouterInfo.Unlock()
}

func RouterInfo() string {
	rwMutexRouterInfo.RLock()
	info := GRouterInfo
	rwMutexRouterInfo.RUnlock()
	return info
}
