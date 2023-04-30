package cache

import (
	"sync"
	"time"
)

type cacheInfo struct {
	size     uint64
	capacity uint64
	sync.RWMutex
}

type cacheItem struct {
	value          interface{}
	expirationTime time.Time
}

type Cache interface {
	Set(key string, value interface{}, timeToLive ...time.Duration) (err error)
	Get(key string) (value interface{}, err error)
	Delete(key string) (err error)
	Purge() (err error)
	DeleteExpired(timeInterval time.Duration)
}
