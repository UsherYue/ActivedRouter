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
	CHECK_ACTIVE_INTERVAL   = 2
	DISPATCH_EVENT_INTERVAL = 10
	BUFFER_SIZE             = 100 * 1024 * 1024
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
			global.GHostInfoTable.UpdateHostTable(addrs[0], data)
			srvmode, _ := global.ConfigMap["srvmode"]
			switch srvmode {
			case "moniter":
				{
					log.Println("-------event begin------------")
					hook.DispatchEvent()
					log.Println("-------event end------------")
				}
			}
		}
	}
}

//停止服务器
func (self *Server) StopServer() {
	//停止服务器之前先关闭所有连接
	//发送关闭消息
	<-self.TaskFlag
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
	//监控
	go func() {
		t1 := time.NewTimer(time.Second * CHECK_ACTIVE_INTERVAL)
		for {
			select {
			case <-t1.C:
				{
					t1.Reset(time.Second * CHECK_ACTIVE_INTERVAL)
					//更新服务器状态
					global.GHostInfoTable.UpdateHostStatus()
				}
			}
		}
	}()
	//dispatch event
	//	dispather := func() {
	//		t1 := time.NewTimer(time.Second * DISPATCH_EVENT_INTERVAL)
	//		for {
	//			select {
	//			case <-t1.C:
	//				{
	//					log.Println("-------event begin------------")
	//					hook.DispatchEvent()
	//					log.Println("-------event end------------")
	//					t1.Reset(time.Second * DISPATCH_EVENT_INTERVAL)
	//				}
	//			}
	//		}
	//	}
	//moniter 模式下分发事件
	//	srvmode, _ := global.ConfigMap["srvmode"]
	//	switch srvmode {
	//	case "moniter":
	//		{
	//			go dispather()
	//		}
	//	}
	//等待任务结束
	self.TaskFlag <- false
}
