package netservice

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
)

var DefaultTCPProxy = &TCPProxy{}

//tcp node
type TCPNode struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

//tcp proxy config
type TCPProxyItem struct {
	ListenAddr string     `json:"addr"`
	Method     string     `json:"method"`
	LbNode     []*TCPNode `json:"lbnode"`
}
type TCPProxyConfig struct {
	TCPSwitch     string          `json:"tcp_switch"`
	TCPProxyitems []*TCPProxyItem `json:"tcp_proxy_items"`
}

type TCPProxy struct {
	TCPProxyConfigData TCPProxyConfig
	TCPLBMap           map[string]*TCPProxyItem
	//connection pool
	TCPConnPool     map[string]net.Conn
	TCPConnMap      map[string][]net.Conn
	ProxyConfigFile string
	Mutex           *sync.RWMutex
}

func NewTCPProxy() *TCPProxy {
	return &TCPProxy{Mutex: &sync.RWMutex{}}
}

func (self *TCPProxy) SaveToFile() bool {
	if bts, err := json.Marshal(self.TCPProxyConfigData); err != nil {
		return false
	} else {
		if file, err := os.OpenFile(self.ProxyConfigFile, os.O_RDWR|os.O_TRUNC, os.ModePerm); err != nil {
			defer file.Close()
			return false
		} else {
			if _, err := file.Write(bts); err != nil {
				return false
			}
		}
	}
	return true
}

//Load tcp proxy config file..
func (self *TCPProxy) LoadTCPProxyConfig(configFile string) {
	if f, err := os.Open(configFile); err != nil {
		goto ERROR
	} else {
		if bts, err := ioutil.ReadAll(f); err != nil {
			goto ERROR
		} else {
			if err := json.Unmarshal(bts, &self.TCPProxyConfigData); err != nil {
				goto ERROR
			} else {
				for _, v := range self.TCPProxyConfigData.TCPProxyitems {
					self.TCPLBMap[v.ListenAddr] = v
				}
				return
			}
		}
	}
ERROR:
	log.Fatalln("Parse TCP Config File error!")
}

//lock
func (self *TCPProxy) ProxySend(localAddr, remoteAddr string, data []byte) (int, error) {
	var conn net.Conn
	if c, ok := self.TCPConnPool[remoteAddr]; !ok {
		conn = c
	} else {
		conns := self.TCPConnMap[localAddr]
		if len(conns) == 0 {
			return 0, errors.New("readerror")
		}
		conn = conns[rand.Intn(len(conns))]
	}
	if n, err := conn.Write(data); err != nil {
		return 0, errors.New("send error")
	} else {
		return n, nil
	}
}

//lock
func (self *TCPProxy) ProxyRecv(localAddr, remoteAddr string, data []byte) (int, error) {
	var conn net.Conn
	if c, ok := self.TCPConnPool[remoteAddr]; !ok {
		conn = c
	} else {
		conns := self.TCPConnMap[localAddr]
		if len(conns) == 0 {
			return 0, errors.New("readerror")
		}
		conn = conns[rand.Intn(len(conns))]
	}
	if n, err := conn.Read(data); err != nil {
		return 0, errors.New("readerror")
	} else {
		return n, nil
	}
}

func (self *TCPProxy) StartTCPProxy() {
	for _, v := range self.TCPProxyConfigData.TCPProxyitems {
		if l, err := net.Listen("tcp", v.ListenAddr); err != nil {
			goto ERROR
		} else {
			go func(l net.Listener, addr string) {
				for {
					if conn, err := l.Accept(); err != nil {
						remoteAddr := conn.RemoteAddr().String()
						loadlAddr := addr
						//GO READ
						go func(conn net.Conn, localAddr, remoteAddr string) {

							data := make([]byte, 1024)
							if rdLen, err := conn.Read(data); err != nil {
								self.ProxySend(loadlAddr, remoteAddr, data[0:rdLen])
							}

						}(conn, loadlAddr, remoteAddr)
						//GO WRITE
						go func(conn net.Conn, localAddr, remoteAddr string) {
							for {
								data := make([]byte, 1024)
								self.ProxySend(loadlAddr, remoteAddr, data)
							}
						}(conn, loadlAddr, remoteAddr)
					}
				}
			}(l, v.ListenAddr)
			return
		}
	}
ERROR:
	log.Fatalln("Start TCP Reverse Proxy Server  error!")
}
