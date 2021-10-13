package middleware

import (
	"crypto/sha256"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/notification"
	log "github.com/sirupsen/logrus"
	"io"
)

type Obfuscate struct {
}

func (o *Obfuscate) Init() {
	log.Debug("Nothing to init in obfuscate middleware")
	return
}

func (o *Obfuscate) Name() string {
	return "Obfuscate"
}

func (o *Obfuscate) Do(event *notification.Event, reader io.ReadCloser, objectInfo minio.ObjectInfo) (io.ReadCloser, minio.ObjectInfo, error) {
	sum := sha256.Sum256([]byte(objectInfo.Key))
	objectInfo.Key = fmt.Sprintf("%x", sum)

	log.Debug("Old Object key ", event.S3.Object.Key, " New Object key ", objectInfo.Key)

	return reader, objectInfo, nil
}
