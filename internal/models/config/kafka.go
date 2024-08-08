package config

type KafkaConfig struct {
	KafkaUse     bool     `env:"KAFKA_USE"`
	KafkaBrokers []string `env:"KAFKA_BROKERS"`
	KafkaTopics  []string `env:"KAFKA_TOPICS"`
	KafkaGroupID string   `env:"KAFKA_GROUP_ID"`
}