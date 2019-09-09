package lib

import (
	"log"
	"os/user"
)

var configDirPath = getUserDir() + "/.pmzero/"
var stateFilePath = configDirPath + "state.ini"

func getUserDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}
