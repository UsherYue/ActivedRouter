//ActivedRouter
//Author:usher.yue
//Amail:usher.yue@gmail.com
//TencentQQ:4223665

package boot

import (
	. "ActivedRouter/netservice"
)

func init() {
	//parse cmd line
	if parseCmdline() {
		//parse config file
		parseConfigfile()
		//	start network
		StartNetworkService()
	}
}
