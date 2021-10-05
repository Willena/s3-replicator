package config

import "S3Replicator/receivers"

type S3Configuration struct {
	Endpoint  string `long:"endpoint" description:"URL to the S3" required:"true" env:"ENDPOINT"`
	Bucket    string `long:"bucket" description:"S3 bucket" required:"true" env:"BUCKET"`
	Class     string `long:"class" description:"S3 Storage Class" env:"CLASS" default:"STANDARD"`
	AccessKey string `long:"access-key" description:"S3 Storage Access Key" required:"true" env:"ACCESS_KEY"`
	SecretKey string `long:"secret-key" description:"S3 Storage Secret Key" required:"true" env:"SECRET_KEY"`
	Region    string `long:"region" description:"S3 Storage Region" required:"false" env:"REGION" default:"us-west"`
}

type Config struct {
	S3 struct {
		Source      S3Configuration `group:"S3 Source Configuration" namespace:"source" env-namespace:"SOURCE" required:"true"`
		Destination S3Configuration `group:"S3 Destination Configuration" namespace:"destination" env-namespace:"DESTINATION" required:"true"`
	} `group:"S3 Configuration" namespace:"s3" env-namespace:"S3"`
	AMQP              receivers.AMQPCommandLineConfig  `group:"AMQP Configuration" namespace:"amqp" env-namespace:"AMQP"`
	Kafka             receivers.KafkaCommandLineConfig `group:"Kafka Configuration" namespace:"kafka" env-namespace:"KAFKA"`
	Http              receivers.HTTPCommandLineConfig  `group:"Http Server Configuration" namespace:"http" env-namespace:"HTTP"`
	EventSource       string                           `long:"queue-type" choice:"kafka" choice:"amqp" choice:"http" required:"true" default:"http"`
	InternalQueueSize uint                             `long:"internal-queue-size" default:"5000" required:"true"`
}
