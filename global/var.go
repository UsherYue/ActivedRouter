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

//GLOBAL CACHE
var GlobalCache cache.Cacher = cache.Newcache("memory")

//Control network service start and stop
var NetworkSwitch = make(chan bool, 0)

//Store server related status
var GHostInfoTable = system.NewHostInfoTable()

//Router System Info
var RouterSysInfo string = ""

//Run Mode
var RunMode = ""

//SRVMODE
var SrvMode = "monitor"

//Global http reverse proxy statistics
var GProxyHttpStatistics = system.NewSysHttpStatistics()

//Read-Write Mutex
var rwMutexRouterInfo = &sync.RWMutex{}

func SetRouterSysInfo(info string) {
	rwMutexRouterInfo.Lock()
	RouterSysInfo = info
	rwMutexRouterInfo.Unlock()
}

func RouterInfo() string {
	rwMutexRouterInfo.RLock()
	info := RouterSysInfo
	rwMutexRouterInfo.RUnlock()
	return info
}
