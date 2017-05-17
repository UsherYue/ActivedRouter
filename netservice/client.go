package netservice

import (
	"log"
	"net"
	"time"

	"ActivedRouter/global"
	"ActivedRouter/system"
	"ActivedRouter/tools"
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

//create client agent
func NewClient(host, port string) *Client {
	return &Client{Host: host, Port: port, TaskFlag: make(chan bool, 0), ConnectFlag: make(chan bool, 0), Closed: false}
}

//connect to server
func (self *Client) ConnectToServer(addr string) {
	//connect  time out 5s
	defer func() {
		if res := recover(); res != nil {
			log.Println("connect to router" + self.Host + ":" + self.Port + " server error!")
			//Re-connect when can't connect to router server
			self.ConnectFlag <- true
		}
	}()
	conn, _ := net.DialTimeout("tcp", addr, time.Second*5)
	self.ConnSocket = conn
	defer conn.Close()

	//conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	//conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	log.Println("The connection to the remote routing server was successful!")
	//Report the system status
	t1 := time.NewTimer(time.Second * HEARTBEAT_INTERVAL)
	for {
		select {
		case <-t1.C:
			{
				//Control to send data
				t1.Reset(time.Second * HEARTBEAT_INTERVAL)
				info := system.SysInfo(global.RunMode, global.Cluster, global.Domain)
				_, err := self.ConnSocket.Write(tools.Base64Encode([]byte(info)))
				//attempt connecting to router server until the client agent connect success
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

//connect to remote routing server
func (self *Client) Run() {
	log.Printf("Connecting to remote routing server, destination address %s:%s........\n", self.Host, self.Port)
	addr := ""
	if self.Host == "*" {
		addr = ":" + self.Port
	} else {
		addr = self.Host + ":" + self.Port
	}
	//Connect to multiple routing servers
	go self.ConnectToServer(addr)
	go func() {
		for {
			//2s try to re-link once
			time.Sleep(time.Second * 2)
			<-self.ConnectFlag
			if !self.Closed {
				log.Println("Attempt Connecting" + addr)
				go self.ConnectToServer(addr)
			} else {
				break
			}
		}
	}()
	//stop task
	self.TaskFlag <- true
}

//Disconnnect
func (self *Client) Disconnect() {
	self.Closed = true
	//Stop all connections before stopping the server
	self.ConnSocket.Close()
	//Send Close Message Exit Close the task
	<-self.TaskFlag
}
