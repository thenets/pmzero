package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"strconv"

	"gopkg.in/yaml.v2"
)

// DeploymentData configs from YAML file
type DeploymentData struct {
	Type    string   `yaml:"type"`
	Name    string   `yaml:"name"`
	CMD     []string `yaml:"cmd"`
	Workdir string
	Env     []struct {
		Name  string
		Value string
	}
	PID    int
	Status string
	CPU    struct {
		Limit int `yaml:"limit"`
	}
	Linux struct {
		User string `yaml:"user"`
	}
}

// GetDeploymentByFilePath and returns the DeploymentData populated.
func GetDeploymentByFilePath(filePath string) DeploymentData {
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
	// fmt.Printf("--- t:\n%v\n\n", t) // DEBUG

	// Add PID
	state := GetState()
	if state.Section("processes").HasKey(t.Name) {
		var pid, _ = strconv.ParseInt(state.Section("processes").Key(t.Name).String(), 0, 0)
		t.PID = int(pid)
	} else {
		t.PID = -1
	}

	// Add status
	if t.PID == -1 {
		t.Status = "stopped"
	} else {
		t.Status = "running"
	}

	return t
}

// RefactorDeploymentFile check all content and update values based on DeploymentData struct
func RefactorDeploymentFile(deploymentName string) {
	var deployment = GetDeploymentByName(deploymentName)

	// Recreate file based on struct
	var filePath = configDirPath + "deployment_" + deployment.Name + ".yaml"
	var configFileContent, err = yaml.Marshal(deployment)
	ioutil.WriteFile(filePath, configFileContent, 0644)
	if err != nil {
		log.Fatalf("[ERROR] Config file can't be write: ./%v\n%v", filePath, err)
	}
	yaml.Marshal(deployment)
}

// RefactorAllDeploymentFile check all content and update values based on DeploymentData struct of all files
func RefactorAllDeploymentFile() {
	var deployments = GetDeployments()

	for _, deployment := range deployments {
		RefactorDeploymentFile(deployment.Name)
	}
}

// LoadDeploymentFile copy deployment file to config dir and validates.
// Returns error message if file doesn't exist or is invalid.
func LoadDeploymentFile(filePath string) error {
	var err error

	// Validates deployment file
	err = validateDeploymentFile(filePath)
	if err != nil {
		log.Fatalf("[ERROR] Invalid config file.\n%v", err)
	}

	var d = GetDeploymentByFilePath(filePath)
	var newDeploymentFilePath = configDirPath + "deployment_" + d.Name + ".yaml"

	// Delete file if already exist.
	// Equivalent to file replace.
	if HasDeployment(d.Name) {
		fmt.Printf("deployment '%s' already exist. Updating...\n", d.Name)
		StopDeployment(d.Name)
		var err = os.Remove(newDeploymentFilePath)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		fmt.Printf("deployment '%s' Adding...\n", d.Name)
	}

	// Copy deployment file to config dir
	CopyFile(filePath, newDeploymentFilePath)

	// Refactor deployment file
	RefactorDeploymentFile(d.Name)

	return nil
}

// TailDeployment is equivalent to "tail -f" for all deployment output
func TailDeployment(deploymentName string) {
	var deployment = GetDeploymentByName(deploymentName)

	if runtime.GOOS == "windows" {
		// Windows
		// fmt.Println("You are running on Windows")
		cmd := exec.Command("powershell", "-c", "Get-Content", "-Path", "\""+configDirPath+"./logs/"+deployment.Name+"\"", "-Wait")
		cmd.Stderr = os.Stdout
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Println("[ERROR] ", err)
		}
	} else {
		// Linux of MacOS
		cmd := exec.Command("tail", "-f", configDirPath+"./logs/"+deployment.Name)
		cmd.Stderr = os.Stdout
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Println("[ERROR] ", err)
		}
		cmd.Wait()
	}

}

// GetDeploymentByName search the deployment by name and returns
// a DeploymentData populated.
func GetDeploymentByName(deploymentName string) DeploymentData {
	return GetDeploymentByFilePath(getDeplomentFilePath(deploymentName))
}

func getDeplomentFilePath(deploymentName string) string {
	return configDirPath + "deployment_" + deploymentName + ".yaml"
}

func validateDeploymentFile(filePath string) error {
	var deploymentData = GetDeploymentByFilePath(filePath)

	// Check if has name
	if len(deploymentData.Name) == 0 {
		return errors.New("deployment: invalid deployment name")
	}

	return nil
}

// HasDeployment if deployment does exist.
func HasDeployment(deploymentName string) bool {
	var state = GetState()

	var hasDeploymentString = state.Section("deployments").Key(deploymentName).String()
	if len(hasDeploymentString) == 0 && hasDeploymentString != "false" {
		return false
	}

	return true
}

// StartDeployment if is not running
func StartDeployment(deploymentName string) {
	var d = GetDeploymentByFilePath(configDirPath + "deployment_" + deploymentName + ".yaml")

	// Check if deployment exists
	UpdateState()
	var state = GetState()
	var pidFromStateFile, _ = strconv.ParseInt(state.Section("processes").Key(d.Name).String(), 0, 0)
	if pidFromStateFile > 0 {
		fmt.Println("process already running")
	} else {
		// Create process
		var pid = createProcess(d)
		state.Section("processes").Key(d.Name).SetValue(strconv.Itoa(pid))
		state.SaveTo(stateFilePath)
	}
}

// GetDeployments return an array of DeploymentData
func GetDeployments() []DeploymentData {
	var deployments []DeploymentData

	var state = GetState()

	for _, d := range state.Section("deployments").Keys() {
		if d.Value() == "true" {
			deployments = append(deployments, GetDeploymentByName(d.Name()))
		}
	}

	return deployments
}

// StopDeployment and returns nil if no error
func StopDeployment(deploymentName string) error {
	var d = GetDeploymentByName(deploymentName)
	return stopProcess(d.PID)
}

// RestartDeployment TODO
func RestartDeployment(deploymentName string) error {
	var err = StopDeployment(deploymentName)
	if err != nil {
		return err
	}

	StartDeployment(deploymentName)

	return nil
}

// DeleteDeployment TODO
func DeleteDeployment(deploymentName string) error {
	var err = os.Remove(getDeplomentFilePath(deploymentName))
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

// ForegroundDeployments keep running and respawn all deployments
func ForegroundDeployments() {
	for {
		time.Sleep(2 * time.Second)
		UpdateState()

		var deployments = GetDeployments()
		for _, d := range deployments {
			if d.PID == -1 {
				log.Printf("[%s] is down. Restarting...", d.Name)
				StartDeployment(d.Name)
			}
		}

		RefactorAllDeploymentFile()
	}
}
