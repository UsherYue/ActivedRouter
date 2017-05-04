package driver

//Mysql实现的内存驱动
type MysqlContainer struct {
}

//创建数据channer
func NewMysqlContainer() *MysqlContainer {
	return nil
}
func (this *MysqlContainer) Exist(k interface{}) bool {
	return true
}
func (this *MysqlContainer) PushKVPair(k, v interface{}) Containerer {
	return nil
}

func (this *MysqlContainer) EraseKVPair(k interface{}) Containerer {
	return nil
}
func (this *MysqlContainer) PushKVMaps(maps ...map[string]interface{}) Containerer {
	return nil
}
func (this *MysqlContainer) ResetKVPair(k string, v interface{}) Containerer {
	return nil
}
func (this *MysqlContainer) ResetOrAddKVPair(k string, v interface{}) Containerer {
	return nil
}

func (this *MysqlContainer) ResetKVPairs(kvMaps map[string]interface{}) Containerer {
	return nil
}

func (this *MysqlContainer) ResetOrAddKVPairs(kvMaps map[string]interface{}) Containerer {
	return nil
}

func (this *MysqlContainer) GetData() *map[string]interface{} {
	return nil
}
