package middleware

import (
	"crypto/sha256"
	"fmt"
	"github.com/minio/minio-go/v7/pkg/notification"
	log "github.com/sirupsen/logrus"
	"github.com/willena/S3Replicator/manifest"
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

func (o *Obfuscate) DoOnCreate(event *notification.Event, readers []io.Reader, objectInfos []manifest.Item) ([]io.Reader, []manifest.Item, error) {

	tmpInfos := make([]manifest.Item, len(objectInfos))

	for i, info := range objectInfos {
		oldKey := info.ObjectId
		sum := sha256.Sum256([]byte(info.ObjectId))
		info.ObjectId = fmt.Sprintf("%x", sum)
		obfuscateLogger.Debug("Old Object key ", oldKey, " New Object key ", info.ObjectId)
		info.ObjectInfo.Key = info.ObjectId
		tmpInfos[i] = info
	}

	return readers, tmpInfos, nil
}

func (c *Obfuscate) DoOnRemove(event *notification.Event, objectInfo []manifest.Item) ([]manifest.Item, error) {
	return objectInfo, nil
}
