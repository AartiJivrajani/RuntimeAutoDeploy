package main

import (
	"RuntimeAutoDeploy/common"
	"RuntimeAutoDeploy/trigger/handlers"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func main() {
	log.Info("RAD Service Started")
	err := handlers.Cleanup(common.GIT_BUILD_FOLDER)
	if err != nil {
		log.Error("Exiting now.")
		return
	}
	go handlers.StartStatusService()

	http.HandleFunc("/trigger", handlers.RADTriggerHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
