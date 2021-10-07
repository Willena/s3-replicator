package eventProcessor

import (
	"S3Replicator/config"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/notification"
	log "github.com/sirupsen/logrus"
	"net/url"
)

type BasicProcessor struct {
	Source      config.S3Configuration
	Destination config.S3Configuration
	srcClient   *minio.Client
	dstClient   *minio.Client
}

func (receiver *BasicProcessor) Init() error {
	log.Info("Initializing main event processor")
	//Create two clients: source client, destclient
	src, err := createMinioClient(receiver.Source)
	if err != nil {
		return err
	}
	receiver.srcClient = src
	log.Info("Created source S3 client")

	dest, err := createMinioClient(receiver.Destination)
	if err != nil {
		return err
	}
	receiver.dstClient = dest
	log.Info("Created dest S3 client")

	return nil
}

func createMinioClient(configuration config.S3Configuration) (*minio.Client, error) {
	minioUrl, err := url.Parse(configuration.Endpoint)
	if err != nil {
		return nil, err
	}

	return minio.New(minioUrl.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(configuration.AccessKey, configuration.SecretKey, ""),
		Secure: minioUrl.Scheme == "https",
		Region: configuration.Region,
	})
}

func (receiver *BasicProcessor) processPostPut(event *notification.Event) error {
	log.Debug("Processing create event... ")

	objectReader, err := receiver.srcClient.GetObject(context.Background(), event.S3.Bucket.Name, event.S3.Object.Key, minio.GetObjectOptions{})
	defer objectReader.Close()
	if err != nil {
		log.Error("Could not get source object from bucket", event.S3.Bucket.Name, "with name", event.S3.Object.Key, "; Err: ", err)
		return err
	}

	objectInfo, err := objectReader.Stat()
	if err != nil {
		log.Error("Could not get Stats for ", event.S3.Object.Key, "; Err: ", err)
		return err
	}

	log.Debug("Source Object is ", objectInfo.Size, " bytes long ! ")

	uploadInfo, err := receiver.dstClient.PutObject(context.Background(), receiver.Destination.Bucket, event.S3.Object.Key, objectReader, objectInfo.Size, minio.PutObjectOptions{
		StorageClass: receiver.Destination.Class,
		UserMetadata: objectInfo.UserMetadata,
		UserTags:     objectInfo.UserTags,
		ContentType:  objectInfo.ContentType,
	})

	if err != nil {
		log.Error("Could not send object !", err)
		return err
	}
	log.Debug("Key: ", uploadInfo.Key, " Uploaded ", uploadInfo.Size, " byes to destination S3 ", uploadInfo.ETag)
	return err
}

func (receiver *BasicProcessor) readKey(event *notification.Event) {
	key, err := url.QueryUnescape(event.S3.Object.Key)
	if err != nil {
		log.Warn("Could not URLDecode key ! Using raw key")
		key = event.S3.Object.Key
	}
	event.S3.Object.Key = key
	log.Debug("Key is : '", event.S3.Object.Key, "'")
}

func (receiver *BasicProcessor) ProcessEvent(event *notification.Event) error {

	receiver.readKey(event)

	log.Info("Processing event with name: ", event.EventName, " Key(unescaped): ", event.S3.Object.Key, " bucket: ", event.S3.Bucket.Name)
	switch event.EventName {
	case notification.ObjectCreatedPut:
		return receiver.processPostPut(event)
	case notification.ObjectCreatedCopy:
		return receiver.processPostPut(event)
	case notification.ObjectCreatedPost:
		return receiver.processPostPut(event)
	case notification.ObjectCreatedCompleteMultipartUpload:
		return receiver.processPostPut(event)
	case notification.ObjectRemovedDelete:
		err := receiver.dstClient.RemoveObject(context.Background(), receiver.Destination.Bucket, event.S3.Object.Key, minio.RemoveObjectOptions{
			VersionID: event.S3.Object.VersionID,
		})
		if err != nil {
			log.Error("Could not remove object ", err)
			return err
		}
		log.Info("Key: ", event.S3.Object.Key, " Object removed !")
	default:
		log.Trace(event.EventName + " Not implemented ")
		return nil
	}
	return nil
}
