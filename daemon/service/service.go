package service

import (
	"log"
	"os"
	"os/user"
	"runtime"

	"github.com/gobuffalo/packr/v2"
)

// SetupService setups service for all compatible OS
func SetupService() {
	if runtime.GOOS == "linux" {
		log.Println("Platform: Linux")
		checkLinuxPrivilegies()
		log.Println("Has root privilegies")
		copyLinuxServiceFile()
		log.Println("Service file copied")
		enableLinuxService()
		log.Println("Service enabled")
	} else {
		log.Fatalln("[ERROR] Operational system not supported.")
	}
}

func checkLinuxPrivilegies() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	if user.Uid != "0" {
		log.Fatalln("[ERROR] This program must be run as root! (sudo)")
	}
}

func enableLinuxService() {
	// TODO
}

func copyLinuxServiceFile() {
	box := getBoxOfStaticFiles()

	// Systemd file
	systemdFileContent, err := box.FindString("linux-systemd.ini")
	if err != nil {
		log.Fatal(err)
	}
	systemdFile, err := os.Create("/etc/systemd/system/pmzero.service")
	if err != nil {
		log.Fatal(err)
	}
	_, err = systemdFile.WriteString(systemdFileContent)
	if err != nil {
		log.Fatal(err)
	}

	// Init file
	initFileContent, err := box.FindString("linux-init.sh")
	if err != nil {
		log.Fatal(err)
	}
	initFile, err := os.Create("/etc/init.d/pmzero")
	if err != nil {
		log.Fatal(err)
	}
	_, err = initFile.WriteString(initFileContent)
	if err != nil {
		log.Fatal(err)
	}
	err = initFile.Chmod(0755)
	if err != nil {
		log.Fatal(err)
	}

}

func getBoxOfStaticFiles() *packr.Box {
	return packr.New("servicesTemplates", "./templates")
}
