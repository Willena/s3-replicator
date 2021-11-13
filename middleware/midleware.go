package middleware

import (
	"github.com/minio/minio-go/v7/pkg/notification"
	"github.com/willena/S3Replicator/manifest"
	"io"
)

type MiddleWare interface {
	Init()
	DoOnCreate(event *notification.Event, readers []io.Reader, objectInfo []manifest.Item) ([]io.Reader, []manifest.Item, error)
	DoOnRemove(event *notification.Event, objectInfo []manifest.Item) ([]manifest.Item, error)
	Name() string
}
