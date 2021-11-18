package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/mhristof/go-stacks/bash"
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

	fmt.Println(fmt.Sprintf("commands: %+v", commands))
	//fmt.Println(fmt.Sprintf("commands[0]: %+v", commands[0]))

	for _, command := range commands {
		fmt.Println(fmt.Sprintf("command: %+v", command))

		_, _ = bash.Run(command)
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
				"git checkout -b foobar1",
			}),
			branches: []string{
				"foobar", "foobar1", "main",
			},
		},
	}

	for _, test := range cases {
		assert.Equal(t, test.branches, Branches(test.path), test.name)
		defer os.Remove(test.path)
	}
}

func TestRebase(t *testing.T) {
	var cases = []struct {
		name   string
		path   string
		branch string
		gitLog string
	}{
		{
			name:   "main ahead of feat1",
			branch: ".*",
			path: mkgit([]string{
				"git init",
				"git commit --allow-empty -m 'empty.commit'",
				"git checkout -b feat1",
				"git commit --allow-empty -m 'feat1.commit'",
				"git checkout main",
				"git commit --allow-empty -m 'main.commit'",
			}),
			gitLog: heredoc.Doc(`
				*  (HEAD -> feat1).'feat1.commit'
				*  (main).'main.commit'
				* .'empty.commit'`),
		},
		{
			name:   "feat1 ahead of feat1.1",
			branch: ".*",
			path: mkgit([]string{
				"git init",
				"git commit --allow-empty -m 'empty.commit'",
				"git commit --allow-empty -m 'empty.commit1'",
				"git checkout -b feat1",
				"git commit --allow-empty -m 'feat1.commit'",
				"git checkout -b feat1.1",
				"git commit --allow-empty -m 'feat1.1.commit'",
				"git checkout feat1",
				"git commit --allow-empty -m 'feat1.commit2'",
			}),
			gitLog: heredoc.Doc(`
				*  (HEAD -> feat1.1).'feat1.1.commit'
				*  (feat1).'feat1.commit2'
				* .'feat1.commit'
				* .'empty.commit1'
				*  (main).'empty.commit1'
				* .'empty.commit'`),
		},
		{
			name:   "no changes",
			branch: ".*",
			path: mkgit([]string{
				"git init",
				"git commit --allow-empty -m 'empty.commit'",
				"git commit --allow-empty -m 'empty.commit1'",
				"git checkout -b feat1",
				"git commit --allow-empty -m 'feat1.commit'",
				"git checkout -b feat1.1",
				"git commit --allow-empty -m 'feat1.1.commit'",
			}),
			gitLog: heredoc.Doc(`
				*  (HEAD -> feat1.1).'feat1.1.commit'
				*  (feat1).'feat1.commit'
				* .'empty.commit1'
				*  (main).'empty.commit1'
				* .'empty.commit'`),
		},
		{
			name:   "limit branch name",
			branch: "feat1",
			path: mkgit(strings.Split(heredoc.Doc(`
				git init
				git commit --allow-empty -m 'empty.commit'
				git commit --allow-empty -m 'empty.commit1'
				git checkout -b feat1
				git commit --allow-empty -m 'feat1.commit'
				git checkout -b feat1.1
				git commit --allow-empty -m 'feat1.1.commit'
				git checkout feat1
				git commit --allow-empty -m 'feat1.commit1'
				git checkout main
				git checkout -b feat2
				git commit --allow-empty -m 'feat2.commit1'
				git checkout -b feat2.1
				git commit --allow-empty -m 'feat2.1.commit1'
				git checkout feat2
				git commit --allow-empty -m 'feat2.commit2'`),
				"\n")),
			gitLog: heredoc.Doc(`
				*  (HEAD -> feat1.1).'feat1.1.commit'
				*  (feat1).'feat1.commit1'
				* .'feat1.commit'
				* .'empty.commit1'
				| *  (feat2).'feat2.commit2'
				| | *  (feat2.1).'feat2.1.commit1'
				| |/  
				| * .'feat2.commit1'
				|/  
				*  (main).'empty.commit1'
				* .'empty.commit'`),
		},
	}

	for _, test := range cases {
		err := os.Chdir(test.path)
		if err != nil {
			panic(err)
		}

		commands, err := Rebase(test.path, test.branch)
		if err != nil {
			panic(err)
		}

		for _, command := range commands {
			fmt.Println(fmt.Sprintf("command: %+v", command))

			bash.Run(command)
		}
		stdout, _ := bash.Run(`git log --graph --pretty=format:%d.%s --all`)
		assert.Equal(t, test.gitLog, stdout, test.name)
		defer os.Remove(test.path)
	}
}
