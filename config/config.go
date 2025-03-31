package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

// Config ...
type Config struct {
	Environment string // develop, staging, production

	UserServiceHost string
	UserServicePort int

	TweetServiceHost string
	TweetServicePort int

	LikeServiceHost string
	LikeServicePort int

	FanoutServiceHost string
	FanoutServicePort int

	// ApiGateWayHost string
	// ApiGateWayPort int

	LogLevel string
	HttpPort string

	RedisHost string
	RedisPort string

	RabbitMqHost     string
	RabbitMqPort     int
	RabbitMqUser     string
	RabbitMqPassword string

	AccessTokenSecret  string
	RefreshTokenSecret string
}

// Load loads environment vars and inflates Config
func Load() Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}

	config := Config{}
	config.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "staging"))
	config.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	config.HttpPort = cast.ToString(getOrReturnDefault("HTTP_PORT", "8000"))

	config.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVICE_HOST", "localhost"))
	config.UserServicePort = cast.ToInt(getOrReturnDefault("USER_SERVICE_PORT", 8082))

	config.TweetServiceHost = cast.ToString(getOrReturnDefault("TWEET_SERVICE_HOST", "localhost"))
	config.TweetServicePort = cast.ToInt(getOrReturnDefault("TWEET_SERVICE_PORT", 8083))

	config.LikeServiceHost = cast.ToString(getOrReturnDefault("LIKE_SERVICE_HOST", "localhost"))
	config.LikeServicePort = cast.ToInt(getOrReturnDefault("LIKE_SERVICE_PORT", 8084))

	config.FanoutServiceHost = cast.ToString(getOrReturnDefault("FANOUT_SERVICE_HOST", "localhost"))
	config.FanoutServicePort = cast.ToInt(getOrReturnDefault("FANOUT_SERVICE_PORT", 8085))

	// config.ApiGateWayHost = cast.ToString(getOrReturnDefault("API_GATEWAY_HOST", "localhost"))
	// config.ApiGateWayPort = cast.ToInt(getOrReturnDefault("API_GATEWAY_PORT", 8080))

	config.RedisHost = cast.ToString(getOrReturnDefault("REDIS_HOST", "localhost"))
	config.RedisPort = cast.ToString(getOrReturnDefault("REDIS_PORT", "6379"))

	// rabbitmq
	config.RabbitMqHost = cast.ToString(getOrReturnDefault("RABBITMQ_HOST", "localhost"))
	config.RabbitMqPort = cast.ToInt(getOrReturnDefault("RABBITMQ_PORT", 5672))
	config.RabbitMqUser = cast.ToString(getOrReturnDefault("RABBITMQ_USER", "guest"))
	config.RabbitMqPassword = cast.ToString(getOrReturnDefault("RABBITMQ_PASSWORD", "guest"))

	config.AccessTokenSecret = cast.ToString(getOrReturnDefault("ACCESS_TOKEN_SECRET", "secret"))
	config.RefreshTokenSecret = cast.ToString(getOrReturnDefault("REFRESH_TOKEN_SECRET", "secret"))

	return config
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
