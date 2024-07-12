package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"homework/internal/models"
	"homework/internal/models/config"
	"log"
	"sync"
)

type KafkaService interface {
	StartConsumer(ctx context.Context, wg *sync.WaitGroup) func() error
	ProduceEvent(event models.Event) error
	GetEvents() chan []byte
}

type kafkaProvider struct {
	consumer      *ConsumerProvider
	producer      *ProducerProvider
	consumerGroup sarama.ConsumerGroup
}

func NewKafkaService(cfg *config.KafkaConfig) KafkaService {
	group, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, cfg.KafkaGroupID, NewConsumerConfig())
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	return &kafkaProvider{
		consumer:      NewConsumerProvider(cfg.KafkaBrokers, cfg.KafkaTopics),
		producer:      NewProducerProvider(cfg.KafkaBrokers, cfg.KafkaTopics, NewProducerConfig),
		consumerGroup: group,
	}
}

func (k *kafkaProvider) StartConsumer(ctx context.Context, wg *sync.WaitGroup) func() error {
	go k.consumer.ConsumeEvents(ctx, wg, k.consumerGroup)

	<-k.consumer.Ready

	return func() error {
		log.Println("Closing consumer")
		return k.consumerGroup.Close()
	}
}

func (k *kafkaProvider) ProduceEvent(event models.Event) error {
	if err := k.producer.ProduceEvent(event); err != nil {
		return err
	}

	return nil
}

func (k *kafkaProvider) GetEvents() chan []byte {
	return k.consumer.Events
}