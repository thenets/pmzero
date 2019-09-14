package lib

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	proc "github.com/shirou/gopsutil/process"
)

// createProcess and add it to a pull
func createProcess(deployment DeploymentData) int {
	var commandName = deployment.CMD[0]
	var args = deployment.CMD[1:]
	var err error

	// Set log files
	stdoutFile, err := os.Create(configDirPath + "./logs/" + deployment.Name)
	if err != nil {
		panic(err)
	}
	defer stdoutFile.Close()
	stderrFile, err := os.Create(configDirPath + "./logs/" + deployment.Name)
	if err != nil {
		panic(err)
	}
	defer stderrFile.Close()

	// Convert env data
	var envs []string
	for _, e := range deployment.Env{
		envs = append(envs, e.Name+"="+e.Value)
	}

	cmd := exec.Command(commandName, args...)
	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile
	cmd.Env = envs
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("Just ran subprocess %d.\n", a.Cyan(cmd.Process.Pid))

	return int(cmd.Process.Pid)
}

// stopProcess returns nil if stop or already stopped
func stopProcess(pid int) error {
	if pid == -1 {
		return nil
	}

	p, err := proc.NewProcess(int32(pid))
	if err != nil {
		log.Fatalf("[ERROR] Trying to get the process PID.\n%v", err)
	}
	return p.Kill()
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
