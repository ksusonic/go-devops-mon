package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address                  string        `json:"address" env:"ADDRESS"`
	FileDropInterval         string        `json:"file_drop_interval" env:"STORE_INTERVAL"`
	FileDropIntervalDuration time.Duration `json:"-" env:"-"`
	StoreFile                string        `json:"store_file" env:"STORE_FILE"`
	RestoreFile              bool          `json:"restore_file" env:"RESTORE"`
	SecretKey                string        `json:"key" env:"KEY"`
	DatabaseURL              string        `json:"database_dsn" env:"DATABASE_DSN"`
	Debug                    bool          `json:"debug" env:"DEBUG"`
	CryptoKeyPath            string        `json:"crypto_key" env:"CRYPTO_KEY"`
	TrustedSubnet            string        `json:"trusted_subnet" env:"TRUSTED_SUBNET"`
}

func NewConfig() (*Config, error) {
	var (
		cfg        Config
		configPath string
	)
	if configPath := findConfigArg(); configPath != nil {
		jsonConfig, err := preloadConfig(*configPath)
		if err != nil {
			return nil, fmt.Errorf("error preloading json-config: %w", err)
		}
		cfg = *jsonConfig
	}

	flag.StringVar(&configPath, "c", "", "json config path")
	flag.StringVar(&configPath, "config", "", "json config path")
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address of server with port if needed")
	flag.BoolVar(&cfg.RestoreFile, "r", true, "bool if restore from file needed")
	flag.StringVar(&cfg.FileDropInterval, "i", "300s", "interval to save metrics to file")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "path to file for metrics save")
	flag.StringVar(&cfg.DatabaseURL, "d", "", "url for database")
	flag.BoolVar(&cfg.Debug, "debug", false, "debug mode")
	flag.StringVar(&cfg.SecretKey, "k", "", "key for metric hash")
	flag.StringVar(&cfg.CryptoKeyPath, "crypto-key", "", "private key for tls")
	flag.StringVar(&cfg.TrustedSubnet, "t", "", "trusted subnets")

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

func preloadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	all, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	var c Config
	err = json.Unmarshal(all, &c)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}
	return &c, nil
}

func findConfigArg() *string {
	path := os.Getenv("CONFIG")
	if path != "" {
		return &path
	}
	flagSet := flag.NewFlagSet("configFlagSet", flag.ContinueOnError)
	cFlag := flagSet.String("c", "", "json config path")
	configFlag := flagSet.String("config", "", "json config path")

	// no variant to disable err logging in flag package
	savedStdErr := os.Stderr
	os.Stderr = nil
	// flag.ContinueOnError returns first error and does not parse further args
	for _, arg := range os.Args[1:] {
		_ = flagSet.Parse([]string{arg})
	}
	os.Stderr = savedStdErr
	if cFlag != nil && *cFlag != "" {
		return cFlag
	}
	if configFlag != nil && *configFlag != "" {
		return configFlag
	}
	return nil
}
