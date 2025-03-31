package etc

import (
	"fmt"

	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	// Create a Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort), // Redis address
		Password: "",                                                 // No password
		DB:       0,                                                  // Default DB
	})

	return rdb
}
