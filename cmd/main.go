package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/api"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/etc"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/logger"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/services"
	"github.com/casbin/casbin/v2"
	rediscache "github.com/golanguzb70/redis-cache"
)

func main() {
	cfg := config.Load()
	fmt.Printf("Users service: %s:%d\n", cfg.UserServiceHost, cfg.UserServicePort)
	fmt.Printf("Tweet service: %s:%d\n", cfg.TweetServiceHost, cfg.TweetServicePort)
	fmt.Printf("Like service: %s:%d\n", cfg.LikeServiceHost, cfg.LikeServicePort)
	fmt.Printf("Fanout service: %s:%d\n", cfg.FanoutServiceHost, cfg.FanoutServicePort)

	log := logger.New(cfg.LogLevel, "Mini_Tweeter_api_gateway")
	gprcClients, err := services.NewGrpcClients(&cfg)
	if err != nil {
		log.Error("Error while getting grpc clients", logger.Error(err))
		return
	}

	cache := rediscache.New(&rediscache.Config{
		RedisHost: cfg.RedisHost,
		RedisPort: cfg.RedisPort,
	})

	modelPath := getConfigPath("model.conf")
	policyPath := getConfigPath("policy.csv")

	casbinEnforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		fmt.Println("Error loading model and policy:", err)
		return
	}

	server := api.New(&api.RouterOptions{
		Log:      log,
		Cfg:      &cfg,
		Services: gprcClients,
		Cache:    cache,
		Enforcer: casbinEnforcer,
		Redis:    etc.NewRedisClient(&cfg),
	})

	server.Run(":" + cfg.HttpPort)
}

func getConfigPath(filename string) string {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		return filename // fallback to relative path
	}

	return filepath.Join(wd, "config", filename)
}
