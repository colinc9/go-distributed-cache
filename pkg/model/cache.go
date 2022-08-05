package model

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

const defaultTimeDuration time.Duration = 30 * 60

type LRUCache struct {
	lock      sync.RWMutex
	capacity  int
	evictList *list.List
	entries   map[interface{}]*list.Element
}

func NewLRUCache(capacity int) (*LRUCache, error) {
	if capacity <= 0 {
		return nil, errors.New("Capacity must be positive")
	}
	c := &LRUCache{
		capacity:  capacity,
		evictList: list.New(),
		entries:   make(map[interface{}]*list.Element),
	}
	return c, nil
}

func (c *LRUCache) Get(key interface{}) (value interface{}, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if ent, ok := c.entries[key]; ok {
		if ent.Value.(*Entry).IsExpired() {
			return nil, false
		}
		c.evictList.MoveToFront(ent)
		return ent.Value.(*Entry).Value, true
	}
	return nil, false
}

func (c *LRUCache) Set(key, value interface{}) (evictedValue interface{}, evicted bool, ok bool) {
	return c.SetEx(key, value, defaultTimeDuration)
}

func (c *LRUCache) SetEx(key, value interface{}, expire time.Duration) (evictedValue interface{}, evicted bool, ok bool) {
	if expire < 0 {
		return nil, false, false
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	var exp time.Time
	if expire != 0 {
		exp = time.Now().Add(expire)
	} else {
		exp = time.Time{}
	}

	if entry, ok := c.entries[key]; ok {
		c.evictList.MoveToFront(entry)
		entry.Value.(*Entry).Value = value
		entry.Value.(*Entry).Expire = &exp
	}

	newEntry := &Entry{Key: key, Value: value, Expire: &exp}
	entry := c.evictList.PushFront(newEntry)
	c.entries[key] = entry

	if c.evictList.Len() > c.capacity {
		evictedValue, evicted := c.removeLRU()
		return evictedValue, evicted, true
	}
	return nil, false, true
}

func (c *LRUCache) removeLRU() (evictedValue interface{}, evicted bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	entry := c.evictList.Back()
	if entry != nil {
		c.evictList.Remove(entry)
		kv := entry.Value.(*Entry)
		delete(c.entries, kv.Key)
		return kv.Value, true
	}
	return nil, false
}
