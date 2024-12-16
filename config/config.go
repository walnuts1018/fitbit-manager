package config

import (
	"log/slog"
	"reflect"
	"time"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	ServerPort string     `env:"SERVER_PORT" envDefault:"8080"`
	LogLevel   slog.Level `env:"LOG_LEVEL"`
	LogType    LogType    `env:"LOG_TYPE" envDefault:"json"`
	ServerURL  ServerURL  `env:"SERVER_URL" envDefault:"https://fitbit-manager.local.walnuts.dev/"`

	UserID         UserID         `env:"USER_ID,required"`
	ClientID       ClientID       `env:"CLIENT_ID,required"`
	ClientSecret   ClientSecret   `env:"CLIENT_SECRET,required"`
	CookieSecret   CookieSecret   `env:"COOKIE_SECRET,required"`
	PSQLDSN        PSQLDSN        `envPrefix:"PSQL_"`
	InfluxDBConfig InfluxDBConfig `envPrefix:"INFLUXDB_"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.ParseWithOptions(&cfg, env.Options{
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(slog.Level(0)):    returnAny(ParseLogLevel),
			reflect.TypeOf(time.Duration(0)): returnAny(time.ParseDuration),
			reflect.TypeOf(LogType("")):      returnAny(ParseLogType),
			reflect.TypeOf(CookieSecret("")): returnAny(ParseCookieSecret),
			reflect.TypeOf(ServerURL("")):    returnAny(ParseServerURL),
		},
	}); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func returnAny[T any](f func(v string) (t T, err error)) func(v string) (any, error) {
	return func(v string) (any, error) {
		t, err := f(v)
		return any(t), err
	}
}
