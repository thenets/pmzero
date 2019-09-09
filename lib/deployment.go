package lib

import (
	"errors"
	"io/ioutil"
	"log"

	"strconv"

	"gopkg.in/yaml.v2"
)

// DeploymentData configs from YAML file
type DeploymentData struct {
	Type string   `yaml:"type"`
	Name string   `yaml:"name"`
	CMD  []string `yaml:"cmd"`
	CPU  struct {
		Limit int `yaml:"limit"`
	}
	Linux struct {
		User string `yaml:"user"`
	}
}

// ReadDeploymentFile and returns the DeploymentData populated.
func ReadDeploymentFile(filePath string) DeploymentData {
	// Read config file content
	configFileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("[ERROR] Config file can't be read: ./%v\n%v", filePath, err)
	}

	// Create config object
	t := DeploymentData{}
	err = yaml.Unmarshal([]byte(configFileContent), &t)
	if err != nil {
		log.Fatalf("[ERROR] DeploymentData struct can't be populated:\n%v", err)
	}
	// fmt.Printf("--- t:\n%v\n\n", t)

	return t
}

// LoadDeploymentFile load deployment file and validates.
// Returns error message if file doesn't exist or is invalid.
func LoadDeploymentFile(filePath string) error {

	// Validates deployment file
	var err = validateDeploymentFile(filePath)
	if err != nil {
		log.Fatalf("[ERROR] Invalid config file.\n%v", err)
	}

	var d = ReadDeploymentFile(filePath)

	// Copy deployment file to config dir
	CopyFile(filePath, configDirPath+"deployment_"+d.Name+".yaml")

	return nil
}

func validateDeploymentFile(filePath string) error {
	var deploymentData = ReadDeploymentFile(filePath)

	// Check if has name
	if len(deploymentData.Name) == 0 {
		return errors.New("deployment: invalid deployment name")
	}

	// Check if already exist in state file
	if HasDeployment(deploymentData.Name) {
		return errors.New("deployment: already exist")
	}

	return nil
}

// HasDeployment if deployment does exist.
func HasDeployment(name string) bool {
	var state = GetState()

	var deploymentName = state.Section("deployments").Key(name).String()
	if len(deploymentName) == 0 && deploymentName != "false" {
		return false
	}

	return true
}

// StartDeployment if is not running
func StartDeployment(deploymentName string) {
	// TODO check if deployment exists

	var d = ReadDeploymentFile(configDirPath + "deployment_" + deploymentName + ".yaml")
	var state = GetState()

	// Create process
	var pid = CreateProcess(d.CMD[0], d.CMD[1:])
	state.Section("processes").Key(d.Name).SetValue(strconv.Itoa(pid))
	state.SaveTo(stateFilePath)
}

// AddDeployment file if doesn't exist and raise
// error if already exist.
func AddDeployment(filepath string) {

}
