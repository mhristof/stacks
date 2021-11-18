package bash

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// Run Run the bash command and return the stdout and stderr.
func Run(command string) (stdout string, stderr string, err error) {
	var stdoutB, stderrB bytes.Buffer

	fields := strings.Fields(command)
	cmd := exec.Command(fields[0], fields[1:]...)
	cmd.Stdout = &stdoutB
	cmd.Stderr = &stderrB
	err = cmd.Run()
	stdout, stderr = stdoutB.String(), stderrB.String()

	if err != nil {
		return stdout, stderr, errors.Wrap(err, "command failed")
	}

	return stdout, stderr, nil
}
