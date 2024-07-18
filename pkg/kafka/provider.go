package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"homework/internal/models"
	"homework/internal/models/config"
	"sync"
)

type KafkaProvider interface {
	StartConsumer(ctx context.Context, wg *sync.WaitGroup) func() error
	ProduceEvent(event models.Event) error
	GetEvents() chan []byte
}

type kafkaProvider struct {
	consumer      *ConsumerProvider
	producer      *ProducerProvider
	consumerGroup sarama.ConsumerGroup
	zapLogger     *zap.SugaredLogger
}

func NewKafkaProvider(cfg *config.KafkaConfig, zap *zap.SugaredLogger) KafkaProvider {
	group, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, cfg.KafkaGroupID, NewConsumerConfig())
	if err != nil {
		zap.Fatalf("Error creating consumer group: %v", err)
	}
	return &kafkaProvider{
		consumer:      NewConsumerProvider(cfg.KafkaBrokers, cfg.KafkaTopics, zap),
		producer:      NewProducerProvider(cfg.KafkaBrokers, cfg.KafkaTopics, NewProducerConfig),
		consumerGroup: group,
		zapLogger:     zap,
	}
}

func (k *kafkaProvider) StartConsumer(ctx context.Context, wg *sync.WaitGroup) func() error {
	go k.consumer.ConsumeEvents(ctx, wg, k.consumerGroup)

	<-k.consumer.Ready

	return func() error {
		k.zapLogger.Debugln("Closing consumer")
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