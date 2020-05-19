package handlers

import (
	"gopkg.in/redis.v5"
)

var (
	statusRoutine *Status
)

type Status struct {
	redisConn  *redis.Client
	statusList []string
}

func StartStatusService() {
	statusRoutine = &Status{
		redisConn:  nil,
		statusList: make([]string, 0),
	}
	statusRoutine.redisConn = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
