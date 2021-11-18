package git

import (
	"fmt"
	"io/ioutil"
	"os"
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

	for _, command := range commands {
		_, _ := bash.Run(command)
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
		gitLog string
	}{
		{
			name: "main ahead of feat1",
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
			name: "feat1 ahead of feat1.1",
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
			name: "no changes",
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
	}

	for _, test := range cases {
		err := os.Chdir(test.path)
		if err != nil {
			panic(err)
		}

		commands, err := Rebase(test.path)
		if err != nil {
			panic(err)
		}

		for _, command := range commands {
			bash.Run(command)
		}
		stdout, _ := bash.Run(`git log --graph --pretty=format:%d.%s --all`)
		assert.Equal(t, test.gitLog, stdout, test.name)
	}
}
