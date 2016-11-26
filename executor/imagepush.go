package executor

import (
	"github.com/Sirupsen/logrus"

	"github.com/Dataman-Cloud/baker/config"
	"github.com/Dataman-Cloud/baker/external/docker"
)

type ImagePushTask struct {
	WorkDir   string
	ImageName string
	Config    *config.DockerRegistry
	Client    *docker.DockerClient
}

// NewImagePushTask is a task to do build image and push image to registry
func NewImagePushTask(workDir, imageName string, config *config.DockerRegistry) *ImagePushTask {
	return &ImagePushTask{
		Client:    docker.NewDockerClient(),
		Config:    config,
		ImageName: imageName,
		WorkDir:   workDir,
	}
}

// Create
func (t *ImagePushTask) Create(c *Collector) func() {
	return func() {
		c.TaskStatus <- StatusRunning
		var err error
		workDir := t.WorkDir
		imageName := t.ImageName
		dockerRegistry := t.Config
		imagePushTask := NewImagePushTask(workDir, imageName, dockerRegistry)

		// dockerLogin
		c.TaskStatus <- StatusDockerLoginStart
		err = imagePushTask.DockerLogin()
		if err != nil {
			logrus.Error("error execute docker login.")
			c.TaskMsg <- err.Error()
			return
		}
		c.TaskStatus <- StatusDockerLoginOK

		// dockerBuild
		c.TaskStatus <- StatusDockerBuildStart
		err = imagePushTask.DockerBuild()
		if err != nil {
			logrus.Error("error execute docker build.")
			c.TaskMsg <- err.Error()
			return
		}
		c.TaskStatus <- StatusDockerBuildOK

		// dockerPush
		c.TaskStatus <- StatusDockerPushStart
		err = imagePushTask.DockerPush()
		if err != nil {
			logrus.Error("error execute docker push.")
			c.TaskMsg <- err.Error()
			return
		}
		c.TaskStatus <- StatusDockerPushOK

		defer func() {
			if r := recover(); r != nil {
				c.TaskStatus <- StatusFailed
			} else {
				c.TaskStatus <- StatusFinished
			}
		}()
	}
}

// DockerLogin
func (t *ImagePushTask) DockerLogin() error {
	config := t.Config
	registry := config.Address
	err := t.Client.DockerLogin(config.Username, config.Password,
		config.Email, registry)
	if err != nil {
		logrus.Error("error docker login to the registry.")
		return err
	}
	return nil
}

// DockerBuild
func (t *ImagePushTask) DockerBuild() error {
	config := t.Config
	registry := config.Address
	repo := config.Repo
	imageAddrAndName := registry + "/" + repo + "/" + t.ImageName
	err := t.Client.DockerBuild(imageAddrAndName, t.WorkDir, "Dockerfile")
	if err != nil {
		logrus.Error("error build image from dockerfile.")
		return err
	}
	return nil
}

// DockerPush
func (t *ImagePushTask) DockerPush() error {
	config := t.Config
	registry := config.Address
	repo := config.Repo
	imageAddrAndName := registry + "/" + repo + "/" + t.ImageName
	err := t.Client.DockerPush(imageAddrAndName, registry)
	if err != nil {
		logrus.Error("error docker push image to the registry.")
		return err
	}
	return nil
}
