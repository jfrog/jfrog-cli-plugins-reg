package utils

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	JfrogCliPluginRegUrl    = "https://github.com/jfrog/jfrog-cli-plugins-reg.git"
	JfrogCliPluginRegBranch = "master"
)

func getGitCloneFlags(branch, tag string) []string {
	flags := []string{"clone"}
	if branch != "" {
		flags = append(flags, "--branch="+branch)
		return flags
	}
	if tag != "" {
		flags = append(flags, "--branch="+tag)
	}
	return flags
}

// Clone the plugin's repository to a local temp directory and return the full path of the plugin's source code.
// 'tempDir' - Temporary dir to which the project will be cloned.
// 'gitRepository' - The GitHub repository to clone.
// 'relativePath' - Relative path inside the repository.
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
	defer func() {
		if deferErr := os.Chdir(currentDir); deferErr != nil {
			log.Print("Failed to change dir to " + currentDir + ". Error:" + deferErr.Error())
		}
	}()
	gitRepository = strings.TrimSuffix(gitRepository, ".git")
	flags := append(getGitCloneFlags(branch, tag), gitRepository+".git")
	if err := RunCommand("git", flags...); err != nil {
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
	if err := RunCommand("git", "remote", "add", uniqueUpstream, JfrogCliPluginRegUrl); err != nil {
		return nil, errors.New("Failed to add git remote for " + uniqueUpstream + " upstream and" + JfrogCliPluginRegUrl + " branch. Error: " + err.Error())
	}
	defer func() {
		if deferErr := RunCommand("git", "remote", "remove", uniqueUpstream); deferErr != nil {
			log.Print("Failed to remove remote upstream " + uniqueUpstream + ". Error:" + deferErr.Error())
		}
	}()
	// Fetch from upstream
	if err := RunCommand("git", "fetch", uniqueUpstream, JfrogCliPluginRegBranch); err != nil {
		return nil, errors.New("Failed to fetch from " + uniqueUpstream + ", branch " + JfrogCliPluginRegBranch + ". Error: " + err.Error())
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
