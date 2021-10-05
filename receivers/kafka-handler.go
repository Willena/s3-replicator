package receivers

import (
	"S3Replicator/queue"
	"fmt"
)

const KAFKA_TYPE_NAME = "kafka"

type KafkaCommandLineConfig struct {
	BrokerList     []string          `long:"broker-list" description:"List of broker to connect to" required:"false" env:"BROKER_LIST"`
	Topics         []string          `long:"topics" description:"List of topic to listen" required:"false" env:"TOPICS"`
	ConsumerConfig map[string]string `long:"consumer-config" description:"Map of parameters for the kafka consumer" required:"false" env:"CONSUMER_CONFIG"`
}

type KafkaEventHandler struct {
	Config              KafkaCommandLineConfig
	NotificationChannel *queue.EventChannelQueue
}

func (k *KafkaEventHandler) GetHandlerName() string {
	return KAFKA_TYPE_NAME
}

func (k *KafkaEventHandler) Init() error {
	//Initialize, connect to Kafka, prepare consumer, ...
	return nil
}

func (k *KafkaEventHandler) Close() error {
	//close
	return fmt.Errorf("")
}
