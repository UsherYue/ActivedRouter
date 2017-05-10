package netservice

import (
	"errors"
	"log"
	"reflect"
	"testing"
)

func Test_proxyConfig(t *testing.T) {
	ProxyHandler.LoadProxyConfig("../config/http_proxy.json")
	log.Println(ProxyHandler.Cfg)
	//ProxyHandler.DeleteDomainConig("www.abcd.com")
	///ProxyHandler.AddDomainConfig("www.abcd.com")
	//ProxyHandler.AddProxyClient("www.abcd.com", "127.0.0.21", "8080")
	//	ProxyHandler.DeleteProxyClient("www.abcd.com", "127.0.0.21", "8080")
	//	result, _ := DeleteSlice([]int{1, 2, 3, 4, 5, 6}, 0)
	//	log.Println(result)
	//ProxyHandler.AddDomainConfig("xxxxxxx")
	//log.Println(ProxyHandler.DomainInfos())
}

func DeleteSlice(slice interface{}, index int) (interface{}, error) {
	sliceValue := reflect.ValueOf(slice)
	length := sliceValue.Len()
	if slice == nil || length == 0 || length < index {
		return nil, errors.New("error")
	}
	if length-1 == index {
		return sliceValue.Slice(0, length-1), nil
	} else if (length - 1) > index {
		return reflect.AppendSlice(sliceValue.Slice(0, index), sliceValue.Slice(index+1, length-1)), nil
	}
	return nil, errors.New("error")
}
