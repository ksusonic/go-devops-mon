package agent

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address        string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	AddressScheme  string `env:"ADDRESS_SCHEME" envDefault:"http"`
	ReportInterval string `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   string `env:"POLL_INTERVAL" envDefault:"2s"`
}

func NewConfig() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &cfg
}
