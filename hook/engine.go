package hook

import (
	"ActivedRouter/global"
	_ "ActivedRouter/system"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
			log.Println(string(bts))
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
		event.ExecCallback()
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
		event.ExecCallback()
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
		event.ExecCallback()
	}
}

//处理cpu event
func processCPUEvent(hostip string, event *Event) {
	log.Println("cpu event")

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

}

//dispatch event
//设计事件分发机制
func DispatchEvent() {
	mapData := GEventQueue.GetEvents()
	for host, eventList := range mapData {
		//跳过没有挂载的主机事件
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
