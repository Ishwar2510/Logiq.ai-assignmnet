package cache

import (
	"errors"
	"time"
	"sync"
	
)

var (
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrMaxLimitReached = errors.New("maximum cache limit reached")
)


type inMemCache struct {
	maxSize int
	cache   map[string]cacheObject
	mu *sync.RWMutex 
}

type cacheObject struct {
	value      interface{}
	expiration time.Time
}

type Cache interface {
	Add(key string, value interface{}, expiration time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	MaxSize(size int)
}

func NewCache(maxSize int) Cache {
	return &inMemCache{
		maxSize: maxSize,
		cache:   make(map[string]cacheObject),
		mu: &sync.RWMutex{},
	}
}

func (c *inMemCache) Add(key string, value interface{}, expiration time.Duration) error {
	

	if c.isFull(){
		c.remove()
	}
	_,expired,err:=c.get(key)
	if err == nil && !expired {
		return ErrAlreadyExists
	}

	 if expired {
		c.delete(key)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheObject{
		value:      value,
		expiration: time.Now().Add(expiration),
	}
	return nil
}

func (c *inMemCache) isFull() bool{
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.maxSize > len(c.cache) {
		return false 
	}
	return true
}

func (c *inMemCache) remove()  {
	
	now := time.Now()
	var oldestKey string
	var oldestTime time.Time
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, obj := range c.cache {
		if obj.expiration.Before(now) {
			delete(c.cache, key)
			return 
		}
		if obj.expiration.Before(oldestTime) {
			oldestTime = obj.expiration
			oldestKey = key
		}
	}
	
		delete(c.cache, oldestKey)
	
	
}

func (c *inMemCache) get(key string) (interface{},bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	obj, isPresent := c.cache[key]
	if !isPresent {
		return nil, false,ErrNotFound
	}
	if obj.expiration.Before(time.Now()) {
		return nil,true, nil
	}
	return obj.value, false,nil
}

func (c *inMemCache) Get(key string) (interface{}, error) {
	 value,expired,err:=c.get(key)
	 if err!=nil{
		return nil,err
	 }

	 if expired {
		c.delete(key)
		return nil,ErrNotFound
	 }
	return value,nil
}

func (c *inMemCache) delete(key string){
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.cache,key)
} 
func (c *inMemCache) MaxSize(size int) {
	c.maxSize = size
}

func (c *inMemCache) Delete(key string) error {
	_,expired,err:=c.get(key)
	 if err!=nil{
		return err
	 }

	 if expired {
		c.delete(key)
		return ErrNotFound
	 }
	 c.delete(key)
	return nil
}
