package git

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	// ErrorDefaultBranchNotFound Error when the default branch is not found.
	ErrorDefaultBranchNotFound = errors.New("default branch not found")
)

// Branches Returns a list of all the branches lex sorted.
func Branches(path string) (ret []string) {
	all := allbranches(path)

	return all
}

// Branch Get the current branch or empty if it could not be found.
func Branch(path string) string {
	head := filepath.Join(path, ".git/HEAD")
	if _, err := os.Stat(head); os.IsNotExist(err) {
		return ""
	}

	branch, err := ioutil.ReadFile(head)
	if err != nil {
		return ""
	}

	return filepath.Base(strings.Fields(string(branch))[1])
}

// Rebase return a list of commands to run in order to rebase required branches.
func Rebase(path, branchName string) (ret []string, err error) {
	var branches = Branches(path)
	var branchRE = regexp.MustCompile(branchName)
	var onto string

	for _, branch := range branches {
		if branch == "main" || branch == "master" {
			onto = branch

			continue
		}
	}

	if onto == "" {
		return ret, ErrorDefaultBranchNotFound
	}

	for _, branch := range branches {
		if branch == "main" || branch == "master" {
			continue
		}

		if !branchRE.Match([]byte(branch)) {
			continue
		}

		ret = append(ret, fmt.Sprintf("git checkout %s", branch))
		ret = append(ret, fmt.Sprintf("git rebase --onto %s %s@{1}", onto, onto))
		onto = branch
	}

	return ret, nil
}

func allbranches(path string) (ret []string) {
	err := filepath.Walk(filepath.Join(path, ".git"),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !strings.Contains(path, ".git/refs/heads/") {
				return nil
			}

			ret = append(ret, filepath.Base(path))

			return nil
		})
	if err != nil {
		return
	}

	sort.Sort(byLength(ret))

	return ret
}

type byLength []string

func (s byLength) Len() int {
	return len(s)
}

func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byLength) Less(i, j int) bool {
	return (s[i]) < (s[j])
}
