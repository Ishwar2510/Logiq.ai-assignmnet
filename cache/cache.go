package cache

import (
	"errors"
	"time"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrMaxLimitReached = errors.New("maximum cache limit reached")
)

type inMemCache struct {
	maxSize int
	cache   map[string]cacheObject
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
	}
}

func (c *inMemCache) Add(key string, value interface{}, expiration time.Duration) error {
	if c.maxSize == len(c.cache) {
		c.remove()
		
	}
	_, err := c.get(key)
	if err == nil {
		return ErrAlreadyExists
	}
	c.cache[key] = cacheObject{
		value:      value,
		expiration: time.Now().Add(expiration),
	}
	return nil
}
func (c *inMemCache) remove() {
	now := time.Now()
	var oldestKey string
	var oldestTime time.Time
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

func (c *inMemCache) get(key string) (interface{}, error) {
	obj, isPresent := c.cache[key]
	if !isPresent {
		return nil, ErrNotFound
	}
	if obj.expiration.Before(time.Now()) {
		delete(c.cache, key)
		return nil, ErrNotFound
	}
	return obj.value, nil
}

func (c *inMemCache) Get(key string) (interface{}, error) {
	return c.get(key)
}

func (c *inMemCache) MaxSize(size int) {
	c.maxSize = size
}

func (c *inMemCache) Delete(key string) error {
	_, err := c.get(key)
	if err != nil {
		return err
	}
	delete(c.cache, key)
	return nil
}
