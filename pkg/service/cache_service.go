package service

import "github.com/colinc9/go-distributed-cache/pkg/model"

type CacheService struct {
	Cache *model.LRUCache
}

func (c *CacheService) Get(key interface{}) (value interface{}, ok bool) {
	return c.Cache.Get(key)
}


func (c *CacheService) Set(key, value interface{}) (evictedValue interface{}, evicted bool, ok bool) {
	return c.Cache.Set(key, value)
}
