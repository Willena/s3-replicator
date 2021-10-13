package eventProcessor

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/notification"
	"io"
)

type Processor interface {
	Init() error
	ProcessEvent(event *notification.Event) error
	sendObject(event *notification.Event, reader io.Reader, objectInfo minio.ObjectInfo) error
}
