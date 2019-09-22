package lib

import (
	"log"
	"os/user"
	"runtime"
)

var configDirPath = getConfigDirPath()
var stateFilePath = configDirPath + "state.ini"

func getConfigDirPath() string {
	if runtime.GOOS == "linux" {
		return "/etc/pmzero/"
	} else {
		return getUserDir() + "/.pmzero/"
	}
}

func getUserDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}
