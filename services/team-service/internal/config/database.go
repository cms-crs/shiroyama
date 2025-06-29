package config

import "time"

type DatabaseConfig struct {
	Host               string        `yaml:"host"`
	User               string        `yaml:"user"`
	Password           string        `yaml:"password"`
	Name               string        `yaml:"name"`
	MaxOpenConnections int           `json:"max_open_connections"`
	MaxIdleConnections int           `json:"max_idle_connections"`
	MaxLifetime        time.Duration `json:"max_lifetime"`
}
