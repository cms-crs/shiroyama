package config

import "time"

type ClientsConfig struct {
	UserServiceAddr string        `yaml:"userServiceAddr" envDefault:"user-service:44044"`
	TeamServiceAddr string        `yaml:"teamServiceAddr" envDefault:"team-service:44045"`
	ClientTimeout   time.Duration `yaml:"clientTimeout" envDefault:"30s"`
}
