package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	HTTPPort    int         `yaml:"http-port"`
	MetricsPort int         `yaml:"metrics-port"`
	LogLevel    string      `yaml:"env" env-default:"local"`
	Kafka       KafkaConfig `yaml:"kafka"`
	Redis       RedisConfig `yaml:"redis"`
	Saga        SagaConfig  `yaml:"saga"`
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
		configPath = os.Getenv("TASK_SERVICE_CONFIG_PATH")
	}

	return configPath
}
