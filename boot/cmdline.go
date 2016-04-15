package boot

import (
	"ActivedRouter/global"
	"bytes"
	"flag"
	"html/template"
	"log"
	"os"
	"strings"
)

//使用方法模板
var usageTemplate = `
ActiveRouter 是一个简单的基于路由分发和反向代理的负载均衡监控服务.
Author:usher.yue
Email:ushe.yue@gmail.com

ActiveRouter Desceiption:
	mode运行模式可以是 server或者client,分别代表路由服务器和客户端 proxy代表反向代理服务器。
	服务器模式下加载server.ini。
	客户端模式下记加载client.ini读取相关配置信息。	
The commands are:
	ActiveRouter --runmode=server或client	
The Help:
    ActiveRouter --help or  -h or -help
`

//运行模式文档
var usageRunmodeTemplate = `
	{{.Msg}}
`

///解析命令行参数
func parseCmdline() {
	runmode := flag.String("runmode", "", "runmode must be server or client or proxy")
	//如果存在--help那么直接退出进程
	flag.Parse()
	if *runmode == "" {
		log.Println(usageTemplate)
		goto EXIT
	} else {
		if strings.ToLower(*runmode) != "server" && strings.ToLower(*runmode) != "client" && strings.ToLower(*runmode) != "proxy" {
			t, _ := template.New("info").Parse(usageRunmodeTemplate)
			buffer := &bytes.Buffer{}
			t.Execute(buffer, struct{ Msg string }{Msg: "runmode参数错误,参考 ActiveRouter --runmode=server或client proxy"})
			//获取buffer内容
			log.Println(string(buffer.Bytes()))
			goto EXIT
		}
	}
	global.RunMode = *runmode
	if strings.ToLower(*runmode) == "client" {
		log.Println("ActivedRouter正在启动client模式")
	} else if strings.ToLower(*runmode) == "server" {
		log.Println("ActivedRouter正在启动server模式")
	} else if strings.ToLower(*runmode) == "proxy" {
		log.Println("ActivedRouter正在启动reserve proxy模式")
	}
	return
EXIT:
	os.Exit(0)
}
