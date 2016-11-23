package executor

import (
	"github.com/Sirupsen/logrus"

	"github.com/Dataman-Cloud/baker/config"
	"github.com/Dataman-Cloud/baker/external/docker"
)

// Execute is a method to execute build image and push image to registry
// with limit a number of workers/executors.
func ImagePush(imageName, path string, config *config.DockerRegistry) error {
	registry := config.Address
	repo := config.Repo
	client := docker.NewDockerClient()
	err := client.DockerLogin(config.Username, config.Password,
		config.Email, registry)
	if err != nil {
		logrus.Fatal("error docker login to the registry.")
		return err
	}
	imageAddrAndName := registry + "/" + repo + "/" + imageName
	err = client.DockerBuild(imageAddrAndName, path, "Dockerfile")
	if err != nil {
		logrus.Fatal("error build image from dockerfile.")
		return err
	}
	err = client.DockerPush(imageAddrAndName, registry)
	if err != nil {
		logrus.Fatal("error docker push image to the registry.")
		return err
	}
	return nil
}
