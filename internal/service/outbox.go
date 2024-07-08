package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"homework/internal/models"
	"homework/internal/storage"
	"homework/pkg/kafka"
	"log"
	"time"
)

type Outbox struct {
	repo     storage.Storage
	UseKafka bool
	Brokers  []string
	Topic    string
}

func NewOutbox(repo storage.Storage, cfg *models.Config) *Outbox {
	return &Outbox{
		repo:     repo,
		Brokers:  cfg.Brokers,
		UseKafka: cfg.UseKafka,
		Topic:    cfg.Topic,
	}
}

func (o *Outbox) StartProcessingEvents(ctx context.Context, done chan struct{}) {
	ticker := time.NewTicker(kafka.Interval)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				events, err := o.repo.GetEvents(ctx)
				if err != nil {
					log.Println("Error getting message from repo:", err)
					continue
				}
				if len(events) == 0 {
					continue
				}

				events, err = o.repo.ProcessEvents(ctx, events)
				if err != nil {
					log.Println("Error processing event:", err)
					continue
				}

				if o.UseKafka {
					err := o.ReceiveKafka()
					if err != nil {
						log.Println("Error receiving message from Kafka:", err)
						continue
					}
					for _, event := range events {
						err = o.SendKafka(event)
						if err != nil {
							log.Println("Error sending message to Kafka:", err)
							continue
						}
					}
				} else {
					o.OutputToConsole(events)
				}
			}
		}
	}()
}

func (o *Outbox) SendKafka(event models.Event) error {
	kafkaProducer, err := kafka.NewProducer(o.Brokers)
	if err != nil {
		log.Fatal(err)
	}
	producer := kafka.NewKafkaSender(kafkaProducer, o.Topic)
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

func (o *Outbox) ReceiveKafka() error {
	kafkaConsumer, err := kafka.NewConsumer(o.Brokers)
	if err != nil {
		fmt.Println(err)
	}

	handlers := map[string]kafka.HandleFunc{
		o.Topic: func(message *sarama.ConsumerMessage) {
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
	err = recv.Subscribe(o.Topic)
	if err != nil {
		return err
	}
	return nil
}

func (o *Outbox) OutputToConsole(events []models.Event) {
	for _, event := range events {
		fmt.Printf("\nOutput to Console: Event ID: %d, Request: %s, Status: %s, AcquiredAt: %v, ProcessedAt: %v\n",
			event.ID, event.Request, event.Status, event.AcquiredAt, event.ProcessedAt)
	}
}