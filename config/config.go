package config

import (
	"RuntimeAutoDeploy/common"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var UserConfig *Config

type Config struct {
	Applications []*Application `json:"applications"`
	Reg          *Registry      `json:"registry"`
}

type Application struct {
	AppName      string `json:"application_name"`
	ReplicaCount int    `json:"replica_count"`
	Dockerfile   string `json:"dockerfile"`
	Port         int    `json:"port"`
}

type Registry struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func ReadUserConfigFile(ctx context.Context) error {
	var (
		err  error
		data []byte
	)
	common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
		fmt.Sprintf(common.STAGE_FORMAT,
			common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
			common.STAGE_STATUS_WIP,
			common.STAGE_READ_USER_CONFIG_FILE),
		false)

	configFilePath := fmt.Sprintf("%s%s", common.GIT_BUILD_FOLDER, common.USER_CONFIG_FILE)
	_, err = os.Stat(configFilePath)
	if err != nil && os.IsNotExist(err) {
		log.WithFields(log.Fields{
			"path": configFilePath,
		}).Error("user config file not found in the repository")

		common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
			fmt.Sprintf(common.STAGE_ERROR_FORMAT,
				common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
				common.STAGE_STATUS_ERROR,
				common.STAGE_READ_USER_CONFIG_FILE,
				"user config file(config.json) not found in the repository"),
			false)
		return err
	}

	data, err = ioutil.ReadFile(configFilePath)
	err = json.Unmarshal(data, &UserConfig)
	common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
		fmt.Sprintf(common.STAGE_FORMAT,
			common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
			common.STAGE_STATUS_DONE,
			common.STAGE_READ_USER_CONFIG_FILE),
		false)

	log.WithFields(log.Fields{
		"config": UserConfig.Applications[0].AppName,
	}).Info("config read from user file")
	return nil
}
