package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var (
	ClientID      *string
	ClientSecret  *string
	UserID        *string
	CookieSecret  *string
	RedisEndpoint *string
	RedisPassword *string
)

func LoadConfig() error {
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

	userID, ok := os.LookupEnv("USER_ID")
	if !ok {
		return fmt.Errorf("USER_ID is not set")
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

	ClientID = &cid
	ClientSecret = &cs
	UserID = &userID
	CookieSecret = &cookieSecret
	RedisEndpoint = &redisEndpoint
	RedisPassword = &redisPassword
	return nil
}
