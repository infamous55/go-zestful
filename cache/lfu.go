package cache

import (
	"container/list"
	"fmt"
	"time"
)

type LFUCache struct {
	cacheInfo
	frequencyList *list.List
	items         map[string]*LFUCacheItem
}

type FrequencyListItem struct {
	value           uint64
	associatedItems map[string]struct{}
}

type LFUCacheItem struct {
	*cacheItem
	frequencyIndicator *list.Element
}

func (c *LFUCache) Set(key string, value interface{}, timeToLive ...time.Duration) (err error) {
	c.Lock()
	defer c.Unlock()

	item, ok := c.items[key]
	if ok {
		item.value = value
	} else {
		if c.capacity != 0 && c.size == c.capacity {
			c.removeBackItems()
		}

		item = &LFUCacheItem{cacheItem: &cacheItem{value: value}}
		c.items[key] = item
		c.size++

		frequencyListBackElement := c.frequencyList.Back()
		frequencyListItem := frequencyListBackElement.Value.(*FrequencyListItem)
		if frequencyListItem == nil {
			frequencyListItem = &FrequencyListItem{value: 0}
		} else if frequencyListItem.value != 0 {
			newFrequencyListItem := &FrequencyListItem{value: 0}
			frequencyListBackElement = c.frequencyList.PushBack(newFrequencyListItem)
			frequencyListItem = newFrequencyListItem
		}

		frequencyListItem.associatedItems[key] = struct{}{}
		item.frequencyIndicator = frequencyListBackElement
	}

	if len(timeToLive) == 1 && timeToLive[0] != 0 {
		item.expirationTime = time.Now().Add(timeToLive[0])
	} else {
		item.expirationTime = time.Time{}
	}

	return nil
}

func (c *LFUCache) removeBackItems() {
	frequencyListBackElement := c.frequencyList.Back()
	if frequencyListBackElement != nil {
		frequencyListItem := frequencyListBackElement.Value.(*FrequencyListItem)
		for item := range frequencyListItem.associatedItems {
			delete(c.items, item)
			c.size--
		}
		c.frequencyList.Remove(frequencyListBackElement)
	}
}

func (c *LFUCache) Get(key string) (value interface{}, err error) {
	c.RLock()
	defer c.RUnlock()

	if item, ok := c.items[key]; ok {
		if !item.expirationTime.IsZero() && time.Now().After(item.expirationTime) {
			c.Lock()
			c.removeCacheItem(item, key)
			c.Unlock()

			return nil, fmt.Errorf("item does not exist")
		}

		c.Lock()
		c.incrementItemFrequency(item, key)
		c.Unlock()

		return item.value, nil
	} else {
		return nil, fmt.Errorf("item does not exist")
	}
}

func (c *LFUCache) removeCacheItem(item *LFUCacheItem, key string) {
	frequencyListElement := item.frequencyIndicator
	frequencyListItem := frequencyListElement.Value.(*FrequencyListItem)

	delete(frequencyListItem.associatedItems, key)
	if len(frequencyListItem.associatedItems) == 0 {
		c.frequencyList.Remove(frequencyListElement)
	}

	delete(c.items, key)
	c.size--
}

func (c *LFUCache) incrementItemFrequency(item *LFUCacheItem, key string) {
	currentFrequencyListElement := item.frequencyIndicator
	currentFrequencyListItem := currentFrequencyListElement.Value.(*FrequencyListItem)
	newFrequencyValue := currentFrequencyListItem.value + 1

	nextFrequencyListElement := currentFrequencyListElement.Next()
	nextFrequencyListItem, ok := nextFrequencyListElement.Value.(*FrequencyListItem)

	if !ok || nextFrequencyListItem.value != newFrequencyValue {
		newFrequencyListItem := &FrequencyListItem{value: newFrequencyValue}
		newFrequencyListItem.associatedItems[key] = struct{}{}
		c.frequencyList.InsertAfter(newFrequencyListItem, currentFrequencyListElement)
	} else {
		nextFrequencyListItem.associatedItems[key] = struct{}{}
	}
}

func (c *LFUCache) Delete(key string) (err error) {
	c.Lock()
	defer c.Unlock()

	if item, ok := c.items[key]; ok {
		c.removeCacheItem(item, key)
		return nil
	} else {
		return fmt.Errorf("item does not exist")
	}
}

func (c *LFUCache) Purge() (err error) {
	c.Lock()
	defer c.Unlock()

	c.frequencyList = &list.List{}
	c.items = make(map[string]*LFUCacheItem)
	c.size = 0
	return nil
}

func (c *LFUCache) DeleteExpired(timeInterval time.Duration) {
	ticker := time.NewTicker(timeInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.Lock()
		for key, item := range c.items {
			if !item.expirationTime.IsZero() && time.Now().After(item.expirationTime) {
				c.removeCacheItem(item, key)
			}
		}
		c.Unlock()
	}
}
