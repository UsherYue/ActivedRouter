package packet

import (
	"ActivedRouter/tools"
	"fmt"
	"testing"
	"time"
)

//Packet encapsulation
func getPackage(rawData string) []byte {
	packetTool := NewDefaultPacket([]byte(rawData)).SetHeader("#####")
	pkg := packetTool.Packet()
	//append package
	pkg = append(pkg, pkg...)
	pkg = append(pkg, pkg...)
	pkg = append(pkg, pkg...)
	pkg = append(pkg, pkg...)
	pkg = append(pkg, []byte{0xA, 0xB}...)
	return pkg
}

func Test_Packet(t *testing.T) {
	rawData := "I am a  Programmer!!"
	//package data
	packetData := getPackage(rawData)
	//1024个缓冲
	readerChan := make(chan []byte, 1024)
	//unpackage
	packetTool := NewDefaultPacket(packetData).SetHeader("#####")
	remainData := packetTool.UnPackage(readerChan)
	//remain data
	fmt.Println("Remain Data:", tools.BytesToHexString(remainData))
	go func(reader chan []byte) {
		for {
			data := <-reader
			fmt.Println("Send:", string(data))
		}
	}(readerChan)
	time.Sleep(time.Second * 5)

}
