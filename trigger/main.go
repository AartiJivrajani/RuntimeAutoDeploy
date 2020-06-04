package main

import (
	"RuntimeAutoDeploy/common"
	"RuntimeAutoDeploy/trigger/handlers"
	"net/http"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

func main() {
	var err error
	log.Info("RAD Service Started")
	err = handlers.Cleanup(common.GIT_BUILD_FOLDER)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error cleaning up")
		return
	}
	handlers.StartStatusService()

	http.HandleFunc("/trigger", handlers.RADTriggerHandler)
	http.HandleFunc("/status", handlers.RADStatusHandler)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			log.Info("Received an interrupt, stopping all connections...")
			//cancel()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
