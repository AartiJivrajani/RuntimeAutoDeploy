package handlers

import (
	"RuntimeAutoDeploy/common"
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

func (routine *Status) addToStatusList(traceId string, status string) {
	var (
		res         string
		err         error
		statusList  []string
		jStatusList []byte
	)
	res, err = routine.redisConn.Get(fmt.Sprintf("%s-%s", common.TRACE_ID, traceId)).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"err":     err.Error(),
			"traceId": traceId,
			"status":  status,
		}).Error("error fetching value from redis")
		return
	}
	err = json.Unmarshal([]byte(res), &statusList)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("error unmarshalling the retrieved value from redis")
		return
	}
	statusList = append(statusList, status)
	jStatusList, _ = json.Marshal(statusList)
	routine.redisConn.Set(fmt.Sprintf("%s-%s", common.TRACE_ID, traceId), jStatusList, 0)
}
