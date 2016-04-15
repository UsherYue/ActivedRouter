package system

import (
	"ActivedRouter/cache"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

//host sync
var hosttableMutex = sync.Mutex{}

const (
	UNACTIVE_TIMEOUT = 5 //unactive seconds
)

const (
	ACTIVE   = "active"
	UNACTIVE = "unactive"
)

//服务器信息
type HostInfo struct {
	Info       *SystemInfo // info
	Status     string      //active noactive
	LastActive int64
}

//active router
type HostRouteInfo struct {
	DNS    string
	IP     string
	Weight int
}

//active ROUTER LIST
type RouterList struct {
	ActiveHostInfo cache.Cache
}

//host info list
type HostInfoTable struct {
	HostsInfo cache.Cache
}

//新创建主机信息存储
func NewRouterList() *RouterList {
	table := &RouterList{}
	table.ActiveHostInfo = cache.Newcache("memory")
	return table
}

//新创建主机信息存储
func NewHostInfoTable() *HostInfoTable {
	table := &HostInfoTable{}
	table.HostsInfo = cache.Newcache("memory")
	return table
}

//update host table status
func (self *HostInfoTable) UpdateHostStatus() {
	cacheMap := *self.HostsInfo.GetMemory().GetData()
	for _, v := range cacheMap {
		hostInfo := v.(*HostInfo)
		//超过最大非活跃时间间隔
		if time.Now().Unix()-hostInfo.LastActive > UNACTIVE_TIMEOUT {
			hostInfo.Status = UNACTIVE
		}
	}
}

//calc host weight
func (self *HostInfoTable) calcHostWeight() {
	hosttableMutex.Lock()
	cacheMap := *self.HostsInfo.GetMemory().GetData()
	for _, v := range cacheMap {
		hostInfo := v.(*HostInfo)
		//计算服务器权重 分为活跃列表和非活跃列表
		hostInfo.LastActive = 1
	}
	hosttableMutex.Unlock()
}

//更新服务器状态  不存在插入 存在更新
func (self *HostInfoTable) UpdateHostTable(ip string, info *SystemInfo) {
	hosttableMutex.Lock()
	//更新服务器列表 此处需要有专门的一个服务定时的计算服务器权重
	//更新状态
	hostInfo := &HostInfo{}
	hostInfo.Info = info
	hostInfo.Status = ACTIVE
	hostInfo.LastActive = time.Now().Unix() //unix timestamp
	self.HostsInfo.Set(ip, hostInfo)
	//self.DumpInfo()
	hosttableMutex.Unlock()
}

//dump
func (self *HostInfoTable) DumpInfo() {

	mapChan := *self.HostsInfo.GetMemory().GetData()
	for k, v := range mapChan {
		fmt.Println(k)
		bts, _ := json.MarshalIndent(v, "", " ")
		fmt.Println(string(bts))
	}

}

//获取服务器信息
//如果服务器存在的话 初始化服务器信息是null
//初始化服务器是status 是 noactive  活跃是 active
func (self *HostInfoTable) GetHostInfo(ip string) *HostInfo {
	hosttableMutex.Lock()
	value, ok := self.HostsInfo.Get(ip)
	if !ok {
		return nil
	}
	systemInfo, ok := value.(*HostInfo)
	if !ok {
		return nil
	}
	//保证返回副本
	safeRet := *systemInfo
	hosttableMutex.Unlock()
	return &safeRet
}
