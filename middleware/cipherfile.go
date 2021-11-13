package middleware

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/minio/minio-go/v7/pkg/notification"
	"github.com/minio/sio"
	log "github.com/sirupsen/logrus"
	"github.com/willena/S3Replicator/manifest"
	"golang.org/x/crypto/hkdf"
	"io"
)

var cipherLogger = log.WithField("middleware", "cipher")

type CipherFile struct {
	masterkey []byte
}

func (c *CipherFile) Init() {
	cipherLogger.Info("Object will be ciphered")

	// the master key used to derive encryption keys
	// this key must be keep secret
	var err error
	//TODO: Make this configurable !
	c.masterkey, err = hex.DecodeString("000102030405060708090A0B0C0D0E0FF0E0D0C0B0A090807060504030201000")
	if err != nil {
		fmt.Printf("Cannot decode hex key: %v", err) // add error handling
		return
	}

}

func (c *CipherFile) Name() string {
	return "CipherFile"
}

func (c *CipherFile) DoOnRemove(event *notification.Event, objectInfo []manifest.Item) ([]manifest.Item, error) {
	return objectInfo, nil
}

func (c *CipherFile) DoOnCreate(event *notification.Event, readers []io.Reader, objectInfo []manifest.Item) ([]io.Reader, []manifest.Item, error) {

	if len(readers) != len(objectInfo) {
		return nil, nil, fmt.Errorf("the number of readers must the same as the objectinfos")
	}

	allReaders := make([]io.Reader, 0)
	allObjectInfos := make([]manifest.Item, 0)

	for i, object := range objectInfo {
		cipherLogger.Debug("Ciphering file ", object.ObjectId, " .... (create specific encrypted reader)")

		// generate a random nonce to derive an encryption key from the master key
		// this nonce must be saved to be able to decrypt the data again - it is not
		// required to keep it secret
		var nonce [32]byte
		if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
			fmt.Printf("Failed to read random data: %v", err) // add error handling
			return nil, nil, err
		}

		// derive an encryption key from the master key and the nonce
		var key [32]byte
		kdf := hkdf.New(sha256.New, c.masterkey, nonce[:], nil)
		if _, err := io.ReadFull(kdf, key[:]); err != nil {
			fmt.Printf("Failed to derive encryption key: %v", err) // add error handling
			return nil, nil, err
		}

		encrypted, err := sio.EncryptReader(readers[i], sio.Config{Key: key[:]})
		if err != nil {
			return nil, nil, err
		}

		finalSize, err := sio.EncryptedSize(object.Size)
		if err != nil {
			return nil, nil, err
		}
		cipherLogger.Debug("Ciphered file will be ", finalSize, " bytes long (original size: ", object.Size, ")")

		//Copy the object :)
		newObject := object
		newObject.Size = finalSize
		newObject.CipherKey = key[:]
		newObject.ObjectInfo.Size = int64(finalSize)

		allReaders = append(allReaders, encrypted)
		allObjectInfos = append(allObjectInfos, newObject)
	}

	return allReaders, allObjectInfos, nil
}
