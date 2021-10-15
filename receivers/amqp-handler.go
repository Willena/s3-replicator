package receivers

import (
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7/pkg/notification"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"github.com/willena/S3Replicator/queue"
	"math/rand"
)

const AMQP_TYPE_NAME = "amqp"

type AMQPCommandLineConfig struct {
	Url       string `long:"url" description:"the url for the AMQP client" required:"false" env:"URL"`
	QueueName string `long:"queue-name" description:"the name of the AMQP queue" required:"false" env:"QUEUE_NAME"`
	Exchange  string `long:"exchange" description:"the name of the AMQP exchange" required:"false" env:" EXCHANGE_NAME"`
}

type AMQPEventHandler struct {
	Config              AMQPCommandLineConfig
	NotificationChannel *queue.EventChannelQueue
	connection          *amqp.Connection
	channel             *amqp.Channel
}

func (A *AMQPEventHandler) Close() error {
	err := A.connection.Close()
	err = A.channel.Close()
	return err
}

func (A *AMQPEventHandler) GetHandlerName() string {
	return AMQP_TYPE_NAME
}

func (A *AMQPEventHandler) messageReceive() error {
	msgs, err := A.channel.Consume(
		A.Config.QueueName, // queue
		"S3Replicator-"+fmt.Sprintf("%04d", rand.Intn(10000)), // consumer
		true,  // auto ack
		true,  // exclusive
		false, // no local
		false, // no wait
		nil,   // args
	)

	if err != nil {
		log.Error(err)
		return err
	}

	for d := range msgs {
		log.Trace("Got new event from AMQP")
		ev := notification.Info{}
		err = json.Unmarshal(d.Body, &ev)
		if err != nil {
			log.Error("ERR: ", err)
		}

		//For each record in the message produce the event
		for _, record := range ev.Records {
			A.NotificationChannel.Queue(record)
			log.Debug("Record queued", record)
			log.Debug("queue size: ", A.NotificationChannel.Size())
		}

	}
	return nil
}

func (A *AMQPEventHandler) Init() error {
	var err error
	A.connection, err = amqp.Dial(A.Config.Url)
	if err != nil {
		return err
	}

	A.channel, err = A.connection.Channel()
	if err != nil {
		return err
	}

	err = A.channel.QueueBind(A.Config.QueueName, "", A.Config.Exchange, false, nil)
	if err != nil {
		return err
	}

	go A.messageReceive()

	return nil
}
