//缓存队列可基于内存  文件 数据库来实现
package cache

//cache接口声明
type Cache interface {
	GetMemory() *MapChan
	Set(k string, v interface{})
	Get(k string) (interface{}, bool)
	Del(k string)
	Has(k string) bool
}
