package common

const (
	TRACE_ID        = "TRACE_ID"
	START_TIMESTAMP = "Start-TS"
	END_TIMESTAMP   = "End-TS"
	//GIT_BUILD_FOLDER = "/Users/aartij17/go/src/RuntimeAutoDeploy/buildRAD/"
	GIT_BUILD_FOLDER = "build/"
	USER_CONFIG_FILE = "config.json"

	// STAGE PRINT
	STAGE_FORMAT       = "[%s]: [%s]: %s" // use this as fmt.Sprintf(STAGE_FORMAT, STAGE_STATUS_WIP, STAGE_GIT)
	STAGE_ERROR_FORMAT = "[%s]: [%s]: [Stage] %s: [Error] %s"
	STAGE_STATUS_WIP   = "IN PROGRESS"
	STAGE_STATUS_DONE  = "COMPLETED"
	STAGE_STATUS_ERROR = "ERROR"

	// STAGES
	STAGE_GIT                   = "clone user git repository"
	STAGE_BUILDING_DOCKER_IMAGE = "build docker image[%s]"
	STAGE_PUSH_IMAGE            = "push docker image to repository"
	STAGE_READ_USER_CONFIG_FILE = "read user config file"

	STAGE_CREATING_DEPLOYMENT = "create kubernetes deployment[%s]"
	STAGE_CREATING_SERVICE    = "create kubernetes service[%s]"
	STAGE_K8S_BOOTSTRAP       = "start k8s client"
)
