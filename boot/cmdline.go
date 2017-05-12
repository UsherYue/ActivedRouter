package boot

import (
	"flag"
	"fmt"
	"log"
	"strings"

	. "ActivedRouter/global"
)

func parseCmdline() bool {
	runmode := flag.String("runmode", "", UsageTemplate)
	flag.Parse()
	//if not set run mode
	if *runmode == "" {
		fmt.Printf(DescTemplate)
		fmt.Println("\033[1m", UsageTemplate, "\033[0m")
		fmt.Println(TheHelpTemplate)
		return false
	} else {
		if strings.ToLower(*runmode) != ServerMode &&
			strings.ToLower(*runmode) != ClientMode &&
			strings.ToLower(*runmode) != ReserveProxyMode {
			log.Println(UsageRunmodeError)
			return false
		} else {
			//set run mode
			RunMode = strings.ToLower(*runmode)
			return true
		}
	}
}
