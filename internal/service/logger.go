package service

import (
	"context"
	"encoding/json"
	"fmt"
	"homework/internal/models"
	"homework/internal/models/config"
	"homework/internal/storage"
	"homework/pkg/kafka"
	"log"
	"sync"
)

type LoggerService interface {
	Start(ctx context.Context, wg *sync.WaitGroup) func() error
	DisplayKafkaEvents()
	CreateEvent(ctx context.Context, input string) (models.Event, error)
	ProcessEvent(ctx context.Context, event models.Event) error
}

type loggerService struct {
	useKafka      bool
	kafkaProvider kafka.KafkaProvider
	repo          storage.Storage
}

func NewLoggerService(cfg *config.KafkaConfig, repo storage.Storage) LoggerService {
	return &loggerService{
		useKafka:      cfg.KafkaUse,
		kafkaProvider: kafka.NewKafkaProvider(cfg),
		repo:          repo,
	}
}

func (logger *loggerService) Start(ctx context.Context, wg *sync.WaitGroup) func() error {
	if logger.useKafka {
		closer := logger.kafkaProvider.StartConsumer(ctx, wg)
		go logger.DisplayKafkaEvents()
		return closer
	}
	return nil
}

func (logger *loggerService) CreateEvent(ctx context.Context, input string) (models.Event, error) {
	event, err := logger.repo.InsertEvent(ctx, input)
	if err != nil {
		return models.Event{}, err
	}
	if logger.useKafka {
		if err := logger.kafkaProvider.ProduceEvent(event); err != nil {
			return models.Event{}, err
		}
	}
	return event, nil
}

func (logger *loggerService) ProcessEvent(ctx context.Context, event models.Event) error {
	event, err := logger.repo.UpdateEvent(ctx, event)
	if err != nil {
		return err
	}

	if logger.useKafka {
		if err := logger.kafkaProvider.ProduceEvent(event); err != nil {
			return err
		}
	} else {
		fmt.Printf("Received event: %v\n", event)
	}
	return nil
}

func (logger *loggerService) DisplayKafkaEvents() {
	for event := range logger.kafkaProvider.GetEvents() {
		var v models.Event
		json.Unmarshal(event, &v)
		log.Printf("\nEvent received: %v\n", v)
	}
}