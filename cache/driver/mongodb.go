package driver

//Mysql实现的内存驱动
type MongoContainer struct {
}

//创建数据channer
func NewMongoContainer() *MongoContainer {
	return nil
}

func (this *MongoContainer) PushKVPair(k, v interface{}) Containerer {
	return nil
}

func (this *MongoContainer) EraseKVPair(k interface{}) Containerer {
	return nil
}
func (this *MongoContainer) PushKVMaps(maps ...map[string]interface{}) Containerer {
	return nil
}
func (this *MongoContainer) ResetKVPair(k string, v interface{}) Containerer {
	return nil
}
func (this *MongoContainer) ResetOrAddKVPair(k string, v interface{}) Containerer {
	return nil
}

func (this *MongoContainer) ResetKVPairs(kvMaps map[string]interface{}) Containerer {
	return nil
}

func (this *MongoContainer) ResetOrAddKVPairs(kvMaps map[string]interface{}) Containerer {
	return nil
}

func (this *MongoContainer) GetData() *map[string]interface{} {
	return nil
}
