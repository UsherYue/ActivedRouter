package system

import (
	"ActivedRouter/cache"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

//host sync
var hosttableMutex = sync.RWMutex{}

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

//host info list
type HostInfoTable struct {
	HostsInfo cache.Cache
	//active router
	ActiveHostList *HostList
}

//新创建主机信息存储
func NewHostInfoTable() *HostInfoTable {
	table := &HostInfoTable{}
	//挂载主机列表
	table.HostsInfo = cache.Newcache("memory")
	//活跃主机列表
	table.ActiveHostList = NewHostList()
	return table
}

//update host table status
func (self *HostInfoTable) UpdateHostStatus() {
	hosttableMutex.Lock()
	defer hosttableMutex.Unlock()
	cacheMap := *self.HostsInfo.GetMemory().GetData()
	for _, v := range cacheMap {
		hostInfo := v.(*HostInfo)
		//超过最大非活跃时间间隔
		if time.Now().Unix()-hostInfo.LastActive > UNACTIVE_TIMEOUT {
			hostInfo.Status = UNACTIVE
			//从活跃主机列表移除
		}
		//更新活跃主机列表 或者添加主机列表,主机不活跃的时候从列表删除
		//unactive下线
		//active 存活
		self.ActiveHostList.UpdateHostList(hostInfo)
	}
	//self.ActiveHostList.DumpInfo()
}

//calc host weight
func (self *HostInfoTable) calcHostWeight() {
	hosttableMutex.Lock()
	defer hosttableMutex.Unlock()
	cacheMap := *self.HostsInfo.GetMemory().GetData()
	for _, v := range cacheMap {
		hostInfo := v.(*HostInfo)
		//计算服务器权重 分为活跃列表和非活跃列表
		hostInfo.LastActive = 1
	}

}

//更新服务器状态  不存在插入 存在更新
func (self *HostInfoTable) UpdateHostTable(ip string, info *SystemInfo) {
	hosttableMutex.Lock()
	defer hosttableMutex.Unlock()
	//更新状态
	hostInfo := &HostInfo{}
	hostInfo.Info = info
	hostInfo.Status = ACTIVE
	hostInfo.LastActive = time.Now().Unix() //unix timestamp
	self.HostsInfo.Set(ip, hostInfo)
	//self.DumpInfo()

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
	hosttableMutex.RLock()
	defer hosttableMutex.RUnlock()
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
	return &safeRet
}
