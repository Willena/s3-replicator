package handlers

import (
	"S3Replicator/queue"
	"fmt"
)

const AMQP_TYPE_NAME = "amqp"

type AMQPCommandLineConfig struct {
	Url       string `long:"url" description:"the url for the AMQP client" required:"false" env:"URL"`
	QueueName string `long:"queue-name" description:"the name of the AMQP queue" required:"false" env:"QUEUE_NAME"`
	Exchange  string `long:"exchange-type" description:"the name of the AMQP exchange" required:"false" env:" EXCHANGE_NAME"`
}

type AMQPEventHandler struct {
	Config              AMQPCommandLineConfig
	NotificationChannel *queue.EventChannelQueue
}

func (A *AMQPEventHandler) Close() error {
	return fmt.Errorf("")
}

func (A *AMQPEventHandler) GetHandlerName() string {
	return AMQP_TYPE_NAME
}

func (A *AMQPEventHandler) Init() error {
	//Do init !
	return nil
}
