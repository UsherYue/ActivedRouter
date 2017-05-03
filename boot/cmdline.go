package boot

import (
	"flag"
	"log"
	"strings"

	. "ActivedRouter/global"
)

func parseCmdline() {
	runmode := flag.String("runmode", "", "runmode must be server or client or reserveproxy")
	flag.Parse()
	//if not set run mode
	if *runmode == "" {
		log.Fatalln(UsageTemplate)
	} else {
		if strings.ToLower(*runmode) != ServerMode &&
			strings.ToLower(*runmode) != ClientMode &&
			strings.ToLower(*runmode) != ReserveProxyMode {
			log.Fatalln(UsageRunmodeError)
		} else {
			//set run mode
			RunMode = strings.ToLower(*runmode)
		}
	}
}
