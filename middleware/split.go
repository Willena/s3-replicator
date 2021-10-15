package middleware

import (
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/notification"
	log "github.com/sirupsen/logrus"
	"io"
	"math"
	"strconv"
)

const DefaultBlockSize128mb = 134217728

var splitLogger = log.WithField("middleware", "split")

type Split struct {
	BlockSize int64
}

func (s *Split) Init() {
	if s.BlockSize <= 0 {
		s.BlockSize = DefaultBlockSize128mb
	}

	splitLogger.Info("Each object will be cut into ", s.BlockSize, " bytes objects ! ")
	return
}

func (s *Split) Do(_ *notification.Event, readers []io.Reader, objectInfos []minio.ObjectInfo) ([]io.Reader, []minio.ObjectInfo, error) {

	if len(readers) != len(objectInfos) {
		return nil, nil, fmt.Errorf("should have the same number of readers and object infos")
	}

	allReaders := make([]io.Reader, 0)
	allObjectInfos := make([]minio.ObjectInfo, 0)

	for i, reader := range readers {

		initialObjectInfo := objectInfos[i]

		if initialObjectInfo.Size > s.BlockSize {
			splitLogger.Debug("The object is bigger (", initialObjectInfo.Size, " bytes)  than the max bloc size (", s.BlockSize, " bytes). It needs to be spliced")
			splitLogger.Debug(initialObjectInfo.Key, " Will be cut into ", math.Ceil(float64(initialObjectInfo.Size)/float64(s.BlockSize)), " parts")
		}

		bytesLeft := initialObjectInfo.Size
		it := 0
		for bytesLeft > 0 {
			//The reader will stop if there are no more bytes to be consumed event if it is less that a block long
			newReader := io.LimitReader(reader, s.BlockSize)
			newObjectInfo := objectInfos[i]
			if bytesLeft > s.BlockSize {
				newObjectInfo.Size = s.BlockSize
				bytesLeft -= s.BlockSize
			} else {
				newObjectInfo.Size = bytesLeft
				bytesLeft -= bytesLeft
			}

			newObjectInfo.Key = newObjectInfo.Key + "_" + strconv.Itoa(it)

			allReaders = append(allReaders, newReader)
			allObjectInfos = append(allObjectInfos, newObjectInfo)
			splitLogger.Debug("Processed part ", it, " of ", initialObjectInfo.Key, " part size: ", newObjectInfo.Size)
			it += 1
		}

	}

	return allReaders, allObjectInfos, nil
}

func (s *Split) Name() string {
	return "Split"
}
