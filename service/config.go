package service

import (
	"log"

	"github.com/danilomarques1/redisexample/cache"
	"github.com/danilomarques1/redisexample/dto"
	"github.com/danilomarques1/redisexample/entity"
)

type ConfigService struct {
	configRepository entity.ConfigRepository
	cacheService     cache.Cache
}

func NewConfigService(configRepository entity.ConfigRepository, cacheService cache.Cache) *ConfigService {
	return &ConfigService{
		configRepository: configRepository,
		cacheService:     cacheService,
	}
}

func (cs *ConfigService) AddConfig(configDto dto.ConfigDto) (entity.Config, error) {
	// check if there is something cached on redis
	// if so, removes it
	if !cs.cacheService.IsCacheEmpty() {
		log.Printf("Cache is not empty when trying to add, removing\n")
		err := cs.cacheService.RemoveCache()
		if err != nil {
			log.Fatalf("Error removing cache %v", err)
		}
	}

	configToBeSaved := entity.Config{
		Timeout:   configDto.Timeout,
		LabelName: configDto.LabelName,
	}

	log.Printf("Adding new data\n")
	config, err := cs.configRepository.Create(configToBeSaved)
	if err != nil {
		return entity.Config{}, err
	}

	return config, nil
}

func (cs *ConfigService) GetConfig() (entity.Config, error) {
	// first we should look up redis, if there is an entry
	// we return that entry, if not, we search the database
	// before returning the database value, we should save on redis
	// for future usage
	if !cs.cacheService.IsCacheEmpty() {
		log.Printf("Cache is not empty\n")
		return cs.cacheService.GetCache()
	}

	log.Printf("Getting data from db\n")
	config, err := cs.configRepository.Get()
	if err != nil {
		return entity.Config{}, err
	}

	log.Printf("Saving data to cache")
	err = cs.cacheService.SaveCache(config)
	if err != nil {
		log.Fatalf("Error saving cache %v", err)
	}

	return config, nil
}

/*
func cacheIsEmpty() bool {
	return (entity.Config{}) == configCache
}

func saveCache(config entity.Config) {
	configCache = config
}

func removeCache() {
	configCache = entity.Config{}
}
*/
