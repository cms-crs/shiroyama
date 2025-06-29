package config

import "time"

type KafkaConfig struct {
	Brokers  []string       `yaml:"brokers"`
	Producer ProducerConfig `yaml:"producer"`
	Consumer ConsumerConfig `yaml:"consumer"`
}

type ProducerConfig struct {
	RetryMax     int           `yaml:"retryMax"`
	FlushTimeout time.Duration `yaml:"flushTimeout"`
	BatchSize    int           `yaml:"batchSize"`
	Compression  string        `yaml:"compression"`
}

type ConsumerConfig struct {
	GroupID           string        `yaml:"groupId"`
	SessionTimeout    time.Duration `yaml:"sessionTimeout"`
	HeartbeatInterval time.Duration `yaml:"heartbeatInterval"`
	AutoOffsetReset   string        `yaml:"autoOffsetReset"`
	MaxProcessingTime time.Duration `yaml:"maxProcessingTime"`
}
