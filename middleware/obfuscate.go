package middleware

import (
	"crypto/sha256"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/notification"
	log "github.com/sirupsen/logrus"
	"io"
)

var obfuscateLogger = log.WithField("middleware", "obfuscate")

type Obfuscate struct {
}

func (o *Obfuscate) Init() {
	obfuscateLogger.Info("Object names will be obfuscated")
	return
}

func (o *Obfuscate) Name() string {
	return "Obfuscate"
}

func (o *Obfuscate) Do(event *notification.Event, readers []io.Reader, objectInfos []minio.ObjectInfo) ([]io.Reader, []minio.ObjectInfo, error) {

	tmpInfos := make([]minio.ObjectInfo, len(objectInfos))

	for i, info := range objectInfos {
		oldKey := info.Key
		sum := sha256.Sum256([]byte(info.Key))
		info.Key = fmt.Sprintf("%x", sum)
		obfuscateLogger.Debug("Old Object key ", oldKey, " New Object key ", info.Key)
		tmpInfos[i] = info
	}

	return readers, tmpInfos, nil
}
