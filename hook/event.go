package hook

import (
	"ActivedRouter/cache"
	"fmt"
	"os/exec"
)

//钩子事件类型
const (
	DISK_EVENT = iota
	CPU_EVENT
	MEM_EVENT
	STATUS_EVENT
	LOAD_EVENT
)

//事件对象映射
var DefaultEventMaps = map[string]int{"disk": DISK_EVENT,
	"cpu":    CPU_EVENT,
	"mem":    MEM_EVENT,
	"status": STATUS_EVENT,
	"load":   LOAD_EVENT,
}

//event
type Event struct {
	EventType      int
	EventTarget    string
	EventAttr      string
	EventCondition map[string]string
	EventCallback  string
	EventHostIP    string
}

//new event

func NewEvent() *Event {
	event := &Event{}
	event.EventCondition = make(map[string]string)
	return event
}

//callback run
func (self *Event) ExecCallback() (string, error) {
	cmd := exec.Command("/bin/sh", "-c", self.EventCallback)
	if bts, err := cmd.Output(); err != nil {
		return "", err
	} else {
		fmt.Println(self.EventHostIP+"触发"+self.EventTarget+"事件－－－－－－－callback执行结果:", string(bts))
		return string(bts), nil
	}
}

//事件列表
type EventQueue struct {
	EventCache cache.Cache
	//email监控
	EmailOpen bool
	EmailUser string
	EmailPwd  string
	SmtpHost  string
	EmailTo   string
}

//new  memory event list
func NewEventQueue() *EventQueue {
	eventQueue := &EventQueue{}
	eventQueue.EventCache = cache.Newcache("memory")
	return eventQueue
}

//push event
func (self *EventQueue) PushEvent(k string, v interface{}) {
	self.EventCache.Set(k, v)
}

//get event
func (self *EventQueue) GetEvent(k string) []*Event {
	mp := *self.EventCache.GetMemory().GetData()
	if eventList, ok := mp[k]; ok {
		if ret, ok := eventList.([]*Event); ok {
			return ret
		}
		return nil
	}
	return nil
}

//events
func (self *EventQueue) GetEvents() map[string]interface{} {
	return *GEventQueue.EventCache.GetMemory().GetData()
}
