package lib

import (
	"fmt"

	"github.com/shirou/gopsutil/host"
)

func init() {
	fmt.Println(host.BootTime())
}
