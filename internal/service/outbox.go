package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"homework/broker/kafka"
	"homework/internal/models"
	"log"
	"strings"
	"time"
)

const topic = "orders"

type Outbox interface {
	CreateEvent(ctx context.Context, request string) error
	GetEvent(ctx context.Context) (models.Event, error)
	ProcessEvent(ctx context.Context, e models.Event) (models.Event, error)
	StartProcessingEvents(ctx context.Context, done chan struct{})
}

type OutboxRepo struct {
	Brokers  []string
	Pool     *pgxpool.Pool
	UseKafka bool
}

func NewOutbox(ctx context.Context, cfg *models.Config) Outbox {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal(err, "db connection error")
	}
	log.Println("Connected to db")

	return &OutboxRepo{
		Pool:     pool,
		Brokers:  cfg.Brokers,
		UseKafka: cfg.UseKafka,
	}
}

func (o *OutboxRepo) StartProcessingEvents(ctx context.Context, done chan struct{}) {
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				event, err := o.GetEvent(ctx)
				if err != nil {
					log.Println("Error getting message from repo:", err)
					continue
				}
				if event.ID == 0 {
					continue
				}

				event, err = o.ProcessEvent(ctx, event)
				if err != nil {
					log.Println("Error processing event:", err)
					continue
				}

				if o.UseKafka {
					err = o.SendKafka(event)
					if err != nil {
						log.Println("Error sending message to Kafka:", err)
						continue
					}

					err = o.ReceiveKafka()
					if err != nil {
						log.Println("Error receiving message from Kafka:", err)
						continue
					}
				} else {
					o.OutputToConsole(event)
				}
			}
		}
	}()
}

func (o *OutboxRepo) CreateEvent(ctx context.Context, request string) error {
	const op = "CreateEvent"
	args := strings.Split(request, " ")

	_, err := o.Pool.Exec(ctx,
		`INSERT INTO events (method_name, request, acquired, acquired_at) 
				VALUES ($1, $2, true, NOW())`,
		args[0], request)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (o *OutboxRepo) GetEvent(ctx context.Context) (models.Event, error) {
	const op = "GetEvent"
	row := o.Pool.QueryRow(ctx,
		`SELECT id, request, acquired, processed, acquired_at
				FROM events WHERE acquired = true AND processed = false`)
	var e models.Event
	err := row.Scan(&e.ID, &e.Request, &e.Acquired, &e.Processed, &e.AcquiredAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return e, nil
		}
		return e, fmt.Errorf("%s: %w", op, err)
	}

	return e, nil
}

func (o *OutboxRepo) ProcessEvent(ctx context.Context, e models.Event) (models.Event, error) {
	const op = "ProcessEvent"
	e.Processed = true
	row := o.Pool.QueryRow(ctx,
		`UPDATE events SET processed = $1, processed_at = NOW()
				WHERE id = $2 returning processed_at`,
		e.Processed, e.ID)
	err := row.Scan(&e.ProcessedAt)
	if err != nil {
		return models.Event{}, fmt.Errorf("%s: %w", op, err)
	}

	return e, nil
}

func (o *OutboxRepo) SendKafka(event models.Event) error {
	kafkaProducer, err := kafka.NewProducer(o.Brokers)
	if err != nil {
		log.Fatal(err)
	}
	producer := kafka.NewKafkaSender(kafkaProducer, topic)
	err = producer.SendMessage(event)
	if err != nil {
		return err
	}
	err = kafkaProducer.Close()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (o *OutboxRepo) ReceiveKafka() error {
	kafkaConsumer, err := kafka.NewConsumer(o.Brokers)
	if err != nil {
		fmt.Println(err)
	}

	handlers := map[string]kafka.HandleFunc{
		topic: func(message *sarama.ConsumerMessage) {
			event := models.Event{}
			err = json.Unmarshal(message.Value, &event)
			if err != nil {
				fmt.Println("Consumer error", err)
			}
			fmt.Println("Topic:", message.Topic, ", Partition:", message.Partition, ", Offset:", message.Offset)
			fmt.Println("Key:", string(message.Key), ", Received Value:", string(message.Value))
		},
	}

	recv := kafka.NewReceiver(kafkaConsumer, handlers)
	err = recv.Subscribe(topic)
	if err != nil {
		return err
	}
	return nil
}

func (o *OutboxRepo) OutputToConsole(event models.Event) {
	fmt.Printf("\nOutput to Console: Event ID: %d, Request: %s, Acquired: %t, Processed: %t, AcquiredAt: %v, ProcessedAt: %v\n",
		event.ID, event.Request, event.Acquired, event.Processed, event.AcquiredAt, event.ProcessedAt)
}