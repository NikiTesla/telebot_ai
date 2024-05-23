package main

import (
	"telebotai/pkg/service"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	srv := service.New()
	if err := srv.Run(); err != nil {
		log.WithError(err).Fatal("running failed")
	}
}
