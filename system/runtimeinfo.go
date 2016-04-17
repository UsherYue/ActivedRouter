//获取系统运行状态信息
package system

import (
	"encoding/json"
	"runtime"
	"strings"

	"ActivedRouter/gopsutil/disk"
	"ActivedRouter/gopsutil/host"
	"ActivedRouter/gopsutil/load"
	"ActivedRouter/gopsutil/mem"
)

//系统信息定义
type SystemInfo struct {
	VM      *mem.VirtualMemoryStat `json:"VM"`      //虚拟内存
	LD      *load.LoadAvgStat      `json:"LD"`      //load average
	DISK    *disk.DiskUsageStat    `json:"DISK"`    //dis
	HOST    *host.HostInfoStat     `json:"HOST"`    //host
	Cluster string                 `json:"Cluster"` //集群分组
	Domain  string                 `json:"Domain"`  //domain
	IP      string                 `json:"IP"`
	CpuNums int                    `json:"CpuNums"` //cpu number
	Weight  int                    ///host weight
	//CPUS    []cpu.CPUInfoStat      `json:"CPUS"`    //cpu
	//CPUTIMES []cpu.CPUTimesStat     `json:"CPUTIMES"` //cpu times
	//SM *mem.SwapMemoryStat    `json:"SM"` //交换内存
	//NC       []net.NetConnectionStat `json:"NC"`       //网络
}

//获取系统信息 返回json
func SysInfo(cluster, domain string) string {
	//内存
	virtualMem, _ := mem.VirtualMemory()
	//交换内存
	//swapMem, _ := mem.SwapMemory()
	//load
	loadAvg, _ := load.LoadAvg()
	info := &SystemInfo{}
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
	//domain
	info.Domain = domain
	//info.CPUS, _ = cpu.CPUInfo()
	//cpu counts
	info.CpuNums = runtime.NumCPU()
	bts, _ := json.Marshal(info)
	return strings.TrimSpace(strings.Trim(strings.Trim(string(bts), "\n"), "\t"))
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
