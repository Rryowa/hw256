package kafka

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"log"
	"sync"
	"time"
)

type ConsumerInterface interface {
	ConsumeEvents(ctx context.Context, wg *sync.WaitGroup, group sarama.ConsumerGroup)
}

func (consumer *ConsumerProvider) ConsumeEvents(ctx context.Context, wg *sync.WaitGroup, group sarama.ConsumerGroup) {
	wg.Add(1)
	defer wg.Done()

	err := group.Consume(ctx, consumer.Topics, consumer)
	if err != nil {
		if errors.Is(err, sarama.ErrClosedConsumerGroup) {
			log.Println("ConsumerProvider group closed")
			return
		}
		log.Panicf("Error from consumer: %v", err)
	}
	if ctx.Err() != nil {
		return
	}
}

func (consumer *ConsumerProvider) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			consumer.Events <- message.Value
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			log.Println("ConsumerProvider stopped")
			return nil
		}
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *ConsumerProvider) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	consumer.Ready <- true
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *ConsumerProvider) Cleanup(sarama.ConsumerGroupSession) error {
	close(consumer.Events)
	return nil
}

type ConsumerProvider struct {
	Ready   chan bool
	GroupId string
	Brokers []string
	Topics  []string
	Events  chan []byte
}

func NewConsumerProvider(brokers, topics []string) *ConsumerProvider {
	return &ConsumerProvider{
		Ready:   make(chan bool),
		GroupId: "FOO",
		Brokers: brokers,
		Topics:  topics,
		Events:  make(chan []byte),
	}
}

func NewConsumerConfig() *sarama.Config {
	consumerConfig := sarama.NewConfig()
	consumerConfig.Version = sarama.DefaultVersion
	consumerConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	consumerConfig.Consumer.MaxWaitTime = 3 * time.Second
	return consumerConfig
}