package netservice

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"ActivedRouter/global"
	. "ActivedRouter/netservice/packet"
	"ActivedRouter/system"
)

//client config
type ClientConfig struct {
	ClusterName string   `json:"cluster"`
	RouterList  []string `json:"router_list"`
}

var ClientConfigData ClientConfig

type Client struct {
	Host        string
	Port        string
	TaskFlag    chan bool //syn channel
	ConnectFlag chan bool //connect syn channel
	Closed      bool
	ConnSocket  net.Conn
}

//create client agent
func NewClient(host, port string) *Client {
	return &Client{Host: host, Port: port, TaskFlag: make(chan bool, 0), ConnectFlag: make(chan bool, 0), Closed: false}
}

/*
{
	"cluster":"测集群",
	"router_list":[
		"127.0.0.1:8888",
		"172.16.200.202:8888"
	]
}
load client json config file*/
func LoadClientConfig(configFilePath string) {
	file, err := os.Open(configFilePath)
	defer file.Close()
	if err != nil {
		log.Fatalln(err.Error())
	}
	if bts, err := ioutil.ReadAll(file); err != nil {
		log.Fatalln(err.Error())
	} else {
		if err := json.Unmarshal(bts, &ClientConfigData); err != nil {
			log.Fatalln("load client.json error")
		}
	}
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
	t1 := time.NewTimer(time.Second * global.CLIENT_HEARTBEAT_INTERVAL)
	for {
		select {
		case <-t1.C:
			{
				//Control to send data
				t1.Reset(time.Second * global.CLIENT_HEARTBEAT_INTERVAL)
				info := system.SysInfo(global.RunMode, ClientConfigData.ClusterName)
				//Encapsulate packets
				dataPackage := NewDefaultPacket([]byte(info)).Packet()
				//fmt.Println(tools.BytesToHexString(dataPackage))
				_, err := self.ConnSocket.Write(dataPackage)
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

//Connect to remote routing server
func (self *Client) Run() {
	log.Printf("Connecting to remote routing server, destination address %s:%s........\n", self.Host, self.Port)
	var addr string
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
