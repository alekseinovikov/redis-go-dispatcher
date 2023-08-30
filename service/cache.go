package service

import (
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto"
)

type CacheService struct {
	cache        *ristretto.Cache
	service      RedisService
	cacheTtl     time.Duration
	cacheKeysKey string
}

func NewCacheService(
	readService RedisService,
	cacheRefreshDuration time.Duration,
	cacheTtl time.Duration,
) *CacheService {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}

	c := &CacheService{
		cache:        cache,
		service:      readService,
		cacheTtl:     cacheTtl,
		cacheKeysKey: "REDIS_GO_DISPATCHER_CACHE_KEYS",
	}

	ticker := time.NewTicker(cacheRefreshDuration)
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

	for _, key := range keys {
		data, err := c.service.GetByKey(key)
		if err != nil {
			fmt.Println(err)
			continue
		}

		c.cache.SetWithTTL(key, data, 0, c.cacheTtl)
	}

	c.cache.SetWithTTL(c.cacheKeysKey, keys, 0, c.cacheTtl)
}

func (c *CacheService) GetById(id string) (string, error) {
	key := c.service.GetPrefix() + id
	result, found := c.cache.Get(key)
	if !found {
		return "", nil
	}

	return result.(string), nil
}

func (c *CacheService) GetAll() ([]string, error) {
	keys, found := c.cache.Get(c.cacheKeysKey)
	if !found {
		// rollback if we have no keys in cache
		return c.service.GetAll()
	}

	keysSlice := keys.([]string)
	result := make([]string, 0, len(keysSlice))
	for _, key := range keysSlice {
		value, found := c.cache.Get(key)
		if !found {
			continue
		}

		result = append(result, value.(string))
	}

	return result, nil
}
