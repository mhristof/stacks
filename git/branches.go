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
	ErrorDefaultBranchNotFound = errors.New("default branch not found")
)

func Branches(path string) (ret []string) {
	all := allbranches(path)

	return all
}

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

	sort.Sort(ByLength(ret))

	return ret
}

type ByLength []string

func (s ByLength) Len() int {
	return len(s)
}

func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByLength) Less(i, j int) bool {
	return (s[i]) < (s[j])
}
