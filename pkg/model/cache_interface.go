package model

import (
	"time"
)

type Cache interface {
	Get(key interface{}) (interface{}, bool)
	Set(key, value interface{}) bool
	SetEx(key, value interface{}, expire time.Duration) bool
}

type Entry struct {
	Key    interface{}
	Value  interface{}
	Expire *time.Time
}

func (e *Entry) IsExpired() bool {
	if e.Expire == nil || e.Expire.IsZero() {
		return false
	}
	return time.Now().After(*e.Expire)
}
