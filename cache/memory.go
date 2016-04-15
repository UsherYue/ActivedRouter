package cache

//基于内存实现的缓存
type MemoryCache struct {
	cacheMap *MapChan //数据缓存
}

func (self *MemoryCache) GetMemory() *MapChan {
	return self.cacheMap
}

//create memory cache
//type "file" or "memory"
func Newcache(cacheType string) Cache {
	switch cacheType {
	case "memory":
		mapChan := NewDataMapChaner()
		cache := &MemoryCache{cacheMap: mapChan}
		return cache
	case "file":
		{

		}
	}
	return nil
}

//set
func (self *MemoryCache) Set(k string, v interface{}) {
	self.cacheMap.PushKVPair(k, v)
}

//get
func (self *MemoryCache) Get(k string) (interface{}, bool) {
	mapData := *self.cacheMap.GetData()
	val, ok := mapData[k]
	return val, ok
}

//erase
func (self *MemoryCache) Del(k string) {
	self.cacheMap.EraseKVPair(k)
}

//has
func (self *MemoryCache) Has(k string) bool {
	mapData := *self.cacheMap.GetData()
	_, ok := mapData[k]
	return ok
}
