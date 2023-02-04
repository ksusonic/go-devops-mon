package server

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address                  string        `env:"ADDRESS"`
	FileDropInterval         string        `env:"STORE_INTERVAL"`
	FileDropIntervalDuration time.Duration `env:"-"`
	StoreFile                string        `env:"STORE_FILE"`
	RestoreFile              bool          `env:"RESTORE"`
	SecretKey                *string       `env:"KEY"`
	DatabaseURL              string        `env:"DATABASE_DSN"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address of server with port if needed")
	flag.BoolVar(&cfg.RestoreFile, "r", true, "bool if restore from file needed")
	flag.StringVar(&cfg.FileDropInterval, "i", "300s", "interval to save metrics to file")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "path to file for metrics save")
	flag.StringVar(&cfg.DatabaseURL, "d", "", "url for database")
	cfg.SecretKey = flag.String("k", "", "key for metric hash")

	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	cfg.FileDropIntervalDuration, err = time.ParseDuration(cfg.FileDropInterval)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
