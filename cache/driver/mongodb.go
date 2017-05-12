package driver

type MongoContainer struct {
}

func NewMongoContainer() *MongoContainer {
	return nil
}
func (this *MongoContainer) Exist(k interface{}) bool {
	return true
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
