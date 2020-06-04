package common

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/redis.v5"
)

var (
	RedisConn *redis.Client
)

func ConnectToRedis() {
	RedisConn = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func AddToStatusList(traceId string, status string, firstAdd bool) {
	var (
		res         string
		err         error
		statusList  []string
		jStatusList []byte
	)
	if firstAdd {
		statusList = make([]string, 0)
		statusList = append(statusList, status)
		jStatusList, _ = json.Marshal(statusList)
		RedisConn.Set(fmt.Sprintf("%s-%s", TRACE_ID, traceId), jStatusList, 0)
		return
	}
	res, err = RedisConn.Get(fmt.Sprintf("%s-%s", TRACE_ID, traceId)).Result()
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
	RedisConn.Set(fmt.Sprintf("%s-%s", TRACE_ID, traceId), jStatusList, 0)
}

func GetStatusList(traceId string) ([]string, error) {
	var (
		res        string
		err        error
		statusList []string
	)
	res, err = RedisConn.Get(fmt.Sprintf("%s-%s", TRACE_ID, traceId)).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"err":     err.Error(),
			"traceId": traceId,
		}).Error("error fetching value from redis")
		return nil, fmt.Errorf("%s", "error fetching value from redis")
	}
	err = json.Unmarshal([]byte(res), &statusList)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("error unmarshalling the retrieved value from redis")
		return nil, fmt.Errorf("%s", "error unmarshalling retreived value from redis")
	}
	return statusList, nil
}

func GetTimestampFormat(st string, et string, op string) string {
	switch op {
	case "diff":
		iSt, _ := strconv.ParseInt(st, 10, 64)
		stParse := time.Unix(iSt, 0)
		iEt, _ := strconv.ParseInt(et, 10, 64)
		etParse := time.Unix(iEt, 0)
		diff := etParse.Sub(stParse).Seconds()
		return fmt.Sprintf("%f", diff)
	default:
		iSt, _ := strconv.ParseInt(st, 10, 64)
		stParse := time.Unix(iSt, 0)
		return stParse.Format(time.ANSIC)
		//return fmt.Sprintf("%d", stParse)
	}
}
