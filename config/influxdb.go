package config

import "net/url"

type InfluxDBConfig struct {
	Endpoint  url.URL `env:"ENDPOINT" envDefault:"https://localhost:8086"`
	AuthToken string  `env:"AUTH_TOKEN,required"`
	Org       string  `env:"ORG" envDefault:"admin"`
	Bucket    string  `env:"BUCKET" envDefault:"fitbit_manager"`
}
