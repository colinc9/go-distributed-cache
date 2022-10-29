package service

import (
	"github.com/colinc9/go-distributed-cache/pkg/model"
	"github.com/colinc9/go-distributed-cache/pkg/service/tcp"
)

type CacheService struct {
	Cache *model.LRUCache
}

func (c *CacheService) Get(key interface{}) (value interface{}, ok bool) {
	msg := tcp.Message{
		Type:  tcp.Get,
		Key:   key,
		Value: nil,
	}
	tcp.DialTcp(&msg)
	return c.Cache.Get(key)
}


func (c *CacheService) Set(key, value interface{}) (evictedValue interface{}, evicted bool, ok bool) {
	msg := tcp.Message{
		Type:  tcp.Set,
		Key:   key,
		Value: value,
	}
	tcp.DialTcp(&msg)
	return c.Cache.Set(key, value)
}
