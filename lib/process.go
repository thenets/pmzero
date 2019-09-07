package lib

import (
	"fmt"
	"log"
	"os/exec"
)

// CreateProcess and add it to a pull
func CreateProcess(commandName string, args []string) int {

	dateCmd := exec.Command(commandName, args...)
	out, err := dateCmd.Output()
	if err != nil {
		fmt.Printf("%v", string(out))
		log.Fatalf("[ERROR] During the command invoke: '%v %v'\n%v\n", commandName, args, err)
	}

	fmt.Printf("%v", string(out))

	return 0
}
