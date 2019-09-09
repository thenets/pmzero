package lib

import (
	"log"
	"path/filepath"
	"strings"
)

// https://ini.unknwon.io/docs/intro/getting_started

func init() {
	CheckBootTime()
	updateDeploymentsState()
}

func updateDeploymentsState() {
	var state = GetState()

	files, err := filepath.Glob(configDirPath + "*")
	if err != nil {
		log.Fatal(err)
	}

	var deploymentNames []string

	var deploymentPrefix = configDirPath + "deployment_"

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
			state.Section("deployments").Key(f.Name()).SetValue("false")
		}
	}

	state.SaveTo(stateFilePath)

}
