package main

import (
    "fmt"
    "log"
    "net/http"
    "RuntimeAutoDeploy/Trigger/handlers"
    "RuntimeAutoDeploy/common"
)


func main() {
    trigger.Cleanup(common.GIT_BUILD_FOLDER)
    fmt.Println("User Preference Service Started")
    http.HandleFunc("/trigger",http.HandlerFunc(trigger.AppPreferencesHandler))
	err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err.Error())
    }
}

