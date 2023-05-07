package cache

import (
	"container/list"
	"fmt"
	"time"
)

type LRUCache struct {
	cacheInfo
	positionList *list.List
	items        map[string]*list.Element
}

func (c *LRUCache) Set(key string, value interface{}, timeToLive ...time.Duration) (err error) {
	c.Lock()
	defer c.Unlock()

	var item *cacheItem
	if listElement, ok := c.items[key]; ok {
		c.positionList.MoveToFront(listElement)
		item = listElement.Value.(*cacheItem)
		item.value = value
	} else {
		if c.capacity != 0 && c.size == c.capacity {
			c.removeBackElement()
		}

		item = &cacheItem{value: value}
		c.items[key] = c.positionList.PushFront(item)
		c.size++
	}

	if len(timeToLive) == 1 && timeToLive[0] != 0 {
		item.expirationTime = time.Now().Add(timeToLive[0])
	} else if c.defaultTtl != 0 {
		item.expirationTime = time.Now().Add(time.Duration(c.defaultTtl))
	} else {
		item.expirationTime = time.Time{}
	}

	return nil
}

func (c *LRUCache) removeBackElement() {
	if listElement := c.positionList.Back(); listElement != nil {
		c.positionList.Remove(listElement)
	}
}

func (c *LRUCache) Get(key string) (value interface{}, err error) {
	c.RLock()
	defer c.RUnlock()

	if listElement, ok := c.items[key]; ok {
		c.positionList.MoveToFront(listElement)
		item := listElement.Value.(*cacheItem)

		if !item.expirationTime.IsZero() && time.Now().After(item.expirationTime) {
			c.Lock()
			c.removeCacheItem(listElement, key)
			c.Unlock()

			return nil, fmt.Errorf("item does not exist")
		}

		return item.value, nil
	} else {
		return nil, fmt.Errorf("item does not exist")
	}
}

func (c *LRUCache) Delete(key string) (err error) {
	c.Lock()
	defer c.Unlock()

	if listElement, ok := c.items[key]; ok {
		c.removeCacheItem(listElement, key)
		return nil
	} else {
		return fmt.Errorf("item does not exist")
	}
}

func (c *LRUCache) Purge() (err error) {
	c.Lock()
	defer c.Unlock()

	c.positionList = &list.List{}
	c.items = make(map[string]*list.Element)
	c.size = 0
	return nil
}

func (c *LRUCache) removeCacheItem(listElement *list.Element, key string) {
	c.positionList.Remove(listElement)
	delete(c.items, key)
	c.size--
}

func (c *LRUCache) DeleteExpired(timeInterval time.Duration) {
	ticker := time.NewTicker(timeInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.Lock()
		for key, listElement := range c.items {
			item := listElement.Value.(*cacheItem)
			if !item.expirationTime.IsZero() && time.Now().After(item.expirationTime) {
				c.removeCacheItem(listElement, key)
			}
		}
		c.Unlock()
	}
}

func (c *LRUCache) Info() (info map[string]interface{}, err error) {
	info = make(map[string]interface{})
	info["size"] = c.size
	info["capacity"] = c.capacity
	info["defaultTtl"] = c.defaultTtl
	return info, nil
}
