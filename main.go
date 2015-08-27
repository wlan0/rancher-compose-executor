package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/go-machine-service/events"
	"github.com/rancher/rancher-compose-executor/handlers"
)

var (
	GITCOMMIT = "HEAD"
)

func main() {

	log.WithFields(log.Fields{
		"gitcommit": GITCOMMIT,
	}).Info("Starting rancher-compose-executor")

	eventHandlers := map[string]events.EventHandler{
		"environment.create": handlers.CreateEnvironment,
		"ping":               handlers.PingNoOp,
	}

	apiUrl := os.Getenv("RANCHER_URL")
	accessKey := os.Getenv("RANCHER_ACCESS_KEY")
	secretKey := os.Getenv("RANCHER_SECRET_KEY")

	router, err := events.NewEventRouter("rancher-compose-executor", 2000, apiUrl, accessKey, secretKey, nil, eventHandlers, "environment", 10)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to create event router")
	} else {
		err := router.Start(nil)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Unable to start event router")
		}
	}
	log.Info("Exiting rancher-compose-executor")
}
