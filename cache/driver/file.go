package driver

//Mysql实现的内存驱动
type FileContainer struct {
}

//创建数据channer
func NewFileContainer() *FileContainer {
	return nil
}

func (this *FileContainer) Exist(k interface{}) bool {
	return true
}
func (this *FileContainer) PushKVPair(k, v interface{}) Containerer {
	return nil
}

func (this *FileContainer) EraseKVPair(k interface{}) Containerer {
	return nil
}
func (this *FileContainer) PushKVMaps(maps ...map[string]interface{}) Containerer {
	return nil
}
func (this *FileContainer) ResetKVPair(k string, v interface{}) Containerer {
	return nil
}
func (this *FileContainer) ResetOrAddKVPair(k string, v interface{}) Containerer {
	return nil
}

func (this *FileContainer) ResetKVPairs(kvMaps map[string]interface{}) Containerer {
	return nil
}

func (this *FileContainer) ResetOrAddKVPairs(kvMaps map[string]interface{}) Containerer {
	return nil
}

func (this *FileContainer) GetData() *map[string]interface{} {
	return nil
}
