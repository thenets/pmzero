package lib

import (
	"log"
	"os/user"
	"runtime"
)

var configDirPath = getConfigDirPath()
var dataDirPath = getDataDirPath()
var stateFilePath = getDataDirPath() + "state.ini"

func getDataDirPath() string {
	if runtime.GOOS == "linux" {
		return "/var/lib/pmzero/"
	} else if runtime.GOOS == "windows" {
		return getUserDir() + "/.pmzero/"
	} else {
		log.Fatal("[ERROR] Operational system not supported")
	}

	return ""
}

func getConfigDirPath() string {
	if runtime.GOOS == "linux" {
		return "/etc/pmzero/"
	} else if runtime.GOOS == "windows" {
		return getUserDir() + "/.pmzero/"
	} else {
		log.Fatal("[ERROR] Operational system not supported")
	}

	return ""
}

func getUserDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}
