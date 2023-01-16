package agent

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	AddressScheme  string `env:"ADDRESS_SCHEME" envDefault:"http"`
	ReportInterval string `env:"REPORT_INTERVAL"`
	PollInterval   string `env:"POLL_INTERVAL"`
}

func NewConfig() *Config {
	var cfg Config

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address of server with port if needed")
	flag.StringVar(&cfg.ReportInterval, "r", "10s", "report interval")
	flag.StringVar(&cfg.PollInterval, "p", "2s", "collect metrics interval")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
