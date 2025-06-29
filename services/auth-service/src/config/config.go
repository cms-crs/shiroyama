package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	Grpc     GRPCConfig     `yaml:"grpc"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
	Kafka    KafkaConfig    `yaml:"kafka"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	fmt.Println(path)
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
		configPath = os.Getenv("AUTH_SERVICE_CONFIG_PATH")
	}

	return configPath
}
