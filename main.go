package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/danilomarques1/redisexample/cache"
	"github.com/danilomarques1/redisexample/dto"
	"github.com/danilomarques1/redisexample/repository"
	"github.com/danilomarques1/redisexample/service"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const db_query = `
        CREATE TABLE IF NOT EXISTS config(
                int serial primary key not null,
                timeout integer not null,
                label_name varchar(60) not null
        );
`

type ConfigHandler struct {
	configService *service.ConfigService
}

func NewConfigHandler(configService *service.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		configService: configService,
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env variables %v\n", err)
	}

	db, err := sql.Open("postgres", fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"),
		os.Getenv("DB_PWD")))

	if err != nil {
		log.Fatalf("Error opening database connection %v", err)
	}
	if _, err = db.Exec(db_query); err != nil {
		log.Fatalf("Error creating table and database %v", err)
	}

	configRepository := repository.NewConfigRepositorySql(db)
	//cache := cache.NewMemoryCache() // in memory cache
	//cache := cache.NewFileCache("redisexample_file_config") // in file cache
	client := redis.NewClient(&redis.Options{ // redis cache
		Addr:     "0.0.0.0:6379",
		Password: "",
		DB:       0,
	})
	cache := cache.NewRedisCache(client)
	configService := service.NewConfigService(configRepository, cache)
	configHandler := NewConfigHandler(configService)

	http.HandleFunc("/config", configHandler.handlerConfig)

	fmt.Println("Application starting at port 8080")
	http.ListenAndServe(":8080", nil)
}

func (ch *ConfigHandler) handlerConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ch.getConfig(w, r)
	} else if r.Method == http.MethodPost {
		ch.addConfig(w, r)
	}
}

func (ch *ConfigHandler) addConfig(w http.ResponseWriter, r *http.Request) {
	var configDto dto.ConfigDto
	json.NewDecoder(r.Body).Decode(&configDto)
	config, err := ch.configService.AddConfig(configDto)
	if err != nil {
		json.NewEncoder(w).Encode(dto.ErrorDto{Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(config)
}

func (ch *ConfigHandler) getConfig(w http.ResponseWriter, r *http.Request) {
	config, err := ch.configService.GetConfig()
	if err != nil {
		json.NewEncoder(w).Encode(dto.ErrorDto{Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(config)
}
