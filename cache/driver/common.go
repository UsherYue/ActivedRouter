//ActivedRouter
//Author:usher.yue
//Amail:usher.yue@gmail.com
//TencentQQ:4223665
//cache driver interface define

package driver

//Container interface define
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

//Cache based on memory implementation
type CacheImpl struct {
	Driver Containerer
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
