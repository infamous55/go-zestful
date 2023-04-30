package cache

import (
	"container/list"
	"fmt"
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

type EvictionPolicy string

const (
	LRU EvictionPolicy = "LRU"
	LFU EvictionPolicy = "LFU"
)

func New(evictionPolicy EvictionPolicy, capacity uint64) (cache Cache, err error) {
	switch {
	case evictionPolicy == LRU:
		return &LRUCache{
			cacheInfo: cacheInfo{
				size:     0,
				capacity: capacity,
			},
			positionList: &list.List{},
			items:        make(map[string]*list.Element),
		}, nil
	case evictionPolicy == LFU:
		return &LFUCache{
			cacheInfo: cacheInfo{
				size:     0,
				capacity: capacity,
			},
			frequencyList: &list.List{},
			items:         make(map[string]*LFUCacheItem),
		}, nil
	default:
		return nil, fmt.Errorf("invalid eviction policy")
	}
}
