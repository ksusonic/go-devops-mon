package agent

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address        string `json:"address" env:"ADDRESS"`
	AddressScheme  string `json:"address_scheme" env:"ADDRESS_SCHEME" envDefault:"http"`
	ReportInterval string `json:"report_interval" env:"REPORT_INTERVAL"`
	PollInterval   string `json:"poll_interval" env:"POLL_INTERVAL"`
	SecretKey      string `json:"key" env:"KEY"`
	RateLimit      int    `json:"rate_limit" env:"RATE_LIMIT"`
	CryptoKeyPath  string `json:"crypto_key" env:"CRYPTO_KEY"`
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
	flag.StringVar(&cfg.ReportInterval, "r", "10s", "report interval")
	flag.StringVar(&cfg.PollInterval, "p", "2s", "collect metrics interval")
	flag.StringVar(&cfg.SecretKey, "k", "", "key for metric hash")
	flag.IntVar(&cfg.RateLimit, "l", 1, "rate limit for sending metrics")
	flag.StringVar(&cfg.CryptoKeyPath, "crypto-key", "", "public key for https requests")

	flag.Parse()

	err := env.Parse(&cfg)
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
