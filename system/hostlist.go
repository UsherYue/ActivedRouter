package system

import (
	"ActivedRouter/cache"
	"encoding/json"
	"fmt"
	"sync"
)

//active ROUTER LIST
type HostList struct {
	ActiveHostInfo cache.Cache
	rwMutex        sync.RWMutex
}

//新创建主机信息存储
func NewHostList() *HostList {
	table := &HostList{}
	table.rwMutex = sync.RWMutex{}
	table.ActiveHostInfo = cache.Newcache("memory")
	return table
}

//更新host状态
//移除unactive主机
func (self *HostList) UpdateHostList(pHostInfo *HostInfo) {
	strTargetIp := pHostInfo.Info.IP
	self.rwMutex.RLock()
	_, ok := self.ActiveHostInfo.Get(strTargetIp)
	self.rwMutex.RUnlock()
	if !ok {
		if pHostInfo.Status == "active" {
			//fmt.Println(strTargetIp + "已经上线.......")
			self.rwMutex.Lock()
			//加入活跃主机
			self.ActiveHostInfo.Set(strTargetIp, pHostInfo)
			self.rwMutex.Unlock()
		}
	} else {
		if pHostInfo.Status == "unactive" {
			self.rwMutex.Lock()
			//fmt.Println(strTargetIp + "已经下线.......")
			//加入活跃主机
			self.ActiveHostInfo.Del(strTargetIp)
			self.rwMutex.Unlock()
		}
	}

}

//dump
func (self *HostList) DumpInfo() {

	mapChan := *self.ActiveHostInfo.GetMemory().GetData()
	for k, v := range mapChan {
		fmt.Println(k)
		bts, _ := json.MarshalIndent(v, "", " ")
		fmt.Println(string(bts))
	}

}
