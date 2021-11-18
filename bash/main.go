package bash

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func Run(command string) (stdout string, stderr string) {
	fields := strings.Fields(command)
	cmd := exec.Command(fields[0], fields[1:]...)
	var stdoutB, stderrB bytes.Buffer
	cmd.Stdout = &stdoutB
	cmd.Stderr = &stderrB
	err := cmd.Run()
	stdout, stderr = string(stdoutB.Bytes()), string(stderrB.Bytes())
	if err != nil {
		fmt.Println(fmt.Sprintf("command: %+v", command))
		fmt.Println(fmt.Sprintf("stdout: %+v", stdout))
		fmt.Println(fmt.Sprintf("stderr: %+v", stderr))
		panic(err)
	}
	return stdout, stderr
}
