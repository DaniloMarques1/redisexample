package cache

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/danilomarques1/redisexample/entity"
)

type FileCache struct {
	fileName string
}

func NewFileCache(fileName string) *FileCache {
	return &FileCache{
		fileName: fileName,
	}
}

func (f *FileCache) IsCacheEmpty() bool {
	file, err := os.Open(f.fileName)
	if err != nil {
		return true
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading cache %v\n", err)
		return true
	}

	return len(bytes) == 0
}

func (f *FileCache) GetCache() (entity.Config, error) {
	file, err := os.Open(f.fileName)
	if err != nil {
		return entity.Config{}, err
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading cache %v\n", err)
		return entity.Config{}, err
	}

	log.Printf("%v\n", string(bytes))
	var config entity.Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Printf("Error generating bytes %v\n", err)
		return entity.Config{}, err
	}

	return config, nil
}

func (f *FileCache) RemoveCache() error {
	file, err := os.OpenFile(f.fileName, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = file.WriteString("")
	if err != nil {
		log.Printf("Error removing cache %v\n", err)
		return err
	}

	return nil
}

func (f *FileCache) SaveCache(config entity.Config) error {
	file, err := os.OpenFile(f.fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
