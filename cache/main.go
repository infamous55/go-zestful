package cache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type cacheInfo struct {
	size       uint64
	capacity   uint64
	defaultTtl TimeToLive
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

func (ep *EvictionPolicy) Set(value string) error {
	switch value {
	case "LRU", "LFU":
		*ep = EvictionPolicy(value)
		return nil
	default:
		return fmt.Errorf("parse error")
	}
}

func (ep *EvictionPolicy) String() string {
	return string(*ep)
}

type TimeToLive time.Duration

func (ttl *TimeToLive) Set(value string) error {
	result, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf("parse error")
	}
	*ttl = TimeToLive(result)
	return nil
}

func (ttl *TimeToLive) String() string {
	return fmt.Sprint(*ttl)
}

func New(capacity uint64, evictionPolicy EvictionPolicy, defaultTtl TimeToLive) (cache Cache, err error) {
	switch {
	case evictionPolicy == LRU:
		return &LRUCache{
			cacheInfo: cacheInfo{
				size:       0,
				capacity:   capacity,
				defaultTtl: defaultTtl,
			},
			positionList: &list.List{},
			items:        make(map[string]*list.Element),
		}, nil
	case evictionPolicy == LFU:
		return &LFUCache{
			cacheInfo: cacheInfo{
				size:       0,
				capacity:   capacity,
				defaultTtl: defaultTtl,
			},
			frequencyList: &list.List{},
			items:         make(map[string]*LFUCacheItem),
		}, nil
	default:
		return nil, fmt.Errorf("invalid value \"%v\" for eviction policy", evictionPolicy)
	}
}
