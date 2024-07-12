package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"homework/internal/models"
	"log"
	"sync"
)

type ProducerInterface interface {
	ProduceEvent(event models.Event) error
}

func (p *ProducerProvider) ProduceEvent(event models.Event) error {
	producerTx := p.borrow()
	defer p.release(producerTx)

	err := producerTx.BeginTxn()
	if err != nil {
		return err
	}

	payload, _ := json.Marshal(event)
	for _, topic := range p.topics {
		producerTx.Input() <- &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(payload),
		}
	}

	err = producerTx.CommitTxn()
	if err != nil {
		log.Printf("Producer: unable to commit txn %s\n", err)
		for {
			if producerTx.TxnStatus()&sarama.ProducerTxnFlagFatalError != 0 {
				// fatal error. need to recreate producer.
				log.Printf("Producer: producer is in a fatal state, need to recreate it")
				break
			}
			// If producer is in abortable state, try to abort current transaction.
			if producerTx.TxnStatus()&sarama.ProducerTxnFlagAbortableError != 0 {
				err = producerTx.AbortTxn()
				if err != nil {
					// If an error occured just retry it.
					log.Printf("Producer: unable to abort transaction: %+v", err)
					continue
				}
				break
			}
			// if not you can retry
			err = producerTx.CommitTxn()
			if err != nil {
				log.Printf("Producer: unable to commit txn %s\n", err)
				continue
			}
		}
		return err
	}
	return nil
}

func (p *ProducerProvider) borrow() (producer sarama.AsyncProducer) {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	if len(p.producers) == 0 {
		for {
			producer = p.ProducerProvider()
			if producer != nil {
				return
			}
		}
	}

	index := len(p.producers) - 1
	producer = p.producers[index]
	p.producers = p.producers[:index]
	return
}

func (p *ProducerProvider) release(producer sarama.AsyncProducer) {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	// If released producer is erroneous close it and don't return it to the producer pool.
	if producer.TxnStatus()&sarama.ProducerTxnFlagInError != 0 {
		// Try to close it
		_ = producer.Close()
		return
	}
	p.producers = append(p.producers, producer)
}

func (p *ProducerProvider) Clear() {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	for _, producer := range p.producers {
		producer.Close()
	}
	p.producers = p.producers[:0]
}

// pool of producers that ensure transactional-id is unique.
type ProducerProvider struct {
	brokers                []string
	topics                 []string
	transactionIdGenerator int32

	producersLock sync.Mutex
	producers     []sarama.AsyncProducer

	ProducerProvider func() sarama.AsyncProducer
}

func NewProducerProvider(brokers, topics []string, producerConfigurationProvider func() *sarama.Config) *ProducerProvider {
	provider := &ProducerProvider{
		brokers: brokers,
		topics:  topics,
	}
	provider.ProducerProvider = func() sarama.AsyncProducer {
		config := producerConfigurationProvider()
		suffix := provider.transactionIdGenerator
		// Append transactionIdGenerator to current config.Producer.Transaction.ID to ensure transaction-id uniqueness.
		if config.Producer.Transaction.ID != "" {
			provider.transactionIdGenerator++
			config.Producer.Transaction.ID = config.Producer.Transaction.ID + "-" + fmt.Sprint(suffix)
		}
		producer, err := sarama.NewAsyncProducer(provider.brokers, config)
		if err != nil {
			return nil
		}
		return producer
	}
	return provider
}

func NewProducerConfig() *sarama.Config {
	producerConfig := sarama.NewConfig()
	producerConfig.Version = sarama.DefaultVersion
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	producerConfig.Producer.Transaction.Retry.Backoff = 10
	producerConfig.Producer.Idempotent = true
	producerConfig.Producer.Transaction.ID = "txn_producer"
	producerConfig.Net.MaxOpenRequests = 1
	return producerConfig
}