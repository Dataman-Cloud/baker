package docker

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

type DockerClient struct {
	Client *docker.Client
}

func NewDockerClient() *DockerClient {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		logrus.Fatal(err)
	}
	return &DockerClient{
		Client: client,
	}
}

// Dockerbuild is a function to create docker image from Dockerfile.
func (c *DockerClient) DockerBuild(name, context, dockerfile string) error {
	auth := loadRegistryAuthFile()
	client := c.Client
	err := client.BuildImage(docker.BuildImageOptions{
		Name:         name,
		Dockerfile:   dockerfile,
		OutputStream: bytes.NewBuffer(nil),
		ContextDir:   context,
		AuthConfigs:  *auth,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	return err
}

// DockerLogin is a function to docker login to the regitstry.
func (c *DockerClient) DockerLogin(username, password, email, registry string) error {
	err := generateRegistryAuthFile(username, password, email, registry)
	if err != nil {
		logrus.Fatal(err)
	}
	return err
}

// DockerPush is a function to docker push image to the registry.
func (c *DockerClient) DockerPush(image, registry string) error {
	dockerAuthConfiguration := getDockerAuthConfigurationFromFile(registry)
	client := c.Client
	err := client.PushImage(docker.PushImageOptions{
		Name:         image,
		OutputStream: bytes.NewBuffer(nil),
	}, *dockerAuthConfiguration)
	if err != nil {
		logrus.Fatal(err)
	}

	return err
}

const (
	registryAuthFile string = ".registry.auth"
)

type RegistryAuth struct {
	AuthConfigurations map[string]AuthConfiguration `json:"configs"`
}

type AuthConfiguration struct {
	Auth  string `json:"auth"`
	Email string `json:"email"`
}

func getDockerAuthConfigurationFromFile(registry string) *docker.AuthConfiguration {
	auth := loadRegistryAuthFile()
	dockerAuthConfiguration := &docker.AuthConfiguration{
		Username:      auth.Configs[registry].Username,
		Email:         auth.Configs[registry].Email,
		Password:      auth.Configs[registry].Password,
		ServerAddress: auth.Configs[registry].ServerAddress,
	}

	return dockerAuthConfiguration
}

func generateRegistryAuthFile(username, password, email, registry string) error {
	encodedCredential := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	c := &RegistryAuth{
		AuthConfigurations: map[string]AuthConfiguration{
			registry: {
				Auth:  encodedCredential,
				Email: email,
			},
		},
	}

	jsonView, _ := json.Marshal(c)
	err := ioutil.WriteFile(registryAuthFile, jsonView, 0600)
	if err != nil {
		logrus.Fatal(err)
	}

	return err
}

func loadRegistryAuthFile() *docker.AuthConfigurations {
	content, _ := os.Open(registryAuthFile)
	auth, err := docker.NewAuthConfigurations(content)
	if err != nil {
		logrus.Fatal(err)
	}
	return auth
}
