package config

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
}
