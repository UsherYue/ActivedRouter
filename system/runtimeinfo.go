//获取系统运行状态信息
package system

import (
	"encoding/json"
	"runtime"
	"time"

	"ActivedRouter/gopsutil/cpu"
	"ActivedRouter/gopsutil/disk"
	"ActivedRouter/gopsutil/host"
	"ActivedRouter/gopsutil/load"
	"ActivedRouter/gopsutil/mem"
	"ActivedRouter/gopsutil/net"
)

//网络信息
//各种tcp网络状态
type TCPNetInfo struct {
	ClosedWaitCount int `json:"CLOSED_WAIT"`
	//ClosedCount      int `json:"CLOSED"`
	ListenCount      int `json:"LISTEN"`
	EstablishedCount int `json:"ESTABLISH"`
	FinWait2Count    int `json:"FIN_WAIT_2"`
	FinWait1Count    int `json:"FIN_WAIT_1"`
	//ClosingCount     int `json:"CLOSING"`
	SynSentCount     int `json:"SYN_SENT"`
	SynReceivedCount int `json:"SYN_RECV"`
	TimeWaitCount    int `json:"TIME_WAIT"`
	//	LastAckCount     int `json:"LAST_ACK"`
	AllConnectCount int `json:"ALLCOUNT"`
}

//系统信息定义
type SystemInfo struct {
	VM      *mem.VirtualMemoryStat `json:"VM"`      //虚拟内存
	LD      *load.LoadAvgStat      `json:"LD"`      //load average
	DISK    *disk.DiskUsageStat    `json:"DISK"`    //dis
	HOST    *host.HostInfoStat     `json:"HOST"`    //host
	Cluster string                 `json:"Cluster"` //集群分组
	RunMode string                 `json:"RunMode"` //runmode
	Domain  string                 `json:"Domain"`  //domain
	IP      string                 `json:"IP"`
	CpuNums int                    `json:"CpuNums"` //cpu number
	Weight  int                    ///host weight
	//CPUS     []cpu.CPUInfoStat      `json:"CPUS"`     //cpu
	//CPUTIMES []cpu.CPUTimesStat     `json:"CPUTIMES"` //cpu times
	//SM *mem.SwapMemoryStat    `json:"SM"` //交换内存
	NC         *TCPNetInfo `json:"NC"`          //网络
	CpuPercent []float64   `json:"CPUPERCENTS"` //cpu percent
}

//获取系统信息 返回json
func SysInfo(runMode, cluster string) string {
	//内存
	virtualMem, _ := mem.VirtualMemory()
	//交换内存
	//swapMem, _ := mem.SwapMemory()
	//load
	loadAvg, _ := load.LoadAvg()
	info := &SystemInfo{}
	info.CpuNums = runtime.NumCPU()
	//runmode
	info.RunMode = runMode
	//vm
	info.VM = virtualMem
	//load
	info.LD = loadAvg
	//disk
	info.DISK, _ = disk.DiskUsage("/")
	//host info
	info.HOST, _ = host.HostInfo()
	//cluster
	info.Cluster = cluster

	//最近5s的状态
	cpuPercent, _ := cpu.CPUPercent(1*time.Second, true)
	info.CpuPercent = cpuPercent
	nc, _ := net.NetConnections("tcp4")
	tcpNc := TCPNetInfo{}
	//所有网络链接
	tcpNc.AllConnectCount = len(nc)
	info.NC = &tcpNc
	//ESTABLISHED  CLOSE_WAIT  LISTEN
	for _, ncItem := range nc {
		switch ncItem.Status {
		case "ESTABLISHED":
			{
				info.NC.EstablishedCount++
			}
		case "CLOSE_WAIT":
			{
				info.NC.ClosedWaitCount++
			}
		case "LISTEN":
			{
				info.NC.ListenCount++
			}
		case "TIME_WAIT":
			{
				info.NC.TimeWaitCount++
			}
		case "FIN_WAIT_1":
			{
				info.NC.FinWait1Count++
			}
		case "FIN_WAIT_2":
			{
				info.NC.FinWait2Count++
			}
		case "SYN_SENT":
			{
				info.NC.SynSentCount++
			}
		case "SYN_RECV":
			{
				info.NC.SynReceivedCount++
			}
		}
		info.NC.AllConnectCount++
	}
	//testNc()
	bts, _ := json.MarshalIndent(info, "", " ")
	return string(bts)
}

func testNc() {
	//nc
	//	nc, _ := net.NetConnections("tcp4")
	//	bts1, _ := json.MarshalIndent(nc, "", " ")
	//fmt.Println(len(nc))
	//fmt.Println(string(bts1))
}

//转换成本地结构体
func DecodeSysinfo(info string) (*SystemInfo, error) {
	sysinfo := &SystemInfo{}
	err := json.Unmarshal([]byte(info), sysinfo)
	if err != nil {
		return nil, err
	}
	return sysinfo, nil
}
