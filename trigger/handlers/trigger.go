package handlers

import (
	"RuntimeAutoDeploy/common"
	"RuntimeAutoDeploy/config"
	"RuntimeAutoDeploy/generateK8S"
	"archive/tar"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	guuid "github.com/google/uuid"

	dockertypes "github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"github.com/go-git/go-git"
	log "github.com/sirupsen/logrus"
)

func Cleanup(dir string) error {
	// Delete the dir folder and all repos inside
	d, err := os.Open(dir)
	if err != nil {
		return nil
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("error reading directories in the current path")
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			log.WithFields(log.Fields{
				"err": err.Error(),
			}).Error("error deleting directories from the current path")
			return err
		}
	}
	return nil
}

func tarDirectory(dir string, buf io.Writer) error {
	// Tar the dir folder
	// Reference: https://gist.github.com/mimoo/25fc9716e0f1353791f5908f94d6e726
	tw := tar.NewWriter(buf)
	filepath.Walk(dir, func(file string, fi os.FileInfo, err error) error {
		// generate tar header

		//fmt.Println(file)
		//fmt.Println(fi.Name())

		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			log.Error("FileInfoHeader")
			return err
		}

		// must provide real name
		// (see https://golang.org/src/archive/tar/common.go?#L626)
		header.Name = filepath.ToSlash(file)

		// write header
		if err := tw.WriteHeader(header); err != nil {
			log.Error("WriteHeader")
			return err
		}
		// if not a dir, write file content
		if !fi.IsDir() {

			data, err := os.Open(file)
			if err != nil {
				return err
			}

			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}

		return nil
	})

	// produce tar
	if err := tw.Close(); err != nil {
		return err
	}

	return nil
}

func downloadGitRepo(ctx context.Context, gitrepo string) bool {

	common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
		fmt.Sprintf(common.STAGE_FORMAT,
			common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
			common.STAGE_STATUS_WIP,
			common.STAGE_GIT), false)

	_, err := git.PlainClone(common.GIT_BUILD_FOLDER, false, &git.CloneOptions{
		URL:      gitrepo,
		Progress: os.Stdout,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error cloning repo")

		common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
			fmt.Sprintf(common.STAGE_ERROR_FORMAT,
				common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
				common.STAGE_STATUS_ERROR,
				common.STAGE_GIT,
				err.Error()), false)
		
		err = Cleanup(common.GIT_BUILD_FOLDER)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("error clearing GIT_BUILD_FOLDER")
		}
		return false
	}

	pattern := common.GIT_BUILD_FOLDER + "Dockerfile*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}
	if len(matches) < 1 {
		// Dockerfile does not exist
		log.Error("git repo missing Dockerfile")

		common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
			fmt.Sprintf(common.STAGE_ERROR_FORMAT,
				common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
				common.STAGE_STATUS_ERROR,
				common.STAGE_GIT,
				"Missing Dockerfile"), false)

		err = Cleanup(common.GIT_BUILD_FOLDER)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("error clearing GIT_BUILD_FOLDER")
		}
		return false
	}
	common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
		fmt.Sprintf(common.STAGE_FORMAT,
			common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
			common.STAGE_STATUS_DONE,
			common.STAGE_GIT), false)
	return true
}

func buildDockerImage(ctx context.Context, path string, conf *config.Application) error {
	// https://stackoverflow.com/questions/38804313/build-docker-image-from-go-code
	//ctx := context.Background()

	common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
		fmt.Sprintf(common.STAGE_FORMAT,
			common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
			common.STAGE_STATUS_WIP,
			fmt.Sprintf(common.STAGE_BUILDING_DOCKER_IMAGE, conf.AppName)), false)

	cli, err := dockerclient.NewClientWithOpts(dockerclient.WithVersion("1.40")) // Max supported API version

	if err != nil {
		log.Fatal(err, " :unable to init client")
		common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
			fmt.Sprintf(common.STAGE_ERROR_FORMAT,
				common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
				common.STAGE_STATUS_ERROR,
				fmt.Sprintf(common.STAGE_BUILDING_DOCKER_IMAGE, conf.AppName),
				"unable to start the docker init, there's an issue with the docker client"), false)
		return err
	}

	buf := new(bytes.Buffer)

	err = tarDirectory(path, buf)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error(" :unable to tar directory")
		common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
			fmt.Sprintf(common.STAGE_ERROR_FORMAT,
				common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
				common.STAGE_STATUS_ERROR,
				fmt.Sprintf(common.STAGE_BUILDING_DOCKER_IMAGE, conf.AppName),
				"unable to tar the directory"), false)
		return err
	}

	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		dockertypes.ImageBuildOptions{
			NoCache:    true,
			Tags:       []string{fmt.Sprintf("%s:%s", config.UserConfig.Reg.Address, conf.AppName)},
			Context:    dockerFileTarReader,
			Dockerfile: fmt.Sprintf("%s%s", common.GIT_BUILD_FOLDER, conf.Dockerfile),
			Remove:     true})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error(" :unable to build docker image")
		common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
			fmt.Sprintf(common.STAGE_ERROR_FORMAT,
				common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
				common.STAGE_STATUS_ERROR,
				fmt.Sprintf(common.STAGE_BUILDING_DOCKER_IMAGE, conf.AppName),
				"error building the docker image"), false)
		return err
	}
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error(err, " :unable to read image build response")
		return err
	}

	common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
		fmt.Sprintf(common.STAGE_FORMAT,
			common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
			common.STAGE_STATUS_DONE,
			fmt.Sprintf(common.STAGE_BUILDING_DOCKER_IMAGE, conf.AppName)), false)

	authConfig := dockertypes.AuthConfig{
		Username: config.UserConfig.Reg.Username,
		Password: config.UserConfig.Reg.Password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	imagePushResponse, err := cli.ImagePush(ctx, fmt.Sprintf("%s:%s", config.UserConfig.Reg.Address, conf.AppName),
		dockertypes.ImagePushOptions{RegistryAuth: authStr})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer imagePushResponse.Close()
	_, err = io.Copy(os.Stdout, imagePushResponse)
	return nil
}

func createK8sArtefacts(ctx context.Context, conf *config.Application) {

	err := generateK8S.CreateDeployment(ctx, conf)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("error creating deployment")
	}
	err = generateK8S.CreateService(ctx, conf)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("error creating service")
	}
}

func startDeployment(ctx context.Context, userRequestConfig *common.RADConfig) bool {
	// Parse the config file
	// Download the git repository into local Trigger/build folder
	// Check for Dockerfile. If does not exist, quit
	// else build docker image and store within Trigger/images
	var (
		err error
	)
	fmt.Println(userRequestConfig.GitRepoLink)

	if !downloadGitRepo(ctx, userRequestConfig.GitRepoLink) {
		return false
	}
	// TODO: [Aarti]: Parse the config file first, then proceed to building the docker image
	err = config.ReadUserConfigFile(ctx)
	if err != nil {
		return false
	}

	for _, conf := range config.UserConfig.Applications {
		err = buildDockerImage(ctx, common.GIT_BUILD_FOLDER, conf)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("error building docker image")
			return false
		}
		createK8sArtefacts(ctx, conf)
	}
	endTime := fmt.Sprintf("%d", time.Now().Unix())
	common.AddToStatusList(fmt.Sprintf("%s-%s", common.END_TIMESTAMP, ctx.Value(common.TRACE_ID).(string)), endTime, true)
	
	err = Cleanup(common.GIT_BUILD_FOLDER)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error clearing GIT_BUILD_FOLDER")
	}
	//deploymentCompleteChan <- true
	return true
}

func RADTriggerHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("received trigger request")
	var (
		data      common.RADConfig
		err       error
		ctx       context.Context
		traceId   guuid.UUID
		startTime string
	)
	if r.Method != "POST" {
		log.Error("error. Received incorrect HTTP method. Expecting POST")
		return
	}
	ctx, _ = context.WithCancel(context.Background())
	// add the unique ID to the context
	traceId = guuid.New()
	ctx = context.WithValue(ctx, common.TRACE_ID, traceId.String())

	// add the start timestamp to the context
	startTime = fmt.Sprintf("%d", time.Now().Unix())
	// add this to the context as well
	common.AddToStatusList(fmt.Sprintf("%s-%s", common.START_TIMESTAMP, traceId.String()), startTime, true)

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("error decoding post body in the trigger handler")
		return
	}
	//err = os.Chmod("setup/rad_management_cluster.sh", 0700)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//output, err := exec.Command("/bin/sh",
	//	"setup/rad_management_cluster.sh").Output()
	//
	//if err != nil {
	//	log.WithFields(log.Fields{
	//		"err": err.Error(),
	//	}).Error("error bootstrapping current cluster as the management cluster")
	//	log.Error(string(output))
	//	return
	//}
	//log.Info(string(output))

	err = generateK8S.GetK8sClient(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("error fetching k8s client")
		return
	}

	go startDeployment(ctx, &data)
	

	// Write back the trace ID for the user os they can request
	// for the status
	_, _ = w.Write([]byte(traceId.String()))
	return
}
