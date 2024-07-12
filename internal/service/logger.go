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
	useKafka     bool
	kafkaService kafka.KafkaService
	repo         storage.Storage
}

// TODO: rename kafka config into logger config
func NewLoggerService(cfg *config.KafkaConfig, repo storage.Storage) LoggerService {
	return &loggerService{
		useKafka:     cfg.KafkaUse,
		kafkaService: kafka.NewKafkaService(cfg),
		repo:         repo,
	}
}

func (l *loggerService) Start(ctx context.Context, wg *sync.WaitGroup) func() error {
	if l.useKafka {
		closer := l.kafkaService.StartConsumer(ctx, wg)
		go l.DisplayKafkaEvents()
		return closer
	}
	return nil
}

func (l *loggerService) CreateEvent(ctx context.Context, input string) (models.Event, error) {
	event, err := l.repo.InsertEvent(ctx, input)
	if err != nil {
		return models.Event{}, err
	}
	if l.useKafka {
		if err := l.kafkaService.ProduceEvent(event); err != nil {
			return models.Event{}, err
		}
	}
	return event, nil
}

func (l *loggerService) ProcessEvent(ctx context.Context, event models.Event) error {
	event, err := l.repo.UpdateEvent(ctx, event)
	if err != nil {
		return err
	}

	if l.useKafka {
		if err := l.kafkaService.ProduceEvent(event); err != nil {
			return err
		}
	} else {
		fmt.Printf("Received event: %v\n", event)
	}
	return nil
}

func (l *loggerService) DisplayKafkaEvents() {
	for event := range l.kafkaService.GetEvents() {
		var v models.Event
		json.Unmarshal(event, &v)
		log.Printf("\nEvent received: %v\n", v)
	}
}