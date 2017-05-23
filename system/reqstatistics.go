package system

import (
	"sync"

	"ActivedRouter/tools"
)

//http 反向代理请求统计
type HttpProxyStatistics struct {
	Timestamp    int64 //时间戳
	RequestCount int64 //all请求次数
}

//http请求分析
type StatisticsMap map[string][]*HttpProxyStatistics

type SysHttpStatistics struct {
	//统计列表
	statistic StatisticsMap
	//当前的节点
	currentNode map[string]*HttpProxyStatistics
	//rw lock
	mutexUpdate *sync.RWMutex
}

//创建一个反向代理统计对象
func newHttpProxyStatistics(timestamp, reqCount int64) *HttpProxyStatistics {
	return &HttpProxyStatistics{Timestamp: timestamp, RequestCount: reqCount}
}

//CREATE SYSHTTPSTATISTICS
func NewSysHttpStatistics() *SysHttpStatistics {
	return &SysHttpStatistics{statistic: make(StatisticsMap),
		currentNode: make(map[string]*HttpProxyStatistics), mutexUpdate: &sync.RWMutex{},
	}
}

//获取统计列表
func (self *SysHttpStatistics) GetStatisticsList() StatisticsMap {
	return self.statistic
}

//更新集群请求统计
//cluster 集群名称
//updateType  0 分别代表请求更新,只增加请求次数     1 时间段更新,添加新的时间段统计。
func (self *SysHttpStatistics) UpdateClusterStatistics(cluster string, updateType int) {
	self.mutexUpdate.Lock()
	//创建统计对象 当该集群统计列表不存在的时候
	if _, ok := self.statistic[cluster]; !ok && updateType == 0 {
		dataTool := tools.DateTool{}
		timestamp := dataTool.CurrentUnixTimestamp()
		var reqCount int64 = 0
		if updateType == 0 {
			reqCount = 1
		}
		self.statistic[cluster] = []*HttpProxyStatistics{newHttpProxyStatistics(timestamp, reqCount)}
		//unlock
		self.mutexUpdate.Unlock()
		return
	}
	//定时增加统计对象边界
	if updateType == 1 {
		dataTool := tools.DateTool{}
		timestamp := dataTool.CurrentUnixTimestamp()
		for k, _ := range self.statistic {
			//长度为100的线性表
			tableLen := len(self.statistic[k])
			//移除线性表投诉
			if tableLen > 100 {
				self.statistic[k] = self.statistic[k][1:]
			}
			//追加
			self.statistic[k] = append(self.statistic[k], newHttpProxyStatistics(timestamp, 0))
		}
	} else {
		//取出最后一个统计索引
		lastIndex := len(self.statistic[cluster]) - 1
		//增加指定集群请求次数
		self.statistic[cluster][lastIndex].RequestCount++
	}
	self.mutexUpdate.Unlock()
}

//存储当天的统计数据到json日志 或者数据库中
func (self *SysHttpStatistics) SaveDataLog() {

}

//清空当天的请求统计日志 并且清空内存
func (self *SysHttpStatistics) ResetData() {

}
