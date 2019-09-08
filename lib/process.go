package lib

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	a "github.com/logrusorgru/aurora"
	proc "github.com/shirou/gopsutil/process"
)

// CreateProcess and add it to a pull
func CreateProcess(commandName string, args []string) int {
	cmd := exec.Command("test")
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Just ran subprocess %d.\n", a.Cyan(cmd.Process.Pid))

	time.Sleep(5 * time.Second)

	p, err := proc.NewProcess(int32(cmd.Process.Pid))
	if err != nil {
		log.Fatalf("[ERROR] Trying to get the process PID.\n%v", err)
	}
	p.Kill()
	fmt.Println(p)

	return 0
}

// SimpleCommandInvoker and add it to a pull
func SimpleCommandInvoker(commandName string, args []string) int {

	dateCmd := exec.Command(commandName, args...)
	out, err := dateCmd.Output()
	if err != nil {
		fmt.Printf("%v", string(out))
		log.Fatalf("[ERROR] During the command invoke: '%v %v'\n%v\n", commandName, args, err)
	}

	fmt.Printf("%v", string(out))

	return 0
}
