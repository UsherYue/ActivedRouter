package netservice

import (
	"ActivedRouter/global"
	"ActivedRouter/hook"
	"ActivedRouter/system"
	"ActivedRouter/tools"
	"log"
	"net"
	"strings"
	"time"
)

//服务器检测时间5s一次
//后期可以配置脚本化
const (
	CHECK_INTERCAL          = 1
	CHECK_ACTIVE_INTERVAL   = 2                 //活跃检测周期
	CHECK_ROUTER_INTERVAL   = 5                 //定期检查路由服务器状态
	DISPATCH_EVENT_INTERVAL = 5                 //事件分发检测
	BUFFER_SIZE             = 100 * 1024 * 1024 //‘缓存buffer大小
)

//路由服务器相关封装
type Server struct {
	Host         string
	Port         string
	TaskFlag     chan bool    //syn channel
	ListenSocket net.Listener //conn
}

//create server
func NewServer(host, port string) *Server {
	return &Server{Host: host, Port: port, TaskFlag: make(chan bool, 0)}
}

//数据接收
func (self *Server) OnDataRecv(c net.Conn) {
	log.Printf("accept connect from %s\n", c.RemoteAddr().String())
	defer c.Close()
	buffer := make([]byte, BUFFER_SIZE)
	for {
		//获取客服端心跳反馈
		n, _ := c.Read(buffer)
		if n > 0 {
			//反序列化数据并且进行统计分析
			decodeData, _ := tools.Base64Decode(buffer[:n])
			data, err := system.DecodeSysinfo(string(decodeData))
			//解析错误不处理
			if err != nil {
				log.Println(err.Error())
				continue
			}
			//获取远程ip
			addrs := strings.Split(c.RemoteAddr().String(), ":")
			//更新服务器状态
			//要做到 thread safe
			//更新服务器列表 如果不存在那么添加到服务器列表
			data.IP = addrs[0]
			global.GHostInfoTable.UpdateHostTable(addrs[0], data)
		}
	}
}

//停止服务器
func (self *Server) StopServer() {
	//停止服务器之前先关闭所有连接
	//发送关闭消息
	<-self.TaskFlag
}

//定时监测路由服务器信息
func (self *Server) checkRouterInfo() {
	timerRouterInfo := time.NewTimer(time.Second * CHECK_ROUTER_INTERVAL)
	for {
		select {
		case <-timerRouterInfo.C:
			{
				routerInfo := system.SysInfo("Router", "")
				global.SetRouterInfo(routerInfo)
				timerRouterInfo.Reset(time.Second * CHECK_ROUTER_INTERVAL)
			}
		}
	}

}

//event dispatch
func (self *Server) dispatcher() {
	closureFunc := func() {
		timerDispathcEvent := time.NewTimer(time.Second * DISPATCH_EVENT_INTERVAL)
		for {
			select {
			case <-timerDispathcEvent.C:
				{
					log.Println("-------event begin------------")
					hook.DispatchEvent()
					log.Println("-------event end------------")
					timerDispathcEvent.Reset(time.Second * DISPATCH_EVENT_INTERVAL)
				}
			}
		}
	}
	srvmode, _ := global.ConfigMap["srvmode"]
	switch srvmode {
	case "moniter":
		{
			go closureFunc()
		}
	}
}

//监控client
func (self *Server) moniterClient() {
	timerCheckActive := time.NewTimer(time.Second * CHECK_ACTIVE_INTERVAL)
	for {
		select {
		case <-timerCheckActive.C:
			{
				timerCheckActive.Reset(time.Second * CHECK_ACTIVE_INTERVAL)
				//更新服务器状态
				global.GHostInfoTable.UpdateHostStatus()
			}
		}
	}
}

//run router server
func (self *Server) Run() {
	log.Printf("开始启动路由服务器服务,%s:%s........\n", self.Host, self.Port)
	addr := ""
	if self.Host == "*" {
		addr = ":" + self.Port
	} else {
		addr = self.Host + ":" + self.Port
	}
	//listen
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	//save listen socket
	self.ListenSocket = l
	//accept
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				//提示错误但不退出
				log.Println(err)
			}
			//数据处理
			go self.OnDataRecv(conn)
		}
	}()
	//监控client
	go self.moniterClient()
	//定期收集路由服务器信息
	go self.checkRouterInfo()
	//moniter 模式下分发事件
	go self.dispatcher()
	//等待任务结束
	self.TaskFlag <- false
}
