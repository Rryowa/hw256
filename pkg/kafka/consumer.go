package kafka

import (
	"github.com/pkg/errors"
	"time"

	"github.com/IBM/sarama"
)

const Interval = 1 * time.Second

type Consumer struct {
	brokers        []string
	SingleConsumer sarama.Consumer
	initialOffset  int64
}

func NewConsumer(brokers []string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = Interval

	initialOffset := sarama.OffsetOldest
	config.Consumer.Offsets.Initial = initialOffset

	consumer, err := sarama.NewConsumer(brokers, config)

	if err != nil {
		return nil, err
	}

	return &Consumer{
		brokers:        brokers,
		SingleConsumer: consumer,
		initialOffset:  initialOffset,
	}, err
}

func (k *Consumer) Close() error {
	err := k.SingleConsumer.Close()
	if err != nil {
		return errors.Wrap(err, "kafka.Consumer.Close")
	}

	return nil
}