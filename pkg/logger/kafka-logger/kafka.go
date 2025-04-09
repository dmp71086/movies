package kafka_logger

import (
	"log"
	"strings"
	"sync"

	"github.com/IBM/sarama"
)

var (
	mx       sync.Mutex
	err      error
	producer sarama.SyncProducer
)

func init() {
	connect()
}

func connect() {
	brokers := strings.Split("env", ",")
	for i, v := range brokers {
		brokers[i] = v + ":9092"
	}

	config := sarama.NewConfig()

	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = 0
	config.Producer.Retry.Max = 3

	mx.Lock()
	producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("failed to connect kafka: %v \n", err.Error())
	}
	mx.Unlock()
}

func Produce(topic string, body []byte) error {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(body),
	}

	_, _, err := producer.SendMessage(message)

	return err
}
