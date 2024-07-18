package kafka

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"homework/internal/models"
	"homework/internal/models/config"
	"homework/internal/storage"
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
	kafkaProvider KafkaProvider
	repo          storage.Storage
	zapLogger     *zap.SugaredLogger
}

func NewLoggerService(cfg *config.KafkaConfig, repo storage.Storage, zap *zap.SugaredLogger) LoggerService {
	return &loggerService{
		useKafka:      cfg.KafkaUse,
		kafkaProvider: NewKafkaProvider(cfg, zap),
		repo:          repo,
		zapLogger:     zap,
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
		logger.displayEvent(event)
	}
	return nil
}

func (logger *loggerService) DisplayKafkaEvents() {
	for event := range logger.kafkaProvider.GetEvents() {
		var v models.Event
		json.Unmarshal(event, &v)
		logger.zapLogger.Infof("Kafka: Received event: %v\n", v)
	}
}

func (logger *loggerService) displayEvent(event models.Event) {
	logger.zapLogger.Infof("IO: Received event: %v\n", event)
}