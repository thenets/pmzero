package lib

import (
	"fmt"
	"os"

	"github.com/shirou/gopsutil/host"
	"gopkg.in/ini.v1"
)

func init() {
	fmt.Println(host.BootTime())

	cfg, err := ini.Load("cache/data.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

    // Classic read of values, default section can be represented as empty string
    fmt.Println("App Mode:", cfg.Section("").Key("app_mode").String())
    fmt.Println("Data Path:", cfg.Section("paths").Key("data").String())

	// Now, make some changes and save it
	cfg.Section("").Key("app_mode").SetValue("production")
	cfg.SaveTo("cache/data.ini")

}
