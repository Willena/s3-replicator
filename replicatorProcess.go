package main

import (
	"S3Replicator/config"
	eventProcessor "S3Replicator/eventProcessor"
	"S3Replicator/middleware"
	"S3Replicator/queue"
	"S3Replicator/receivers"
	"fmt"
	"github.com/minio/minio-go/v7/pkg/notification"
	log "github.com/sirupsen/logrus"
	"strings"
)

func HandlerFromType(handlerType string, config config.Config, notificationChannel *queue.EventChannelQueue) (receivers.S3EventReceiver, error) {
	var ev receivers.InitializableEventReceiver

	switch handlerType {
	case receivers.AMQP_TYPE_NAME:
		ev = &receivers.AMQPEventHandler{Config: config.AMQP, NotificationChannel: notificationChannel}
		break
	case receivers.HTTP_TYPE_NAME:
		ev = &receivers.HTTPEventHandler{Config: config.Http, NotificationChannel: notificationChannel}
		break
	case receivers.KAFKA_TYPE_NAME:
		ev = &receivers.KafkaEventHandler{Config: config.Kafka, NotificationChannel: notificationChannel}
		break
	default:
		return nil, fmt.Errorf("%s is not a recognized handler type", handlerType)
	}

	//Initialize the collector...
	err := ev.Init()
	if err != nil {
		return nil, err
	}

	return ev.(receivers.S3EventReceiver), nil

}

type Replicator struct {
	handler            receivers.S3EventReceiver
	notificationChanel *queue.EventChannelQueue
	stop               bool
	eventProcessor     eventProcessor.Processor
}

func (receiver *Replicator) processEventThread() {
	log.Debug("Event receiver started ")
	for !receiver.stop {
		event := receiver.notificationChanel.Dequeue()
		var err error
		for i := 0; i < 50; i++ {
			err = receiver.eventProcessor.ProcessEvent(&event)
			if err == nil {
				break
			}
			log.Error("Retry", i+1, "/50 : ", err)
		}
	}
}

func (receiver *Replicator) Start() {
	receiver.processEventThread()
}

func (receiver *Replicator) Stop() {
	err := receiver.handler.Close()
	if err != nil {
		log.Error("Error while stopping", err)
	}
}

func shouldFilterEvent(event notification.Event) bool {
	if strings.HasPrefix(event.EventName, "s3:ObjectCreated") || strings.HasPrefix(event.EventName, "s3:ObjectRemoved") {
		return false
	}
	return true
}

func NewReplicator(config config.Config) (*Replicator, error) {
	channel := queue.NewEventChannelQueue(config.InternalQueueSize, shouldFilterEvent)

	handler, err := HandlerFromType(config.EventSource, config, channel)
	if err != nil {
		return nil, err
	}

	processor := &eventProcessor.MultiThreadProcessor{Processor: &eventProcessor.BasicProcessor{
		Source:      config.S3.Source,
		Destination: config.S3.Destination,
		MiddleWares: []middleware.MiddleWare{&middleware.Obfuscate{}, &middleware.CipherFile{}},
	}, WorkerNumber: 30}

	err = processor.Init()
	if err != nil {
		return nil, err
	}

	r := &Replicator{
		handler:            handler,
		notificationChanel: channel,
		stop:               false,
		eventProcessor:     processor,
	}

	return r, nil
}
