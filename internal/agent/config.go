package agent

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	AddressScheme  string `env:"ADDRESS_SCHEME" envDefault:"http"`
	ReportInterval string `env:"REPORT_INTERVAL"`
	PollInterval   string `env:"POLL_INTERVAL"`
	SecretKey      string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address of server with port if needed")
	flag.StringVar(&cfg.ReportInterval, "r", "10s", "report interval")
	flag.StringVar(&cfg.PollInterval, "p", "2s", "collect metrics interval")
	flag.StringVar(&cfg.SecretKey, "k", "", "key for metric hash")
	flag.IntVar(&cfg.RateLimit, "l", 1, "rate limit for sending metrics")

	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
