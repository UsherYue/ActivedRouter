package driver

type RedisContainer struct {
}

func NewRedisContainer() *RedisContainer {
	return nil
}
func (this *RedisContainer) Exist(k interface{}) bool {
	return true
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
