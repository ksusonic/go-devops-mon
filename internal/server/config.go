package server

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address          string `env:"ADDRESS"`
	FileDropInterval string `env:"STORE_INTERVAL"`
	StoreFile        string `env:"STORE_FILE"`
	RestoreFile      bool   `env:"RESTORE"`
}

func NewConfig() *Config {
	var cfg Config

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address of server with port if needed")
	flag.BoolVar(&cfg.RestoreFile, "r", true, "bool if restore from file needed")
	flag.StringVar(&cfg.FileDropInterval, "i", "300s", "interval to save metrics to file")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "path to file for metrics save")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
