package main

import (
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"github.com/willena/S3Replicator/config"
	"os"
)

func main() {

	log.SetLevel(log.DebugLevel)

	var opts config.Config
	log.Info("Starting S3Replicator...")
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Error(err)
		return
	}

	replicator, err := NewReplicator(opts)
	if err != nil || replicator == nil {
		log.Error("Could not initialize Replicate", err)
		return
	}

	replicator.Start()

}
