package repository

import (
	"database/sql"
	"fmt"

	"github.com/danilomarques1/redisexample/entity"
)

type ConfigRepositorySql struct {
	db *sql.DB
}

func NewConfigRepositorySql(db *sql.DB) *ConfigRepositorySql {
	return &ConfigRepositorySql{
		db: db,
	}
}

func (cr *ConfigRepositorySql) Create(config entity.Config) (entity.Config, error) {
	configDb, _ := cr.Get()

	fmt.Println("config from db = ", configDb)

	if (entity.Config{}) == configDb {
		return cr.createFlow(config)
	} else {
		return cr.updateFlow(configDb, config)
	}
}

func (cr *ConfigRepositorySql) updateFlow(oldConfig, config entity.Config) (entity.Config, error) {
	if cr.shouldUpdate(oldConfig, config) {
		return cr.updateConfig(config)
	}
	config.Id = oldConfig.Id
	return config, nil
}

func (cr *ConfigRepositorySql) createFlow(config entity.Config) (entity.Config, error) {
	return cr.addConfig(config)
}

func (cr *ConfigRepositorySql) addConfig(config entity.Config) (entity.Config, error) {
	stmt, err := cr.db.Prepare("INSERT INTO config(timeout, label_name) VALUES($1, $2)")
	if err != nil {
		return entity.Config{}, err
	}
	defer stmt.Close()
	fmt.Println("config to be added", config)

	_, err = stmt.Exec(config.Timeout, config.LabelName)
	if err != nil {
		return entity.Config{}, err
	}

	return config, nil
}

func (cr *ConfigRepositorySql) updateConfig(config entity.Config) (entity.Config, error) {
	stmt, err := cr.db.Prepare("update config set timeout = $1, label_name = $2")
	if err != nil {
		return entity.Config{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(config.Timeout, config.LabelName)
	if err != nil {
		return entity.Config{}, err
	}

	return config, nil
}

func (cr *ConfigRepositorySql) Get() (entity.Config, error) {
	stmt, err := cr.db.Prepare("select * from config")
	if err != nil {
		return entity.Config{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRow()
	if row.Err() != nil {
		return entity.Config{}, row.Err()
	}
	var config entity.Config
	row.Scan(&config.Id, &config.Timeout, &config.LabelName)

	return config, nil
}

func (cr *ConfigRepositorySql) shouldUpdate(oldConfig, newConfig entity.Config) bool {
	if oldConfig.Timeout != newConfig.Timeout || oldConfig.LabelName != newConfig.LabelName {
		return true
	}

	return false
}
