package config

import "time"

type RedisConfig struct {
	Port int           `yaml:"port"`
	Host string        `yaml:"host"`
	TTL  time.Duration `yaml:"ttl"`
}
