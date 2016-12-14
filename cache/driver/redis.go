package driver

//Mysql实现的内存驱动
type RedisContainer struct {
}

//创建数据channer
func NewRedisContainer() *RedisContainer {
	return nil
}

func (this *RedisContainer) PushKVPair(k, v interface{}) Containerer {
	return nil
}

func (this *RedisContainer) EraseKVPair(k interface{}) Containerer {
	return nil
}
func (this *RedisContainer) PushKVMaps(maps ...map[string]interface{}) Containerer {
	return nil
}
func (this *RedisContainer) ResetKVPair(k string, v interface{}) Containerer {
	return nil
}
func (this *RedisContainer) ResetOrAddKVPair(k string, v interface{}) Containerer {
	return nil
}

func (this *RedisContainer) ResetKVPairs(kvMaps map[string]interface{}) Containerer {
	return nil
}

func (this *RedisContainer) ResetOrAddKVPairs(kvMaps map[string]interface{}) Containerer {
	return nil
}

func (this *RedisContainer) GetData() *map[string]interface{} {
	return nil
}
