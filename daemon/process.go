package lib

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	proc "github.com/shirou/gopsutil/process"
)

// createProcess and add it to a pull
func createProcess(deployment DeploymentData) int {
	var commandName = deployment.CMD[0]
	var args = deployment.CMD[1:]
	var err error

	// TODO Make sure the process logs files exist

	// Set log files
	stdoutFile, err := os.Create(dataDirPath + "./logs/" + deployment.Name)
	if err != nil {
		panic(err)
	}
	defer stdoutFile.Close()
	stderrFile, err := os.Create(dataDirPath + "./logs/" + deployment.Name)
	if err != nil {
		panic(err)
	}
	defer stderrFile.Close()

	// Convert env data
	var envs = os.Environ()
	for _, e := range deployment.Env {
		envs = append(envs, e.Name+"="+e.Value)

		// If Windows, set OS env
		if runtime.GOOS == "windows" {
			os.Setenv(e.Name, e.Value)
		}
	}
	cmd := exec.Command(commandName, args...)
	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile
	cmd.Dir = deployment.Workdir
	cmd.Env = envs
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("Just ran subprocess %d.\n", a.Cyan(cmd.Process.Pid))

	// Append exit code to the log file
	defer appendExitCode(cmd.ProcessState.ExitCode(), deployment)

	return int(cmd.Process.Pid)
}

func appendExitCode(exitCode int, deployment DeploymentData) {
	AppendToFile(dataDirPath+"./logs/"+deployment.Name, "Exited "+strconv.Itoa(exitCode))
}

// AppendToFile appends "text" to the file
func AppendToFile(filePath string, text string) {
	f, err := os.OpenFile(filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(text + "\n"); err != nil {
		log.Println(err)
	}
}

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
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
