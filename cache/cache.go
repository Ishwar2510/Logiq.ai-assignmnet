package cache

import (
	"errors"
	"fmt"
	"time"
)

type Cache struct{
	maxSize int
	cache map[string]cacheObject
}
type cacheObject struct{
	value interface{}
	expiration time.Time
}

type CacheMethods interface{
	Add(key string, value interface{}, expiration time.Time) error
	Get(key string) (interface{},error)
	Delete(key string) error
}

func NewCache(maxSize int) *Cache{
	return &Cache{
		maxSize: maxSize,
		cache: make(map[string]cacheObject),
	}
}

func (c *Cache) Add(key string, value interface{}, expiration time.Duration) error {
	if (c.maxSize == len(c.cache)){
			return fmt.Errorf("Max Size reached, please increase the capacity")
	}
		_,err:= c.get(key)
		if err==nil{
			return errors.New("Already exists")
		}
	expTime := time.Now().Add(expiration)

	c.cache[key] = cacheObject{value,expTime};
	return nil
}
func (c* Cache) get(key string) (interface{},error){
	obj, isPresent := c.cache[key]
		if !isPresent{
			return nil,errors.New("not found")
		}
		if obj.expiration.Before(time.Now()){
			c.Delete( key)
			return  nil,errors.New("not found")
		}
	return obj.value,nil
}


func (c *Cache)Get(key string)  (interface{},error){
		return c.get(key)
}

func (c *Cache) Delete(key string) error{
		_,err:=c.get(key)
		if err!=nil{
			return err
		}
		delete(c.cache,key)
	return nil
}




