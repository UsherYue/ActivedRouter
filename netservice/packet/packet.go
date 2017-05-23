//ActivedRouter
//Author:usher.yue
//Amail:usher.yue@gmail.com
//TencentQQ:4223665
// Provide http/https, tcp and reverse proxy  services

package packet

import (
	"bytes"
	"encoding/binary"
)

const (
	DEFAULE_HEADER           = "[**********]"
	DEFAULT_HEADER_LENGTH    = 12
	DEFAULT_SAVE_DATA_LENGTH = 4
)

type Packet struct {
	Header         string
	HeaderLengh    int32
	SaveDataLength int32
	Data           []byte
}

//set delimiter header
func (self *Packet) SetHeader(header string) *Packet {
	self.Header = header
	self.HeaderLengh = int32(len([]byte(header)))
	return self
}

//create default package
func NewDefaultPacket(data []byte) *Packet {
	return &Packet{DEFAULE_HEADER, DEFAULT_HEADER_LENGTH, DEFAULT_SAVE_DATA_LENGTH, data}
}

//convert to net package
func (self *Packet) Packet() []byte {
	return append(append([]byte(self.Header), self.IntToBytes(int32(len(self.Data)))...), self.Data...)
}

//return value is sticky data
func (self *Packet) UnPacket(readerChannel chan []byte) []byte {
	dataLen := int32(len(self.Data))
	var i int32
	for i = 0; i < dataLen; i++ {
		//Termiate for loop when the remaining data is insufficient .
		if dataLen < i+self.HeaderLengh+self.SaveDataLength {
			break
		}
		//find Header
		if string(self.Data[i:i+self.HeaderLengh]) == self.Header {
			saveDataLenBeginIndex := i + self.HeaderLengh
			actualDataLen := self.BytesToInt(self.Data[saveDataLenBeginIndex : saveDataLenBeginIndex+self.SaveDataLength])
			//The remaining data is less than one package
			if dataLen < i+self.HeaderLengh+self.SaveDataLength+actualDataLen {
				break
			}
			//Get a packet
			packageData := self.Data[saveDataLenBeginIndex+self.SaveDataLength : saveDataLenBeginIndex+self.SaveDataLength+actualDataLen]
			//send pacakge data to reader channel
			readerChannel <- packageData
			//get next package index
			i += self.HeaderLengh + self.SaveDataLength + actualDataLen - 1
		}
	}
	//Reach the end
	if i >= dataLen {
		return []byte{}
	}
	//Returns the remaining data
	return self.Data[i:]
}

func (self *Packet) IntToBytes(i int32) []byte {
	byteBuffer := bytes.NewBuffer([]byte{})
	binary.Write(byteBuffer, binary.BigEndian, i)
	return byteBuffer.Bytes()
}

func (self *Packet) BytesToInt(data []byte) int32 {
	var val int32
	byteBuffer := bytes.NewBuffer(data)
	binary.Read(byteBuffer, binary.BigEndian, &val)
	return val
}
