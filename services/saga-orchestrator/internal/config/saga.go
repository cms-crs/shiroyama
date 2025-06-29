package config

import "time"

type SagaConfig struct {
	Timeout         time.Duration `yaml:"timeout"`
	RetryInterval   time.Duration `yaml:"retryInterval"`
	MaxRetries      int           `yaml:"maxRetries"`
	CleanupInterval time.Duration `yaml:"cleanupInterval"`
}
