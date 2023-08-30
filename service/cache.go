package service

import (
	"fmt"
	"sync"
	"time"
)

type ReadService interface {
	GetByKey(id string) (string, error)
	GetAll() ([]string, error)
	GetAllKeys() ([]string, error)
	GetPrefix() string
}

type CacheService struct {
	cache   map[string]string
	service ReadService
	rwLock  sync.RWMutex
}

func NewCacheService(readService ReadService, duration time.Duration) *CacheService {
	c := &CacheService{
		cache:   make(map[string]string),
		service: readService,
		rwLock:  sync.RWMutex{},
	}

	ticker := time.NewTicker(duration)
	go c.warmUpCacheJob(ticker)

	return c
}

func (c *CacheService) warmUpCacheJob(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			c.warmUpCache()
		}
	}
}

func (c *CacheService) warmUpCache() {
	keys, err := c.service.GetAllKeys()
	if err != nil {
		fmt.Println(err)
		return
	}

	c.rwLock.Lock()
	defer c.rwLock.Unlock()

	c.cache = make(map[string]string)
	for _, key := range keys {
		data, err := c.service.GetByKey(key)
		if err != nil {
			fmt.Println(err)
			continue
		}

		c.cache[key] = data
	}
}

func (c *CacheService) GetById(id string) (string, error) {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()

	key := c.service.GetPrefix() + id
	result, found := c.cache[key]
	if !found {
		return "", nil
	}

	return result, nil
}

func (c *CacheService) GetAll() ([]string, error) {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()

	result := make([]string, 0, len(c.cache))
	for _, value := range c.cache {
		result = append(result, value)
	}

	return result, nil
}
