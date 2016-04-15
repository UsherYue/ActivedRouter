/*
* global data
 */
package global

import (
	"ActivedRouter/cache"
	"ActivedRouter/system"
)

//global mem cache
var GlobalCache cache.Cache = cache.Newcache("memory")

//存储服务器相关状态
var GHostInfoTable = system.NewHostInfoTable()

//active router list
var GActiveRouterList = system.NewRouterList()

//dns配置脚本
var DnsScript []map[string]interface{}

//配置选项
var ConfigMap map[string]string = map[string]string{}

//运行模式
var RunMode = ""

//客户端的集群分组
var Cluster = ""

// 客户端配置下的域名
var Domain = "www.api1.com"
