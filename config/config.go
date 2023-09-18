package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var (
	ClientID     string
	ClientSecret string
	CookieSecret string
	PSQLEndpoint string
	PSQLPort     string
	PSQLDatabase string
	PSQLUser     string
	PSQLPassword string
	ServerPort   string
)

func LoadConfig() error {
	serverport := flag.String("port", "8080", "server port")
	flag.Parse()

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

	psqlEndpoint, ok := os.LookupEnv("PSQL_ENDPOINT")
	if !ok {
		return fmt.Errorf("PSQL_ENDPOINT is not set")
	}
	psqlPort, ok := os.LookupEnv("PSQL_PORT")
	if !ok {
		return fmt.Errorf("PSQL_PORT is not set")
	}

	psqlDatabase, ok := os.LookupEnv("PSQL_DATABASE")
	if !ok {
		return fmt.Errorf("PSQL_DATABASE is not set")
	}

	psqlUser, ok := os.LookupEnv("PSQL_USER")
	if !ok {
		return fmt.Errorf("PSQL_USER is not set")
	}

	psqlPassword, ok := os.LookupEnv("PSQL_PASSWORD")
	if !ok {
		return fmt.Errorf("PSQL_PASSWORD is not set")
	}

	ClientID = cid
	ClientSecret = cs
	CookieSecret = cookieSecret
	PSQLEndpoint = psqlEndpoint
	PSQLPort = psqlPort
	PSQLDatabase = psqlDatabase
	PSQLUser = psqlUser
	PSQLPassword = psqlPassword
	ServerPort = *serverport
	return nil
}
