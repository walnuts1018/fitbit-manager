package config

import (
	"fmt"
)

type PSQLDSN struct {
	Dsn          string `env:"DSN" envDefault:""`
	PSQLHost     string `env:"ENDPOINT" envDefault:"localhost"`
	PSQLPort     string `env:"PORT" envDefault:"5432"`
	PSQLDatabase string `env:"DATABASE" envDefault:"fitbit_manager"`
	PSQLUser     string `env:"USER" envDefault:"postgres"`
	PSQLPassword string `env:"PASSWORD,required"`
	PSQLSSLMode  string `env:"SSL_MODE" envDefault:"disable"`
	PSQLTimeZone string `env:"TIMEZONE" envDefault:"Asia/Tokyo"`
}

func (p PSQLDSN) String() string {
	if p.Dsn != "" {
		return p.Dsn
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s", p.PSQLHost, p.PSQLPort, p.PSQLUser, p.PSQLPassword, p.PSQLDatabase, p.PSQLSSLMode, p.PSQLTimeZone)
}
