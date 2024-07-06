package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"homework/internal/models"
)

type KafkaSender struct {
	producer *Producer
	topic    string
}

func NewKafkaSender(producer *Producer, topic string) *KafkaSender {
	return &KafkaSender{
		producer,
		topic,
	}
}

func (s *KafkaSender) SendMessage(event models.Event) error {
	kafkaMsg, err := s.buildMessage(event)
	if err != nil {
		fmt.Println("Send event marshal error", err)
		return err
	}

	_, _, err = s.producer.SendSyncMessage(kafkaMsg)

	if err != nil {
		fmt.Println("Send order connector error", err)
		return err
	}

	return nil
}

func (s *KafkaSender) buildMessage(event models.Event) (*sarama.ProducerMessage, error) {
	msg, err := json.Marshal(event)

	if err != nil {
		fmt.Println("Send order marshal error", err)
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic:     s.topic,
		Value:     sarama.ByteEncoder(msg),
		Partition: -1,
		Key:       sarama.StringEncoder(fmt.Sprint(event.ID)),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("test-header"),
				Value: []byte("test-value"),
			},
		},
	}, nil
}