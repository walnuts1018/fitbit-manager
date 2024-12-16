package config

type InfluxDBConfig struct {
	Endpoint  string `env:"ENDPOINT" envDefault:"localhost:8086"`
	AuthToken string `env:"AUTH_TOKEN,required"`
	Org       string `env:"ORG" envDefault:"admin"`
	Bucket    string `env:"BUCKET" envDefault:"fitbit_manager"`
}
