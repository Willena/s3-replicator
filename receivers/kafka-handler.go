package receivers

import (
	"S3Replicator/queue"
	"context"
	"encoding/json"
	"github.com/minio/minio-go/v7/pkg/notification"
	kafka "github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

const KAFKA_TYPE_NAME = "kafka"

type KafkaCommandLineConfig struct {
	BrokerList     []string          `long:"broker-list" description:"List of broker to connect to" required:"false" env:"BROKER_LIST"`
	Topics         []string          `long:"topics" description:"Topics to listen" required:"false" env:"TOPICs"`
	ConsumerConfig map[string]string `long:"consumer-config" description:"Map of parameters for the kafka consumer" required:"false" env:"CONSUMER_CONFIG"`
	GroupId        string            `long:"group-id" description:"Consumer group id" required:"true" env:"GROUP_ID"`
}

type KafkaEventHandler struct {
	Config              KafkaCommandLineConfig
	NotificationChannel *queue.EventChannelQueue
	kafkaReader         *kafka.Reader
}

func (k *KafkaEventHandler) GetHandlerName() string {
	return KAFKA_TYPE_NAME
}

func (k *KafkaEventHandler) receiveMessages() {
	ctx := context.Background()
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := k.kafkaReader.ReadMessage(ctx)
		if err != nil {
			log.Error("Error while reading message from Kafka... ")
		}

		ev := notification.Info{}
		err = json.Unmarshal(msg.Value, &ev)
		if err != nil {
			log.Error("ERR: ", err)
		}

		//For each record in the message produce the event
		for _, record := range ev.Records {
			k.NotificationChannel.Queue(record)
			log.Debug("Record queued", record)
			log.Debug("queue size: ", k.NotificationChannel.Size())
		}

	}
}

func (k *KafkaEventHandler) Init() error {
	//Initialize, connect to Kafka, prepare consumer, ...
	k.kafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     k.Config.BrokerList,
		GroupTopics: k.Config.Topics,
		GroupID:     k.Config.GroupId,
	})

	log.Info("Starting Kafka consumer... Topics: ", k.Config.Topics, " GroupId: ", k.Config.GroupId, " Borkers: ", k.Config.BrokerList)
	go k.receiveMessages()

	return nil
}

func (k *KafkaEventHandler) Close() error {
	return k.kafkaReader.Close()
}
