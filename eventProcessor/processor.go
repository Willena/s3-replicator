package eventProcessor

import (
	"github.com/minio/minio-go/v7/pkg/notification"
	"github.com/willena/S3Replicator/manifest"
	"io"
)

type Processor interface {
	Init() error
	ProcessEvent(event *notification.Event) error
	sendObject(event *notification.Event, reader io.Reader, objectInfo manifest.Item) (manifest.Item, error)
}
