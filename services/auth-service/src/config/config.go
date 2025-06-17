package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	Grpc     GRPCConfig     `yaml:"grpc"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
}

func MustLoad() *Config {
	config, err := Load()

	if err != nil {
		panic(err)
	}

	return config
}

func Load() (*Config, error) {
	path, err := fetchConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	var config Config

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func fetchConfigPath() (string, error) {
	var configPath string

	flag.StringVar(&configPath, "config", "", "config file path")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath, nil
}
