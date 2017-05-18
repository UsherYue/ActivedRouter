//ActivedRouter
//Author:usher.yue
//Amail:usher.yue@gmail.com
//TencentQQ:4223665

package boot

import (
	"ActivedRouter/boot/cmdline"
	"ActivedRouter/boot/config"
	"ActivedRouter/netservice"
)

func init() {
	//parse cmd line
	if cmdline.ParseCmdline() {
		//parse config file
		config.ParseConfigfile()
		//	start network
		netservice.StartNetworkService()
	}
}
