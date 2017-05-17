package netservice

import (
	"log"
	"net"
	"strings"
	"time"

	"ActivedRouter/global"
	"ActivedRouter/hook"
	"ActivedRouter/system"
	"ActivedRouter/tools"
)

const (
	CHECK_INTERCAL          = 1
	CHECK_ACTIVE_INTERVAL   = 2                 // check  active interval
	CHECK_ROUTER_INTERVAL   = 5                 //Periodically check the routing server status
	DISPATCH_EVENT_INTERVAL = 5                 //check event dispatch  interval
	BUFFER_SIZE             = 100 * 1024 * 1024 //Maximum size of cache
)

//Router server infomational encapsulation
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

//Data Receive
func (self *Server) OnDataRecv(c net.Conn) {
	//return
	log.Printf("accept connect from %s\n", c.RemoteAddr().String())
	defer c.Close()
	buffer := make([]byte, BUFFER_SIZE)
	for {
		//Get client heartbeat feedback
		if n, err := c.Read(buffer); err != nil {
			//The server shuts down the client-to-server connection when
			//the client actively shuts down the connection
			if err.Error() == "EOF" {
				c.Close()
			}
			//Repair cpu overload bug
			return
		} else {
			if n > 0 {
				//Deserialize the data and perform statistical analysis
				///Note.......
				//There may be sticky packets of bug
				decodeData, _ := tools.Base64Decode(buffer[:n])
				data, err := system.DecodeSysinfo(string(decodeData))
				//Parsing errors but not handling
				if err != nil {
					log.Println(err.Error())
					//Close connection to server when json decode  error  occured
					c.Close()
					return
				}
				//get remote host ip
				addrs := strings.Split(c.RemoteAddr().String(), ":")
				//update host status
				data.IP = addrs[0]
				global.GHostInfoTable.UpdateHostTable(addrs[0], data)
			}
		}
	}
}

//Send Service Stop Message We Shoule Close all connections before stopping the server。
func (self *Server) StopServer() {
	<-self.TaskFlag
}

//Timing monitoring of routing server information
func (self *Server) checkRouterInfo() {
	timerRouterInfo := time.NewTimer(time.Second * CHECK_ROUTER_INTERVAL)
	for {
		select {
		case <-timerRouterInfo.C:
			{
				routerInfo := system.SysInfo(global.RunMode, "ActivedRouterInfo", "")
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

//monitoring client
func (self *Server) moniterClient() {
	timerCheckActive := time.NewTimer(time.Second * CHECK_ACTIVE_INTERVAL)
	for {
		select {
		case <-timerCheckActive.C:
			{
				timerCheckActive.Reset(time.Second * CHECK_ACTIVE_INTERVAL)
				//update host status
				global.GHostInfoTable.UpdateHostStatus()
			}
		}
	}
}

//run router server
func (self *Server) Run() {
	log.Printf("Begin Running Router Service,%s:%s........\n", self.Host, self.Port)
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
				log.Println(err)
			}
			//Data Recv
			go self.OnDataRecv(conn)
		}
	}()
	//Run Monitor client
	go self.moniterClient()
	//Check the server status regularly
	go self.checkRouterInfo()
	//dispatch monitor event
	go self.dispatcher()
	//wait for stop message
	self.TaskFlag <- false
}
