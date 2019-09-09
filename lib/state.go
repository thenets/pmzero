package lib

import (
	"log"
	"math"
	"os"
	"strconv"

	"github.com/shirou/gopsutil/host"
	"gopkg.in/ini.v1"
)

// CheckBootTime update boot time in state.ini and delete all processes
// from "processes" section if boot time changes.
func CheckBootTime() {
	var cfg = GetState()
	bootTimeTmp, err := host.BootTime()
	if err != nil {
		log.Fatalf("[ERROR] Erro during get the host's boot time.\n%v", err)
	}
	bootTimeFromDataFileTmp, err := strconv.ParseInt(cfg.Section("").Key("boot_time").String(), 0, 0)
	if err != nil {
		log.Println("[WARNING] Can't parse boot time from state.ini file as integer.")
		bootTimeFromDataFileTmp = 0
	}
	bootTime := float64(bootTimeTmp)
	bootTimeFromDataFile := float64(bootTimeFromDataFileTmp)

	if math.Abs(bootTimeFromDataFile-bootTime) > 4 {
		deleteProcesses()
		cfg = GetState()

		// Save current boot time
		cfg.Section("").Key("boot_time").SetValue(strconv.Itoa(int(bootTime)))
		cfg.SaveTo(stateFilePath)
	}
}

func deleteProcesses() {
	var cfg = GetState()
	cfg.DeleteSection("processes")
	cfg.SaveTo(stateFilePath)
}

// GetState return the ini object.
func GetState() *ini.File {
	// Create cache/ dir if not exist
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		os.Mkdir(configDirPath, os.ModeDir)
	}

	// Create state.ini if not exist
	if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
		emptyFile, err := os.Create(stateFilePath)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(emptyFile)
		emptyFile.Close()
	}

	// Load state.ini file as cfg
	cfg, err := ini.Load(stateFilePath)
	if err != nil {
		log.Fatalf("[ERROR] Fail to read file: %v\n", err)
		os.Exit(1)
	}

	return cfg
}

// StringInSlice returns if string exist in slice.
func StringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}