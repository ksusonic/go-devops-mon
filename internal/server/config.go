package server

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address          string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	FileDropInterval string `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile        string `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	RestoreFile      bool   `env:"RESTORE" envDefault:"true"`
}

func NewConfig() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &cfg
}
