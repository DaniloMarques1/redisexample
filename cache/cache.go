package cache

import (
	"github.com/danilomarques1/redisexample/entity"
)

type Cache interface {
	IsCacheEmpty() bool
	SaveCache(entity.Config) error // TODO receive interface
	RemoveCache() error
	GetCache() (entity.Config, error) // TODO return interface
}
