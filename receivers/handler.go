package receivers

import "io"

type S3EventReceiver interface {
	io.Closer
	GetHandlerName() string
}

type InitializableEventReceiver interface {
	S3EventReceiver
	Init() error
}
