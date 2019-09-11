package lib

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/host"
	proc "github.com/shirou/gopsutil/process"
)

func init() {
	createConfigDir()
	UpdateState()
}

// UpdateState checks all changes and update the state file
func UpdateState() {
	CheckBootTime()
	updateDeploymentsState()
	updateProcessState()
	RefactorAllDeploymentFile()
}

func createConfigDir() {
	// Create cache/ dir if not exist
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		os.Mkdir(configDirPath, os.ModePerm)
	}

	// Create cache/logs/ dir if not exist
	if _, err := os.Stat(configDirPath + "logs"); os.IsNotExist(err) {
		os.Mkdir(configDirPath + "logs", os.ModePerm)
	}

	// Create state.ini if not exist
	if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
		emptyFile, err := os.Create(stateFilePath)
		if err != nil {
			log.Fatal(err)
		}
		emptyFile.Close()

		// Save current boot time
		bootTime, err := host.BootTime()
		if err != nil {
			log.Fatalf("[ERROR] Erro during get the host's boot time.\n%v", err)
		}
		var state = GetState()
		state.Section("").Key("boot_time").SetValue(strconv.Itoa(int(bootTime)))
		state.SaveTo(stateFilePath)
	}

}

func updateDeploymentsState() {
	var state = GetState()

	files, err := filepath.Glob(configDirPath + "*")
	if err != nil {
		log.Fatal(err)
	}

	var deploymentNames []string

	var deploymentPrefix = strings.ReplaceAll(configDirPath+"deployment_", "\\", "/")

	for _, f := range files {
		f = strings.ReplaceAll(f, "\\", "/")
		if strings.Contains(f, deploymentPrefix) {
			f = strings.ReplaceAll(f, deploymentPrefix, "")
			deploymentName := strings.Split(f, ".")[0]

			// Check if deploymentName exist in state file
			var tmp = state.Section("deployments").Key(deploymentName).String()
			if len(tmp) == 0 || tmp == "false" {
				state.Section("deployments").Key(deploymentName).SetValue("true")
			}

			deploymentNames = append(deploymentNames, deploymentName)
		}
	}

	var keys = state.Section("deployments").Keys()
	for _, f := range keys {
		if !StringInSlice(f.Name(), deploymentNames) {
			state.Section("deployments").DeleteKey(f.Name())
		}
	}

	state.SaveTo(stateFilePath)
}

func updateProcessState() {
	var state = GetState()

	for _, p := range state.Section("processes").Keys() {
		var pid, _ = strconv.ParseInt(p.Value(), 0, 0)
		var pidExist, _ = proc.PidExists(int32(pid))
		if !pidExist {
			state.Section("processes").DeleteKey(p.Name())
		}
	}

	state.SaveTo(stateFilePath)
}
