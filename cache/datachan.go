package cache

//Map Chan  Tools
type MapChan struct {
	data map[string]interface{}
}

//创建数据channer
func NewDataMapChaner() *MapChan {
	return &MapChan{data: make(map[string]interface{})}
}

func (this *MapChan) PushKVPair(k, v interface{}) *MapChan {
	if key, ok := k.(string); !ok {
		panic("key必须是string类型!")
	} else {
		this.data[key] = v
	}
	return this
}

func (this *MapChan) EraseKVPair(k interface{}) *MapChan {
	if key, ok := k.(string); !ok {
		panic("key必须是string类型!")
	} else {
		delete(this.data, key)
	}
	return this
}
func (this *MapChan) PushKVMaps(maps ...map[string]interface{}) *MapChan {
	for _, itemMap := range maps {
		for itemKey, itemValue := range itemMap {
			this.PushKVPair(itemKey, itemValue)
		}
	}
	return this
}
func (this *MapChan) ResetKVPair(k string, v interface{}) *MapChan {
	if _, ok := this.data[k]; ok {
		this.data[k] = v
	}
	return this
}
func (this *MapChan) ResetOrAddKVPair(k string, v interface{}) *MapChan {
	this.data[k] = v
	return this
}

func (this *MapChan) ResetKVPairs(kvMaps map[string]interface{}) *MapChan {
	for k, v := range kvMaps {
		if _, ok := this.data[k]; ok {
			this.data[k] = v
		}
	}
	return this
}

func (this *MapChan) ResetOrAddKVPairs(kvMaps map[string]interface{}) *MapChan {
	for k, v := range kvMaps {
		this.data[k] = v
	}
	return this
}

func (this *MapChan) GetData() *map[string]interface{} {
	return &this.data
}
