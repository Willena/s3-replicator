package middleware

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/notification"
	log "github.com/sirupsen/logrus"
	"io"
)

type Split struct {
	blockSize int64
}

func (s *Split) Init() {
	s.blockSize = 128
	return
}

func (s *Split) Do(event *notification.Event, reader io.ReadCloser, objectInfo minio.ObjectInfo) (io.ReadCloser, minio.ObjectInfo, error) {
	if objectInfo.Size > s.blockSize {
		log.Debug("The object is bigger than ", s.blockSize, ". It needs to be spliced")
	}
	return reader, objectInfo, nil
}

func (s *Split) Name() string {
	return "Split"
}
