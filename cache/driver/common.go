//ActivedRouter
//Author:usher.yue
//Amail:usher.yue@gmail.com
//TencentQQ:4223665
//缓存驱动接口定义
//可以自定义扩展驱动类型 默认是 memory cache

package driver

//容器接口声明
type Containerer interface {
	PushKVPair(k, v interface{}) Containerer
	EraseKVPair(k interface{}) Containerer
	PushKVMaps(maps ...map[string]interface{}) Containerer
	ResetKVPair(k string, v interface{}) Containerer
	ResetOrAddKVPair(k string, v interface{}) Containerer
	ResetKVPairs(kvMaps map[string]interface{}) Containerer
	ResetOrAddKVPairs(kvMaps map[string]interface{}) Containerer
	Exist(k interface{}) bool
	GetData() *map[string]interface{}
}

//基于内存实现的缓存
type CacheImpl struct {
	Driver Containerer //数据缓存驱动
}

func (self *CacheImpl) Exist(k interface{}) bool {

	return self.Driver.Exist(k)
}

func (self *CacheImpl) GetStorage() Containerer {
	return self.Driver
}

//set
func (self *CacheImpl) Set(k string, v interface{}) {
	self.Driver.PushKVPair(k, v)
}

//get
func (self *CacheImpl) Get(k string) (interface{}, bool) {
	mapData := *self.Driver.GetData()
	val, ok := mapData[k]
	return val, ok
}

//erase
func (self *CacheImpl) Del(k string) {
	self.Driver.EraseKVPair(k)
}

//has
func (self *CacheImpl) Has(k string) bool {
	mapData := *self.Driver.GetData()
	_, ok := mapData[k]
	return ok
}
