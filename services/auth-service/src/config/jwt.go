package config

import "time"

type JWTConfig struct {
	AccessTokenTTL  time.Duration `yaml:"accessTTL"`
	RefreshTokenTTL int           `yaml:"refreshTTL"`
}
