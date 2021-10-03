package eventProcessor

import (
	"S3Replicator/poolWorker"
	"github.com/minio/minio-go/v7/pkg/notification"
	log "github.com/sirupsen/logrus"
)

type MultiThreadProcessor struct {
	BasicProcessor
	WorkerNumber uint
	processPool  *poolWorker.Pool
}

type job struct {
	processor *BasicProcessor
	event     *notification.Event
}

func (j *job) Do() {
	for i := 0; i < 50; i++ {
		if j.processor.ProcessEvent(j.event) == nil {
			return
		}
		log.Warn("Retrying ", i, "/50: ", j.event.EventName, " Key: ", j.event.S3.Object.Key)
	}
}

func (receiver *MultiThreadProcessor) Init() error {
	if receiver.WorkerNumber <= 0 {
		receiver.WorkerNumber = 1
	}
	err := receiver.BasicProcessor.Init()
	receiver.processPool = poolWorker.NewWorkerPool(receiver.WorkerNumber)
	receiver.processPool.Start()
	log.Info("Started a pool of ", receiver.WorkerNumber, " workers... ")
	return err
}

func (receiver *MultiThreadProcessor) ProcessEvent(event *notification.Event) error {
	receiver.processPool.Submit(&job{processor: &receiver.BasicProcessor, event: event})
	return nil
}
