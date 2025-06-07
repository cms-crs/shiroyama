package config

import (
	"errors"
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DB   DatabaseConfig `yaml:"db"`
	Grpc GRPCConfig     `yaml:"grpc"`
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

	err := godotenv.Load(".env")
	if err != nil {
		return "", errors.New("error loading .env file")
	}

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath, nil
}
