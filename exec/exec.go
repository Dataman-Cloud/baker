package exec

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// Blueprint is a project specification file.
type Blueprint struct {
	CompileBuild CompileBuild `yaml:"compilebuild"`
	BinaryBuild  BinaryBuild  `yaml:"binarybuild"`
	Canary       Canary       `yaml:"canary"`
	Rollback     Rollback     `yaml:"rollback"`
}

// CompileBuild type.
type CompileBuild struct {
	BaseImage string            `yaml:"image"`
	Command   string            `yaml:"command"`
	EnvVars   map[string]string `yaml:"envVars"`
}

// BinaryBuild type.
type BinaryBuild struct {
	BaseImage   string `yaml:"image"`
	LibraryFile string `yaml:"libraryfile"`
}

// Canary type.
type Canary struct {
	Rule string `yaml:"rule"`
}

// Rollback type.
type Rollback struct {
	Trigger string `yaml:"trigger"`
}

// Job type is a baker job.
type Job interface {
	//Build    Build    // Build job including compilebuild and binarybuild.
	//Canary   Canary   // Canary to execute rolling upgrade.
	//Rollback rollback // Rollback to execute rollback publish.
}

// New is the constructor of Job.
func NewJob(yamlFile string) *Job {
	// retrieve blueprint
	blueprint, err := getBlueprint(yamlFile)
	if err != nil {
		log.Printf("Error getting the project specification. %s", err)
		return err
	}
	return &Job{}
}

// ExecuteBuild is a worker to execute a build job.
func (j *Job) ExecuteBuild() error {
	log.Printf("Build Job to be executed: ")

	return nil
}

// ExecuteCanary is a worker to execute a rolling upgrade job.
func (j *Job) ExecuteCanary() error {
	log.Printf("Rolling Upgrade Job to be executed: ")

	return nil
}

// ExecuteRollback is a worker to execute a rollback job.
func (j *Job) ExecuteRollback() error {
	log.Printf("Rollback Job to be executed: ")

	return nil
}

// parseBlueprint to retrieve .blueprint.yml.
func parseBlueprint(yamlFile string) (*Buildprint, error) {
	contents, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %s", yamlFile, err)
	}
	bp := &Blueprint{}
	if err = yaml.Unmarshal(contents, bp); err != nil {
		return nil, fmt.Errorf("could not parse blueprint file: %s", err)
	}
	return bp, nil
}
