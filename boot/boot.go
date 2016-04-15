package boot

import (
	"ActivedRouter/hook"
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
	//hook script
	hook.ParseHookScript("config/hook.json")
	//start network
	netservice.StartNetworkService()

}
