package common

const (
	TRACE_ID         = "TRACE_ID"
	GIT_BUILD_FOLDER = "build/"

	// STAGE PRINT
	STAGE_FORMAT       = "[%s]: %s" // use this as fmt.Sprintf(STAGE_FORMAT, STAGE_STATUS_WIP, STAGE_GIT)
	STAGE_ERROR_FORMAT = "[%s]: [Stage] %s: [Error] %s"
	STAGE_STATUS_WIP   = "IN PROGRESS"
	STAGE_STATUS_DONE  = "COMPLETED"
	STAGE_STATUS_ERROR = "ERROR"

	// STAGES
	STAGE_GIT                   = "Git Repo Cloned"
	STAGE_COLLECTING_DOCKER_REQ = "Collecting build requirements"
	STAGE_BUILDING_DOCKER_IMAGE = "Building docker image"

	STAGE_CREATING_PODS = "Creating Kubernetes Artefacts"
)
