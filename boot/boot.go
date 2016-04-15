package boot

import (
	"ActivedRouter/netservice"
)

/**
* 初始化
 */
func init() {
	//parse cmd line
	parseCmdline()
	//parse config file
	parseConfigfile()
	//start network
	netservice.StartNetworkService()

}
