package executor

const (
	// task status.
	StatusStarting = 0
	StatusRunning  = 1
	StatusFailed   = 2
	StatusExpired  = 3
	StatusFinished = 4
	// ImagePushTask status.
	StatusDockerLoginStart = 10
	StatusDockerLoginOK    = 11
	StatusDockerBuildStart = 12
	StatusDockerBuildOK    = 13
	StatusDockerPushStart  = 14
	StatusDockerPushOK     = 15
)

var TaskStatusEnum = map[int]string{
	StatusStarting: "Starting",
	StatusRunning:  "Running",
	StatusFailed:   "Failed",
	StatusExpired:  "Expired",
	StatusFinished: "Finished",
	// ImagePushTask status map.
	StatusDockerLoginStart: "DockerLoginStart",
	StatusDockerLoginOK:    "DockerLoginOK",
	StatusDockerBuildStart: "DockerBuildStart",
	StatusDockerBuildOK:    "DockerBuildOK",
	StatusDockerPushStart:  "DockerPushStart",
	StatusDockerPushOK:     "DockerPushOK",
}
