package boot

import (
	"bytes"
	"flag"
	"html/template"
	"log"
	"os"
	"strings"

	. "ActivedRouter/global"
)

///解析命令行参数
func parseCmdline() {
	runmode := flag.String("runmode", "", "runmode must be server or client or reserveproxy")
	flag.Parse()
	//如果没设置运行模式
	if *runmode == "" {
		log.Println(UsageTemplate)
		os.Exit(0)
	} else {
		if strings.ToLower(*runmode) != ServerMode &&
			strings.ToLower(*runmode) != ClientMode &&
			strings.ToLower(*runmode) != ReserveProxyMode &&
			strings.ToLower(*runmode) != MixMode {
			t, _ := template.New("info").Parse(UsageRunmodeTemplate)
			buffer := &bytes.Buffer{}
			t.Execute(buffer, struct{ Msg string }{Msg: "runmode参数错误,参考 ActiveRouter --runmode=Client或Reserveproxy,Server"})
			log.Println(string(buffer.Bytes()))
		} else {
			//set run mode
			RunMode = strings.ToLower(*runmode)
		}
	}
}
