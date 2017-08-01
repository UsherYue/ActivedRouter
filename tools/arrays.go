package tools

import (
	"errors"
	"fmt"
	"reflect"
)

//DeleteSlice selete slice
func DeleteSlice(slice interface{}, index int) (interface{}, error) {
	sliceValue := reflect.ValueOf(slice)
	length := sliceValue.Len()
	if slice == nil || length == 0 || (length-1) < index {
		return nil, errors.New("error")
	}
	if length-1 == index {
		return sliceValue.Slice(0, index).Interface(), nil
	} else if (length - 1) >= index {
		return reflect.AppendSlice(sliceValue.Slice(0, index), sliceValue.Slice(index+1, length)).Interface(), nil
	}
	return nil, errors.New("error")
}

//BytesToHexString Convert byte array to hex string .
func BytesToHexString(data []byte) []string {
	ret := []string{}
	for _, byteData := range data {
		ret = append(ret, fmt.Sprintf("%X", byteData))
	}
	return ret
}

//Convert string  to hex string
func StringToHexString(str string) []string {
	return BytesToHexString([]byte(str))
}
