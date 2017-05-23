package netservice

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"

	"strings"
	"time"

	"ActivedRouter/global"
	"ActivedRouter/hook"
	. "ActivedRouter/netservice/packet"
	"ActivedRouter/system"
)

//server config
var ServerConfigData ServerConfig

//server config mapping struct
type ServerConfig struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	ServerMode string `json:"srvmode"`
	HttpHost   string `json:"httphost"`
	HttpPort   string `json:"httpport"`
}

//loade serverConfig
func LoadServerJsonConfig(config string) {
	if file, err := os.Open(config); err == nil {
		if bts, err := ioutil.ReadAll(file); err == nil {
			if json.Unmarshal(bts, &ServerConfigData) == nil {
				global.SrvMode = ServerConfigData.ServerMode
				return
			}
		} else {
			goto Exit
		}
	} else {
		goto Exit
	}
Exit:
	log.Fatalln("server config load error!")
}

//Router server infomational encapsulation
type Server struct {
	Host         string
	Port         string
	TaskFlag     chan bool    //syn channel
	ListenSocket net.Listener //conn
}

//create server
func NewServer() *Server {
	return &Server{TaskFlag: make(chan bool, 0)}
}

//Data Receive
func (self *Server) OnDataRecv(c net.Conn) {
	//return
	log.Printf("accept connect from %s\n", c.RemoteAddr().String())
	defer c.Close()
	buffer := make([]byte, global.SERVER_BUFFER_SIZE)
	//Declare a pipe for receiving unpacked data
	readerChannel := make(chan []byte, 1024)
	//Store truncated data
	remainBuffer := make([]byte, 0)
	//read unpackage data from buffered channel
	go func(reader chan []byte) {
		for {
			packageData := <-reader
			//decodeData, _ := tools.Base64Decode(packageData)
			data, err := system.DecodeSysinfo(string(packageData))
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
	}(readerChannel)
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
				//fix sticky bugs...
				remainBuffer = NewDefaultPacket(append(remainBuffer, buffer[:n]...)).UnPacket(readerChannel)
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
	timerRouterInfo := time.NewTimer(time.Second * global.SERVER_CHECK_ROUTER_INTERVAL)
	for {
		select {
		case <-timerRouterInfo.C:
			{
				routerInfo := system.SysInfo(global.RunMode, "ActivedRouterInfo")
				global.SetRouterSysInfo(routerInfo)
				timerRouterInfo.Reset(time.Second * global.SERVER_CHECK_ROUTER_INTERVAL)
			}
		}
	}
}

//event dispatch
func (self *Server) dispatcher() {
	closureFunc := func() {
		timerDispathcEvent := time.NewTimer(time.Second * global.SERVER_DISPATCH_EVENT_INTERVAL)
		for {
			select {
			case <-timerDispathcEvent.C:
				{
					log.Println("-------event begin------------")
					hook.DispatchEvent()
					log.Println("-------event end------------")
					timerDispathcEvent.Reset(time.Second * global.SERVER_DISPATCH_EVENT_INTERVAL)
				}
			}
		}
	}
	srvmode := global.SrvMode
	switch srvmode {
	case "moniter":
		{
			go closureFunc()
		}
	}
}

//monitoring client
func (self *Server) moniterClient() {
	timerCheckActive := time.NewTimer(time.Second * global.SERVER_CHECK_ACTIVE_INTERVAL)
	for {
		select {
		case <-timerCheckActive.C:
			{
				timerCheckActive.Reset(time.Second * global.SERVER_CHECK_ACTIVE_INTERVAL)
				//update host status
				global.GHostInfoTable.UpdateHostStatus()
			}
		}
	}
}

//run router server
func (self *Server) Run(host, port string) {
	self.Host = host
	self.Port = port
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
