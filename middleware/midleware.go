package middleware

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/notification"
	"io"
)

type MiddleWare interface {
	Init()
	Do(event *notification.Event, reader io.ReadCloser, objectInfo minio.ObjectInfo) (io.ReadCloser, minio.ObjectInfo, error)
	Name() string
}
