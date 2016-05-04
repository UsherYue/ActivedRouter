package system

//http 反向代理请求统计
type HttpProxyStatistics struct {
	Timestamp    int64 //时间戳
	RequestCount int64 //all请求次数
}

//http请求分析
type SysHttpStatistics map[string][]*HttpProxyStatistics

//创建一个反向代理统计对象
func newHttpProxyStatistics(timestamp, reqCount int64) *HttpProxyStatistics {
	return &HttpProxyStatistics{Timestamp: timestamp, RequestCount: reqCount}
}

//CREATE SYSHTTPSTATISTICS
func NewSysHttpStatistics() SysHttpStatistics {
	return make(SysHttpStatistics)
}

//更新集群请求统计
func (self SysHttpStatistics) UpdateClusterStatistics(cluster string, timestamp, reqCount int64) {
	//创建统计对象
	if _, ok := self[cluster]; !ok {
		self[cluster] = []*HttpProxyStatistics{newHttpProxyStatistics(timestamp, reqCount)}
		return
	}
	//如果存在那么直接更新统计数据
	//时间肯定是生序的
	self[cluster] = append(self[cluster], newHttpProxyStatistics(timestamp, reqCount))
}
