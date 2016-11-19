package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config type.
type Config struct {
	Server         ServerConfig         `yaml:"server"`
	Users          map[string]Authorize `yaml:"users"`
	DockerRegistry DockerRegistry       `yaml:"dockerRegistry"`
}

// ServerConfig type.
type ServerConfig struct {
	ListenAddress string `yaml:"addr,omitempty"`
}

// Password type.
type Authorize struct {
	Password string `yaml:"password"`
}

// dockerRegistry type.
type DockerRegistry struct {
	Address  string `yaml:"address"`
	Repo     string `yaml:"repo"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Email    string `yaml:"email"`
}

func validate(c *Config) error {
	return nil
}

// Decode the configuration in a file path.
func Decode(path string) (*Config, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %s", path, err)
	}
	c := &Config{}
	if err = yaml.Unmarshal(contents, c); err != nil {
		return nil, fmt.Errorf("could not parse config: %s", err)
	}
	if err = validate(c); err != nil {
		return nil, fmt.Errorf("invalid config: %s", err)
	}
	return c, nil
}
