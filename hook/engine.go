package hook

import (
	"ActivedRouter/global"
	_ "ActivedRouter/system"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"strings"
)

//钩子脚本
var _hookScript map[string]interface{} = nil

//event queue
var GEventQueue = NewEventQueue()

//syntax
var GScriptSyntax = NewDefaultSyntax()

//load HOOK script
func loadHookScript(routerFile string) {
	file, err := os.Open(routerFile)
	defer func() {
		file.Close()
	}()
	if err != nil {
		log.Fatalln(err.Error())
	}
	//reader
	if bts, err := ioutil.ReadAll(file); err != nil {
		log.Fatalln(err.Error())
	} else {
		var hookScript map[string]interface{}
		err := json.Unmarshal(bts, &hookScript)
		if err != nil {
			log.Fatalln(err.Error())
		} else {
			_hookScript = hookScript
		}
	}
}

//检测脚本语法
func checkScriptItem(scriptItem map[string]interface{}) {
	for k, _ := range scriptItem {
		if !GScriptSyntax.CheckSyntakKeyWords(k) {
			log.Fatalln("Script Syntax Error :", k, "Unknow Syntax")
		}
	}
}

//解析钩子脚本
func ParseHookScript(configfile string) {
	loadHookScript(configfile)
	//发送通知的smtp 和xxxxxxxxxx
	if v, ok := _hookScript["email_open"]; ok {
		mailOpen, ok := v.(string)
		if !ok {
			GEventQueue.EmailOpen = false
		} else {
			//是否开启email事件通知
			if mailOpen == "1" {
				mailTo := _hookScript["emailto"].(string)
				mailUsername := _hookScript["username"].(string)
				mailPwd := _hookScript["password"].(string)
				mailSmtpServer := _hookScript["smtp_server"].(string)
				GEventQueue.EmailOpen = true
				GEventQueue.EmailUser = mailUsername
				GEventQueue.EmailPwd = mailPwd
				GEventQueue.SmtpHost = mailSmtpServer
				GEventQueue.EmailTo = mailTo
			} else {
				GEventQueue.EmailOpen = false
			}
		}
	}
	scriptList := _hookScript["script"]
	eventList := scriptList.([]interface{})
	for _, event := range eventList {
		eventMap, _ := event.(map[string]interface{})
		host := ""
		// 事件对象列表
		var eventObjList []*Event
		for k, v := range eventMap {
			switch k {
			case "host":
				{
					host, _ = v.(string)
				}
			case "hookscript":
				{
					scriptItems, _ := v.([]interface{})
					for _, scriptItem := range scriptItems {
						subScriptItem, _ := scriptItem.(map[string]interface{})
						//检测脚本语法
						checkScriptItem(subScriptItem)
						eventItem := NewEvent()
						eventItem.EventHostIP = host
						for subK, subV := range subScriptItem {
							switch subK {
							case "attr":
								{
									eventItem.EventAttr, _ = subV.(string)
								}
							case "callback":
								{
									eventItem.EventCallback, _ = subV.(string)
								}
							case "eventtarget":
								{
									eventItem.EventTarget, _ = subV.(string)
									eventItem.EventType = DefaultEventMaps[eventItem.EventTarget]
								}
							default:
								{
									eventItem.EventCondition[subK] = subV.(string)
								}
							}
						}
						eventObjList = append(eventObjList, eventItem)
					}
				}
			}
		}
		GEventQueue.PushEvent(host, eventObjList)
	}
}

//发送邮件
func SendEmailTo(user, password, host, to, subject, body, mailtype string) error {
	hostInfo := strings.Split(host, ":")
	//AUTH
	auth := smtp.PlainAuth("", user, password, hostInfo[0])
	//MIME TYPE
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	//email body
	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

//generate notify content
func sendNotifyContent(subject, hostip, info, sysinfo string) {
	template, _ := template.New("notify").Parse(NotifyTemplate)
	buffer := &bytes.Buffer{}
	var data struct {
		IP      string
		INFO    string
		SYSINFO string
	}
	data.IP = hostip
	data.INFO = info
	data.SYSINFO = sysinfo
	template.Execute(buffer, &data)
	//发送通知邮件到管理员邮箱
	SendEmailTo(GEventQueue.EmailUser, GEventQueue.EmailPwd, GEventQueue.SmtpHost, GEventQueue.EmailTo, subject, string(buffer.Bytes()), "html")
}

//通告模板
var NotifyTemplate = `
<html>  
        <head>  
			<title>ActivedRouter邮件通知</title>
        </head>  
        <body >  
            尊敬的管理员:<br/>
			您的服务器收到如下阀值报警 <br/>
			IP:{{.IP}}<br/>
			报警信息:{{.INFO}}<br/>
			当前服务器状态:<br/>
			<pre>
				{{.SYSINFO}}
			</pre>
			<br/>
        </body>  
</html>  
`

//处理disk event
func processDiskEvent(hostip string, event *Event) {
	//	//获取服务器
	info := global.GHostInfoTable.GetHostInfo(hostip)
	//获取失败返回
	if info == nil {
		return
	}
	exprData := GScriptSyntax.GetExpt(event)
	var used float64
	switch event.EventAttr {
	case "used":
		{
			used = (info.Info.DISK.UsedPercent)
		}
	}
	//触发事件执行
	if GScriptSyntax.CheckFloadValue(exprData, used) {
		//event.ExecCallback()
		log.Println("Send Disk Email......")
		//发送警报
		tipinfo := fmt.Sprintf("%s%.2f%s", "服务器磁盘使用率:", used, "超出限制阀值")
		bts, _ := json.MarshalIndent(info, "", " ")
		sendNotifyContent("ActivedRouter Disk警报", hostip, tipinfo, string(bts))
	}
}

//处理mem event
func processMemEvent(hostip string, event *Event) {
	//获取服务器
	info := global.GHostInfoTable.GetHostInfo(hostip)
	//获取失败返回
	if info == nil {
		return
	}
	exprData := GScriptSyntax.GetExpt(event)
	var used float64
	switch event.EventAttr {
	case "used":
		{
			used = (info.Info.VM.UsedPercent)
		}
	}
	//触发事件执行
	if GScriptSyntax.CheckFloadValue(exprData, used) {
		//event.ExecCallback()
		log.Println("Send VM Email......")
		//发送警报
		tipinfo := fmt.Sprintf("%s%.2f%s", "服务器虚拟内存使用率:", used, "超出限制阀值")
		bts, _ := json.MarshalIndent(info, "", " ")
		sendNotifyContent("ActivedRouter VM警报", hostip, tipinfo, string(bts))
	}
}

//处理load event
func processLoadEvent(hostip string, event *Event) {
	//获取服务器
	info := global.GHostInfoTable.GetHostInfo(hostip)
	//获取失败返回
	if info == nil {
		return
	}
	exprData := GScriptSyntax.GetExpt(event)
	var load1, load5, load15 float64
	switch event.EventAttr {
	case "load":
		{
			load1 = (info.Info.LD.Load1)
			load5 = (info.Info.LD.Load5)
			load15 = (info.Info.LD.Load15)
		}
	}
	//触发事件执行
	if GScriptSyntax.CheckFloadValue(exprData, load1) ||
		GScriptSyntax.CheckFloadValue(exprData, load5) ||
		GScriptSyntax.CheckFloadValue(exprData, load15) {
		//event.ExecCallback()
		log.Println("Send load Email......")
		//发送警报
		tipinfo := fmt.Sprintf("%s%.2f,%.2f,%.2f%s", "服务器负载:", load1, load5, load15, " 超出限制阀值")
		bts, _ := json.MarshalIndent(info, "", " ")
		sendNotifyContent("ActivedRouter Load Average警报", hostip, tipinfo, string(bts))
	}
}

//处理cpu event
func processCPUEvent(hostip string, event *Event) {
	log.Println("cpu event")
	//获取服务器信息
	info := global.GHostInfoTable.GetHostInfo(hostip)
	//获取失败返回
	if info == nil {
		return
	}
	//表达式数据
	exprData := GScriptSyntax.GetExpt(event)
	var cpuPercent float64 = 0
	for _, p := range info.Info.CpuPercent {
		cpuPercent += p
	}
	//cpu 利用率
	cpuPercent = cpuPercent / float64(info.Info.CpuNums)
	//除法cpu报警事件
	if GScriptSyntax.CheckFloadValue(exprData, cpuPercent) {
		log.Println("Send cpu info Email......")
		//发送警报
		tipinfo := fmt.Sprintf("%s%.2f%s", "服务器CPU使用率:", cpuPercent, "超出限制阀值")
		bts, _ := json.MarshalIndent(info, "", " ")
		sendNotifyContent("ActivedRouter CPU使用率警报", hostip, tipinfo, string(bts))
	}

}

//处理status event
func processStatusEvent(hostip string, event *Event) {
	log.Println("status event")
	//获取服务器
	info := global.GHostInfoTable.GetHostInfo(hostip)
	//获取失败返回
	if info == nil {
		return
	}
	switch info.Status {
	case "active":
		{
			log.Println("Send  info active Email......")
		}
	case "unactive":
		{
			log.Println("Send  info unactive Email......")
			//发送警报
			tipinfo := fmt.Sprintf("%s", "服务器下线或者宕机!!!")
			bts, _ := json.MarshalIndent(info, "", " ")
			sendNotifyContent("ActivedRouter 服务器状态通告", hostip, tipinfo, string(bts))
		}
	}

}

//dispatch event
//设计事件分发机制
func DispatchEvent() {
	mapData := GEventQueue.GetEvents()
	for host, eventList := range mapData {
		//跳过没有挂载的主机事件
		//如果host 是* 那么标示对所有主机进行监控
		if info := global.GHostInfoTable.GetHostInfo(host); info == nil {
			log.Println("没有发现主机,忽略主机", host, "的钩子脚本!")
			continue
		} else {
			//触发事件
			eventArr, _ := eventList.([]*Event)
			for _, eventItem := range eventArr {
				switch eventItem.EventType {
				case DISK_EVENT:
					{
						processDiskEvent(host, eventItem)
					}
				case MEM_EVENT:
					{
						processMemEvent(host, eventItem)
					}
				case CPU_EVENT:
					{
						processCPUEvent(host, eventItem)
					}
				case STATUS_EVENT:
					{
						processStatusEvent(host, eventItem)
					}
				case LOAD_EVENT:
					{
						processLoadEvent(host, eventItem)
					}
				}
			}
		}

	}

}
