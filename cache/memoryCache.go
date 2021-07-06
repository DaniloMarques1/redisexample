package cache

import "github.com/danilomarques1/redisexample/entity"


/*
        // TODO next
        // allow to save any kind of data using a map

        map[config] -> configValue
        map[user] -> userValue
*/

// TODO receive interface from the save cache
// TODO return interface from the get cache

type MemoryCache struct {
	cacheValue entity.Config
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		cacheValue: entity.Config{},
	}
}

func (m *MemoryCache) IsCacheEmpty() bool {
	return (entity.Config{}) == m.cacheValue
}

func (m *MemoryCache) SaveCache(value entity.Config) error {
	m.cacheValue = value
	return nil
}

func (m *MemoryCache) RemoveCache() error {
	m.cacheValue = entity.Config{}
	return nil
}

func (m *MemoryCache) GetCache() (entity.Config, error) {
	return m.cacheValue, nil
}
