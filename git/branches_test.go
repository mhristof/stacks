package git

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mkgit(commands []string) string {
	dir, err := ioutil.TempDir("", "git")
	if err != nil {
		panic(err)
	}

	err = os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("dir: %+v", dir))

	for _, command := range commands {
		fields := strings.Fields(command)
		cmd := exec.Command(fields[0], fields[1:]...)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

		if err != nil {
			fmt.Println(fmt.Sprintf("outStr: %+v", outStr))

			fmt.Println(fmt.Sprintf("errStr: %+v", errStr))
			panic(err)
		}
	}
	return dir
}

func TestBranches(t *testing.T) {
	var cases = []struct {
		name     string
		branches []string
		path     string
	}{
		{
			name: "non git folder",
			path: mkgit([]string{}),
		},
		{
			name: "git folder with main",
			path: mkgit([]string{
				"git init",
				"git commit --allow-empty -m 'empty.commit'",
			}),
			branches: []string{
				"main",
			},
		},
		{
			name: "git folder with a couple of branches",
			path: mkgit([]string{
				"git init",
				"git commit --allow-empty -m 'empty.commit'",
				"git checkout -b foobar",
			}),
			branches: []string{
				"foobar", "main",
			},
		},
	}

	for _, test := range cases {
		assert.Equal(t, test.branches, Branches(test.path), test.name)
		defer os.Remove(test.path)
	}
}
