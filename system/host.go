package system

import (
	"ActivedRouter/cache"
	"container/list"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"
)

//host sync
var hosttableMutex = sync.RWMutex{}

//权重
var hostWeightMutex = sync.RWMutex{}

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

//host info list
type HostInfoTable struct {
	HostsInfo cache.Cache
	//active router
	ActiveHostList *HostList
	//active host weight list,sorted by weight
	ActiveHostWeightList *list.List
}

//新创建主机信息存储
func NewHostInfoTable() *HostInfoTable {
	table := &HostInfoTable{}
	//挂载主机列表
	table.HostsInfo = cache.Newcache("memory")
	//活跃主机列表
	table.ActiveHostList = NewHostList()
	//权重列表 倒序
	table.ActiveHostWeightList = list.New()
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
		} else {
			//如果活跃重新计算权重
			self.CalcHostWeight(hostInfo)
		}
		//更新活跃主机列表 或者添加主机列表,主机不活跃的时候从列表删除
		self.ActiveHostList.UpdateHostList(hostInfo)
		//对计算过后的活跃主机进行按照权重插入排序 unactive主机删除
		self.InsertSortHostWeight(*hostInfo)
		self.DumpSortedByWeightInfo()
	}
}

//主机是否在排序列表中
//根据ip对比
func (self *HostInfoTable) HostInfoSortList(hostinfo HostInfo) (*list.Element, bool) {
	//遍历列表
	for e := self.ActiveHostWeightList.Front(); e != nil; e = e.Next() {
		//获取hostinfo
		hostItem := e.Value.(HostInfo)
		if hostItem.Info.IP == hostinfo.Info.IP {
			return e, true
		}
	}
	return nil, false
}

//根据权重对主机进行排序  只有权重改变的主机才进行排序
func (self *HostInfoTable) InsertSortHostWeight(hostinfo HostInfo) {
	fmt.Println("对服务器进行权重排序...........")
	//列表空
	if self.ActiveHostWeightList.Len() == 0 && hostinfo.Status == ACTIVE {
		self.ActiveHostWeightList.PushFront(hostinfo)
	} else if self.ActiveHostWeightList.Len() > 0 {
		//如果在主机列表 先摘除
		//可添加标志位在权重不变的情况下 不排序
		if e, ok := self.HostInfoSortList(hostinfo); ok {
			// 先移除服务器
			self.ActiveHostWeightList.Remove(e)
			//不活跃直接退出
			if hostinfo.Status == UNACTIVE {
				return
			}
			//同一台服务器 列表中只有一台
			if 0 == self.ActiveHostWeightList.Len() {
				self.ActiveHostWeightList.PushFront(hostinfo)
				return
			}
		}
		//遍历列表
		for e := self.ActiveHostWeightList.Front(); e != nil; e = e.Next() {
			//获取hostinfo
			hostItem := e.Value.(HostInfo)
			//降序排序 插入并返回
			if hostItem.Info.Weight <= hostinfo.Info.Weight {
				self.ActiveHostWeightList.InsertBefore(hostinfo, e)
				fmt.Println("SortListLen4:" + strconv.Itoa(self.ActiveHostWeightList.Len()))
				return
			}
		}
		//最后权重最低插入到最后
		self.ActiveHostWeightList.PushBack(hostinfo)
	}
}

//calc host weight
//根据服务器负载状态  计算 权重
// cpu  0 1 2 3
// load 0 1 2 3
// mem  0 1 2 3
//
func (self *HostInfoTable) CalcHostWeight(hostInfo *HostInfo) {
	hostWeight := 0
	//cpu percent
	cpuPercent := 0.0
	cpuPercents := hostInfo.Info.CpuPercent
	for _, v := range cpuPercents {
		cpuPercent += v
	}
	cpuPercent = (cpuPercent / float64(len(cpuPercents)))
	hostWeight += 3 - int(int(cpuPercent)/30)
	//mem used
	mem := hostInfo.Info.VM.UsedPercent
	hostWeight += 3 - int(int(mem)/30)
	//net connections
	//网络连接暂时预留
	nc := hostInfo.Info.NC.AllConnectCount
	//load average
	load := hostInfo.Info.LD.Load1
	//good load
	goodLoad := float64(hostInfo.Info.CpuNums) * 0.9
	if load > 0.0 && load < (goodLoad/3) {
		hostWeight += 3
	} else if load >= (goodLoad/3) && load < (goodLoad*2/3) {
		hostWeight += 2
	} else {
		hostWeight += 1
	}

	hostInfo.Info.Weight = hostWeight
	fmt.Println("Calc Host Weight.......")
	fmt.Println("Mem Used:", mem)
	fmt.Println("Cpu Used:", cpuPercent)
	fmt.Println("Load:", load)
	fmt.Println("Good Load:", goodLoad)
	fmt.Println("Net Connection:", nc)
	fmt.Println("Host Weight:", hostInfo.Info.Weight)
	//设置权重

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

//dump挂载的服务器列表
func (self *HostInfoTable) DumpInfo() {
	mapChan := *self.HostsInfo.GetMemory().GetData()
	for k, v := range mapChan {
		fmt.Println(k)
		bts, _ := json.MarshalIndent(v, "", " ")
		fmt.Println(string(bts))
	}
}

//打印根据权重排序后的服务器信息
func (self *HostInfoTable) DumpSortedByWeightInfo() {
	for e := self.ActiveHostWeightList.Front(); e != nil; e = e.Next() {
		bts, _ := json.MarshalIndent(e.Value, "", " ")
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
