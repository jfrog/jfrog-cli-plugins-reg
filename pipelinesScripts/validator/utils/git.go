package utils

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	jfrogCliPluginRegUrl    = "https://github.com/jfrog/jfrog-cli-plugins-reg.git"
	jfrogCliPluginRegBranch = "master"
)

func getGitCloneFlags(branch, tag string) (flags string) {
	if branch != "" {
		flags = flags + "--branch=" + branch
		return
	}
	if tag != "" {
		flags = flags + "--branch=" + tag
	}
	return
}

// Clone the plugin's repository to a local temp directory and return the full path pointing to the plugin's code relative path.
// 'tempDir' - Temporary folder to which the project will be cloned.
// 'gitRepository' - The GitHub repository to clone.
// 'relativePath' - Relative path in the repository to be chained in the returned path.
// 'branch' - If specified, override the default branch with the input branch.
// 'tag' - If specified, checkout to the input tag.
// returns: (project-path, error)
func CloneRepository(tempDir, gitRepository, relativePath, branch, tag string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", errors.New("Couldn't get current directory: " + err.Error())
	}
	err = os.Chdir(tempDir)
	if err != nil {
		return "", errors.New("Couldn't get change directory to" + tempDir + ": " + err.Error())
	}
	defer os.Chdir(currentDir)
	gitRepository = strings.TrimSuffix(gitRepository, ".git")
	flags := getGitCloneFlags(branch, tag)
	if err := RunCommand("git", "clone", flags, gitRepository+".git"); err != nil {
		return "", errors.New("Failed to run git clone for " + gitRepository + ", error:" + err.Error())
	}
	repositoryName := gitRepository[strings.LastIndex(gitRepository, "/")+1:]

	return filepath.Join(tempDir, repositoryName, relativePath), nil
}

func GetModifiedFiles() ([]string, error) {
	// Create unique upstream and branch names
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	uniqueUpstream := "remote-origin-" + timestamp

	// Add remote.
	if err := RunCommand("git", "remote", "add", uniqueUpstream, jfrogCliPluginRegUrl); err != nil {
		return nil, errors.New("Failed to add git remote for " + uniqueUpstream + " upstream and" + jfrogCliPluginRegUrl + " branch. Error: " + err.Error())
	}
	defer RunCommand("git", "remote", "remove", uniqueUpstream)

	// Fetch from upsream
	if err := RunCommand("git", "fetch", uniqueUpstream, jfrogCliPluginRegBranch); err != nil {
		return nil, errors.New("Failed to fetch from " + uniqueUpstream + ", branch " + jfrogCliPluginRegBranch + ". Error: " + err.Error())
	}
	return runGitDiff(uniqueUpstream + "/master")
}

func runGitDiff(compareTo string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--no-commit-id", "--name-only", "-r", compareTo+"...HEAD")
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.New("Failed to run git diff command: " + err.Error())
	}
	var fullPathCommittedFiles []string
	for _, file := range strings.Split(string(output), "\n") {
		if file != "" {
			fullPathCommittedFiles = append(fullPathCommittedFiles, file)
		}
	}
	return fullPathCommittedFiles, nil
}
