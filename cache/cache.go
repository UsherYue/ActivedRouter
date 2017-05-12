//ActivedRouter
//Author:usher.yue
//Amail:usher.yue@gmail.com
//TencentQQ:4223665
//cache driver interface define
//we can customize the drive type,default driver is memory

package cache

import (
	"ActivedRouter/cache/driver"
)

//cache  interface declare
type Cacher interface {
	GetStorage() driver.Containerer
	Set(k string, v interface{})
	Get(k string) (interface{}, bool)
	Del(k string)
	Has(k string) bool
}

//create memory cache
//type "file" or "memory"
func Newcache(cacheType string) Cacher {
	switch cacheType {
	case "memory":
		return &driver.CacheImpl{Driver: driver.NewMapContainer()}
	case "file":
		{
			return &driver.CacheImpl{Driver: driver.NewFileContainer()}
		}
	case "mysql":
		{
			return &driver.CacheImpl{Driver: driver.NewMysqlContainer()}
		}
	case "redis":
		{
			return &driver.CacheImpl{Driver: driver.NewRedisContainer()}
		}
	case "mongodb":
		{
			return &driver.CacheImpl{Driver: driver.NewMongoContainer()}
		}
	}
	return nil
}
