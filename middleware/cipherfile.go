package middleware

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/notification"
	log "github.com/sirupsen/logrus"
	"io"
)

type CipherFile struct {
}

func (c *CipherFile) Init() {
	log.Debug("Cipher middleware")
}

func (c *CipherFile) Name() string {
	return "CipherFile"
}

func (c *CipherFile) Do(event *notification.Event, reader io.ReadCloser, objectInfo minio.ObjectInfo) (io.ReadCloser, minio.ObjectInfo, error) {
	log.Debug("[WIP] Ciphering file ", event.S3.Object.Key, "....")
	return reader, objectInfo, nil
}
