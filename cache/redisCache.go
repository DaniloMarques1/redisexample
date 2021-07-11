package cache

import (
	"encoding/json"

	"github.com/danilomarques1/redisexample/entity"
	"github.com/go-redis/redis"
)

type RedisCache struct {
	client *redis.Client
	keyName string // should not be mocked, but for now, it will do
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
		keyName: "configuration-cache",
	}
}

func (rc *RedisCache) IsCacheEmpty() bool {
	value, err := rc.client.Get(rc.keyName).Result()
	if err != nil {
		return true
	}

	return value == ""
}

func (rc *RedisCache) SaveCache(config entity.Config) error {
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = rc.client.Set(rc.keyName, bytes, 0).Err()

	return err
}

func (rc *RedisCache) RemoveCache() error {
	_, err := rc.client.Del(rc.keyName).Result()
	return err
}

func (rc *RedisCache) GetCache() (entity.Config, error) {
	value, err := rc.client.Get(rc.keyName).Result()
	if err != nil {
		return entity.Config{}, err
	}
	var config entity.Config
	err = json.Unmarshal([]byte(value), &config)

	if err != nil {
		return entity.Config{}, err
	}

	return config, nil
}
