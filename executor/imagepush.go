package executor

import (
	"github.com/Sirupsen/logrus"

	"github.com/Dataman-Cloud/baker/config"
	"github.com/Dataman-Cloud/baker/external/docker"
)

type ImagePush struct {
	WorkDir   string
	ImageName string
	Config    *config.DockerRegistry
	Client    *docker.DockerClient
}

// NewImagePush is a task to do build image and push image to registry
func NewImagePush(workDir, imageName string, config *config.DockerRegistry) *ImagePush {
	return &ImagePush{
		Client:    docker.NewDockerClient(),
		Config:    config,
		ImageName: imageName,
		WorkDir:   workDir,
	}
}

// Before
func (t *ImagePush) Before(c *Collector) func() {
	return func() {
		c.TaskStats <- &TaskStats{Code: StatusRunning}
	}
}

// After
func (t *ImagePush) After(c *Collector) func() {
	return func() {
		if r := recover(); r != nil {
			c.TaskStats <- &TaskStats{Code: StatusFailed}
		} else {
			c.TaskStats <- &TaskStats{Code: StatusFinished}
		}
	}
}

// DockerLogin
func (t *ImagePush) DockerLogin(c *Collector) func() {
	return func() {
		c.TaskStats <- &TaskStats{Code: StatusDockerLoginStart}
		config := t.Config
		registry := config.Address
		err := t.Client.DockerLogin(config.Username, config.Password,
			config.Email, registry)
		if err != nil {
			logrus.Error("error docker login to the registry.")
			c.TaskStats <- &TaskStats{Code: StatusFailed, Message: err.Error()}
		}
		c.TaskStats <- &TaskStats{Code: StatusDockerLoginStart}
	}
}

// DockerBuild
func (t *ImagePush) DockerBuild(c *Collector) func() {
	return func() {
		c.TaskStats <- &TaskStats{Code: StatusDockerBuildStart}
		config := t.Config
		registry := config.Address
		repo := config.Repo
		imageAddrAndName := registry + "/" + repo + "/" + t.ImageName
		err := t.Client.DockerBuild(imageAddrAndName, t.WorkDir, "Dockerfile")
		if err != nil {
			logrus.Error("error build image from dockerfile.")
			c.TaskStats <- &TaskStats{Code: StatusFailed, Message: err.Error()}
		}
		c.TaskStats <- &TaskStats{Code: StatusDockerBuildOK}
	}
}

// DockerPush
func (t *ImagePush) DockerPush(c *Collector) func() {
	return func() {
		c.TaskStats <- &TaskStats{Code: StatusDockerPushStart}
		config := t.Config
		registry := config.Address
		repo := config.Repo
		imageAddrAndName := registry + "/" + repo + "/" + t.ImageName
		err := t.Client.DockerPush(imageAddrAndName, registry)
		if err != nil {
			logrus.Error("error docker push image to the registry.")
			c.TaskStats <- &TaskStats{Code: StatusFailed, Message: err.Error()}
		}
		c.TaskStats <- &TaskStats{Code: StatusDockerPushOK}
	}
}
