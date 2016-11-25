package executor

import (
	"github.com/Sirupsen/logrus"

	"github.com/Dataman-Cloud/baker/config"
	"github.com/Dataman-Cloud/baker/external/docker"
)

type ImagePushTask struct {
	Config    *config.DockerRegistry
	ImageName string
	WorkDir   string
}

// NewImagePushTask is a task to do build image and push image to registry
func NewImagePushTask(imageName, workDir string, config *config.DockerRegistry) *ImagePushTask {
	return &ImagePushTask{
		Config:    config,
		ImageName: imageName,
		WorkDir:   workDir,
	}
}

// DockerLogin
func (t *ImagePushTask) DockerLogin() error {
	config := t.Config
	registry := config.Address
	client := docker.NewDockerClient()
	err := client.DockerLogin(config.Username, config.Password,
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
	client := docker.NewDockerClient()
	err := client.DockerBuild(imageAddrAndName, t.WorkDir, "Dockerfile")
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
	client := docker.NewDockerClient()
	err := client.DockerPush(imageAddrAndName, registry)
	if err != nil {
		logrus.Error("error docker push image to the registry.")
		return err
	}
	return nil
}
