package queue

import (
	"github.com/minio/minio-go/v7/pkg/notification"
	log "github.com/sirupsen/logrus"
)

type EventChannelQueue struct {
	notificationChannel chan notification.Event
	shouldFilter        func(event notification.Event) bool
}

func (e *EventChannelQueue) Dequeue() notification.Event {
	event := <-e.notificationChannel
	return event
}

func (e *EventChannelQueue) Queue(event notification.Event) error {
	if e.shouldFilter(event) {
		log.Debug("Event filtered: ", event.EventName)
		return nil
	}
	e.notificationChannel <- event
	return nil
}

func (e *EventChannelQueue) Size() int {
	return len(e.notificationChannel)
}

func NewEventChannelQueue(size uint, filterFunction func(event notification.Event) bool) *EventChannelQueue {
	return &EventChannelQueue{
		notificationChannel: make(chan notification.Event, size),
		shouldFilter:        filterFunction,
	}
}
