/*
* global data
 */
package global

import (
	"ActivedRouter/cache"
	"ActivedRouter/system"
	"sync"
)

//global mem cache
var GlobalCache cache.Cache = cache.Newcache("memory")

//存储服务器相关状态
var GHostInfoTable = system.NewHostInfoTable()

//路由服务器信息
var GRouterInfo string = ""

//dns配置脚本
var DnsScript []map[string]interface{}

//配置选项
var ConfigMap map[string]string = map[string]string{}

//运行模式
var RunMode = ""

//客户端的集群分组
var Cluster = ""

// 客户端配置下的域名
var Domain = ""

//全局http反向代理统计
var GHttpStatistics = system.NewSysHttpStatistics()

//安全的设置GRouterInfo
var rwMutexRouterInfo = &sync.RWMutex{}

//安全的设置安全的设置GRouterInfo
func SetRouterInfo(info string) {
	rwMutexRouterInfo.Lock()
	GRouterInfo = info
	rwMutexRouterInfo.Unlock()
}

//安全的检索GRouterInfo
func RouterInfo() string {
	rwMutexRouterInfo.RLock()
	info := GRouterInfo
	rwMutexRouterInfo.RUnlock()
	return info
}
