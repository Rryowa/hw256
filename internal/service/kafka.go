package service

import (
	"context"
	"github.com/IBM/sarama"
	"homework/internal/models"
	"homework/internal/models/config"
	"homework/internal/storage"
	"homework/pkg/kafka"
	"log"
	"sync"
)

type KafkaService interface {
	UseKafka() bool
	StartConsumer(ctx context.Context, wg *sync.WaitGroup) func() error
	ProduceEvent(ctx context.Context, input string) error
	CreateEvent(ctx context.Context, input string) (models.Event, error)
	GetEvents() chan []byte
}

type kafkaProvider struct {
	useKafka      bool
	consumer      *kafka.ConsumerProvider
	producer      *kafka.ProducerProvider
	consumerGroup sarama.ConsumerGroup
	repository    storage.Storage
}

func NewKafkaService(cfg *config.KafkaConfig, repo storage.Storage) KafkaService {
	group, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, cfg.KafkaGroupID, kafka.NewConsumerConfig())
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	return &kafkaProvider{
		useKafka:      cfg.KafkaUse,
		consumer:      kafka.NewConsumerProvider(cfg.KafkaBrokers, cfg.KafkaTopics),
		producer:      kafka.NewProducerProvider(cfg.KafkaBrokers, cfg.KafkaTopics, kafka.NewProducerConfig),
		consumerGroup: group,
		repository:    repo,
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

func (k *kafkaProvider) ProduceEvent(ctx context.Context, input string) error {
	event, err := k.CreateEvent(ctx, input)

	err = k.producer.ProduceEvent(event)
	if err != nil {
		return err
	}

	return nil
}

func (k *kafkaProvider) UseKafka() bool {
	return k.useKafka
}

func (k *kafkaProvider) CreateEvent(ctx context.Context, input string) (models.Event, error) {
	return k.repository.InsertEvent(ctx, input)
}

func (k *kafkaProvider) GetEvents() chan []byte {
	return k.consumer.Events
}