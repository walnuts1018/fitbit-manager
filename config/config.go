package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var (
	ClientID      string
	ClientSecret  string
	CookieSecret  string
	RedisEndpoint string
	RedisPassword string
	Port          string
)

func LoadConfig() error {
	port := flag.String("port", "8080", "server port")
	flag.Parse()
	slog.Info(fmt.Sprintf("port: %v", *port))

	err := godotenv.Load(".env")
	if err != nil {
		slog.Warn("Error loading .env file")
	}

	cid, ok := os.LookupEnv("CLIENT_ID")
	if !ok {
		return fmt.Errorf("CLIENT_ID is not set")
	}

	cs, ok := os.LookupEnv("CLIENT_SECRET")
	if !ok {
		return fmt.Errorf("CLIENT_SECRET is not set")
	}

	cookieSecret, ok := os.LookupEnv("COOKIE_SECRET")
	if !ok {
		return fmt.Errorf("COOKIE_SECRET is not set")
	}

	redisEndpoint, ok := os.LookupEnv("REDIS_ENDPOINT")
	if !ok {
		return fmt.Errorf("REDIS_ENDPOINT is not set")
	}

	redisPassword, ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		return fmt.Errorf("REDIS_PASSWORD is not set")
	}

	ClientID = cid
	ClientSecret = cs
	CookieSecret = cookieSecret
	RedisEndpoint = redisEndpoint
	RedisPassword = redisPassword
	Port = *port
	return nil
}
