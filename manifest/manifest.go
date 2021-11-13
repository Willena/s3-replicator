package manifest

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/notification"
	log "github.com/sirupsen/logrus"
	"os"
)

var currentManifestObject *Manifest = nil

const (
	S3Storage                = "s3"
	FileStorage              = "file"
	createStructureStatement = `
create table if not exists objects
(
    id         TEXT not null,
    sourceName TEXT not null,
    versionId    TEXT,
    parts      INTEGER default 1 not null,
    part       INTEGER default 1 not null,
    key        BLOB,
    size       INTEGER default 0 not null,
    constraint objects_pk
        primary key (id, sourceName, part)
);
create unique index if not exists objects_id_uindex
	on objects (id);
`
	createNewRecordStatement = `INSERT INTO objects (id, sourceName,version, parts, part, key, size) VALUES (?, ?, ?, ?, ?, ?)`
)

type FileParts struct {
	File    string
	Objects []Item
}

type Item struct {
	ObjectId     string
	OriginalName string
	Parts        uint64
	Part         uint64
	CipherKey    []byte
	Size         uint64
	VersionId    string
	ObjectInfo   minio.ObjectInfo
}

func (i Item) Commit() {
	Get().eventChannel <- i
}

type Manifest struct {
	manifestPath string
	db           *sql.DB
	eventChannel chan Item
}

func Get() *Manifest {
	if currentManifestObject == nil {
		currentManifestObject = &Manifest{}
		currentManifestObject.init()
		log.Debug("Created new Manifest object")
	}
	return currentManifestObject
}

func (m *Manifest) init() {
	m.manifestPath = "./foo.db"
	os.Remove(m.manifestPath)

	var err error
	m.db, err = sql.Open("sqlite3", "./foo.db")
	//	defer db.Close() do not close :)
	if err != nil {
		log.Fatal(err)
	}

	_, err = m.db.Exec(createStructureStatement)
	if err != nil {
		log.Error("Could not create initial tables: ", err)
		return
	}

	m.eventChannel = make(chan Item)
	go m.readItems()
}

func (m *Manifest) Put(item Item) {
	log.Debug("Putting item in db...")

	tx, err := m.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(createNewRecordStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.ObjectId, item.OriginalName, item.Parts, item.Part, item.CipherKey, item.Size)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}

func (m *Manifest) GetObject(objectName string) Item {
	return Item{}
}

func (m *Manifest) GetPartsFromFileName(fileName string) FileParts {
	return FileParts{}
}

func (m *Manifest) GetPartsFromObjectId(objectId string) FileParts {
	item := m.GetObject(objectId)
	return m.GetPartsFromFileName(item.OriginalName)
}

func (m *Manifest) readItems() {
	for item := range m.eventChannel {
		m.Put(item)
	}
}

func WrapObjectInfo(info minio.ObjectInfo) Item {
	return Item{
		ObjectId:     info.Key,
		OriginalName: info.Key,
		Parts:        1,
		Part:         1,
		CipherKey:    nil,
		Size:         uint64(info.Size),
		ObjectInfo:   minio.ObjectInfo{},
	}
}

func WrapEvent(event *notification.Event) Item {
	return Item{
		ObjectId:     event.S3.Object.Key,
		OriginalName: event.S3.Object.Key,
		Parts:        1,
		Part:         1,
		CipherKey:    nil,
		Size:         uint64(event.S3.Object.Size),
		VersionId:    event.S3.Object.VersionID,
		ObjectInfo:   minio.ObjectInfo{},
	}
}
