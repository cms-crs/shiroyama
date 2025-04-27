package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env  string         `yaml:"env" env-default:"local"`
	Grpc GRPCConfig     `yaml:"grpc"`
	DB   DatabaseConfig `yaml:"db"`
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found: " + path)
	}

	var config Config

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic("error reading config file: " + err.Error())
	}

	return &config
}

func fetchConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "config", "", "config file path")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("TASK_SERVICE_CONFIG_PATH")
	}

	return configPath
}
