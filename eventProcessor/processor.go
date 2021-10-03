package eventProcessor

import (
	"github.com/minio/minio-go/v7/pkg/notification"
)

type Processor interface {
	Init() error
	ProcessEvent(event *notification.Event) error
}
