package middleware

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/notification"
	"io"
)

type MiddleWare interface {
	Init()
	Do(event *notification.Event, readers []io.Reader, objectInfo []minio.ObjectInfo) ([]io.Reader, []minio.ObjectInfo, error)
	Name() string
}
