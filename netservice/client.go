package netservice

import (
	"ActivedRouter/global"
	"ActivedRouter/system"
	"ActivedRouter/tools"
	"log"
	"net"
	"time"
)

type Client struct {
	Host        string
	Port        string
	TaskFlag    chan bool //syn channel
	ConnectFlag chan bool //connect syn channel
	Closed      bool
	ConnSocket  net.Conn
}

const (
	HEARTBEAT_INTERVAL = 5
)

//创建http服务
func NewClient(host, port string) *Client {
	return &Client{Host: host, Port: port, TaskFlag: make(chan bool, 0), ConnectFlag: make(chan bool, 0), Closed: false}
}

//connect to server
func (self *Client) ConnectToServer(addr string) {
	//connect  time out 5s
	defer func() {
		if res := recover(); res != nil {
			log.Println("connect to router" + self.Host + ":" + self.Port + " server error!")
			//链接失败尝试重复
			self.ConnectFlag <- true
		}
	}()
	conn, _ := net.DialTimeout("tcp", addr, time.Second*5)
	self.ConnSocket = conn
	defer conn.Close()

	//短连接设置
	//conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	//conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	log.Println("连接远程路由服务器成功!")
	//定时获取系统状态进行上报
	t1 := time.NewTimer(time.Second * HEARTBEAT_INTERVAL)
	for {
		select {
		case <-t1.C:
			{
				//控制发送数据
				//ch <- true
				t1.Reset(time.Second * HEARTBEAT_INTERVAL)
				info := system.SysInfo(global.RunMode, global.Cluster, global.Domain)
				_, err := self.ConnSocket.Write([]byte(tools.Base64Encode([]byte(info))))
				//如果断开连接重复连接 直到连接到路由服务器为止
				if err != nil && !self.Closed {
					conn, err := net.DialTimeout("tcp", addr, time.Second*5)
					if err == nil {
						self.ConnSocket = conn
					}
				}
			}
		}
	}
}

//connect to remote router server
func (self *Client) Run() {
	log.Printf("正在连接远程路由服务器,目标地址%s:%s........\n", self.Host, self.Port)
	addr := ""
	if self.Host == "*" {
		addr = ":" + self.Port
	} else {
		addr = self.Host + ":" + self.Port
	}
	//此处应该是连接到多个路由服务器
	go self.ConnectToServer(addr)
	//多次尝试
	go func() {
		for {
			//2s尝试一次重新链接
			time.Sleep(time.Second * 2)
			<-self.ConnectFlag
			if !self.Closed {
				log.Println("尝试链接" + addr)
				go self.ConnectToServer(addr)
			} else {
				break
			}
		}
	}()
	//stop task
	self.TaskFlag <- true
}

//停止服务器
func (self *Client) Disconnect() {
	self.Closed = true
	//停止服务器之前先关闭所有连接
	self.ConnSocket.Close()
	//发送关闭消息退出关闭任务
	<-self.TaskFlag
}
