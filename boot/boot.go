//ActivedRouter
//Author:usher.yue
//Amail:usher.yue@gmail.com
//TencentQQ:4223665

package boot

import . "ActivedRouter/netservice"

/**
* 初始化
 */
func init() {
	//parse cmd line
	parseCmdline()
	//parse config file
	parseConfigfile()
	//	start network
	StartNetworkService()

}
