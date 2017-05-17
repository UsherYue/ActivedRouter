package datapackage

import (
	"bytes"
	"encoding/binary"
)

const (
	DEFAULE_HEADER           = "[**********]"
	DEFAULT_HEADER_LENGTH    = 12
	DEFAULT_SAVE_DATA_LENGTH = 4
)

type Package struct {
	Header         string
	HeaderLengh    int32
	SaveDataLength int32
	Data           []byte
}

//create default package
func NewDefaultPackage(data []byte) *Package {
	return &Package{DEFAULE_HEADER, DEFAULT_HEADER_LENGTH, DEFAULT_SAVE_DATA_LENGTH, data}
}

//convert to net package
func (self *Package) Package() []byte {
	return append(append([]byte(self.Header), self.IntToBytes(self.SaveDataLength)...), self.Data...)
}

func (self *Package) UnPackage() []byte {
	return nil
}

func (self *Package) IntToBytes(i int32) []byte {
	byteBuffer := bytes.NewBuffer([]byte{})
	binary.Write(byteBuffer, binary.BigEndian, i)
	return byteBuffer.Bytes()
}

func (self *Package) BytesToInt(data []byte) int32 {
	var val int32
	byteBuffer := bytes.NewBuffer(data)
	binary.Read(byteBuffer, binary.BigEndian, &val)
	return val
}
