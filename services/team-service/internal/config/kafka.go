package config

type KafkaConfig struct {
	Brokers []string `yaml:"brokers" env-default:"localhost:29092"`
	GroupID string   `yaml:"group_id" env-default:"team-service-group"`
}
