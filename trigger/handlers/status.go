package handlers

import (
	"RuntimeAutoDeploy/common"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gopkg.in/redis.v5"
)

var (
	StatusRoutine *Status
)

type Status struct {
	redisConn  *redis.Client
	statusList []string
}

func RADStatusHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("received status request")
	// When the user pings this, return the status
	var (
		err        error
		statusList []string
	)
	if r.Method != "GET" {
		log.Error("error. Received incorrect HTTP method. Expecting GET")
		return
	}
	keys, OK := r.URL.Query()[common.TRACE_ID]
	if !OK || len(keys[0]) < 1 {
		log.Println("Url Param 'TRACE_ID' is missing")
		return
	}
	traceId := r.FormValue(common.TRACE_ID)

	// Fetch the status list from REDIS
	statusList, err = common.GetStatusList(traceId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal error, please try again"))
		return
	}

	for i := range statusList {
		_, _ = w.Write([]byte(statusList[i]))
		_, _ = w.Write([]byte("\n"))
	}
	startTimestamp, _ := common.GetStatusList(fmt.Sprintf("%s-%s", common.START_TIMESTAMP, traceId))
	endTimestamp, _ := common.GetStatusList(fmt.Sprintf("%s-%s", common.END_TIMESTAMP, traceId))

	_, _ = w.Write([]byte(fmt.Sprintf("start time: %s\n",
		common.GetTimestampFormat(startTimestamp[0], "", ""))))

	if endTimestamp != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("end time: %s\n",
			common.GetTimestampFormat(endTimestamp[0], "", ""))))

		_, _ = w.Write([]byte(fmt.Sprintf("time taken: %s seconds\n",
			common.GetTimestampFormat(startTimestamp[0], endTimestamp[0], "diff"))))
		return
	} else {
		_, _ = w.Write([]byte(fmt.Sprintf("end timestamp cannot be found, either the pipeline is running, OR there has been an error")))
		return
	}
}

func StartStatusService() {
	common.ConnectToRedis()
	StatusRoutine = &Status{
		redisConn:  common.RedisConn,
		statusList: make([]string, 0),
	}
}
