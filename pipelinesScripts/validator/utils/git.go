package utils

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	jfrogCliPluginRegUrl    = "https://github.com/jfrog/jfrog-cli-plugins-reg"
	jfrogCliPluginRegBranch = "master"
)

func getGitCloneFlags(branch string) (flags string) {
	if branch != "" {
		flags = flags + "--branch=" + branch
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
	flags := getGitCloneFlags(branch)
	if err := RunCommand("git", "clone", flags, gitRepository+".git"); err != nil {
		return "", errors.New("Failed to run git clone for " + gitRepository + ", error:" + err.Error())
	}
	repositoryName := gitRepository[strings.LastIndex(gitRepository, "/")+1:]
	if tag != "" {
		err = os.Chdir(repositoryName)
		if err != nil {
			return "", errors.New("Fail to get change directory to" + repositoryName + ", error:" + err.Error())
		}
		if err := RunCommand("git", "checkout", "tags/"+tag); err != nil {
			return "", errors.New("Failed to checkout tag" + tag + ", error:" + err.Error())
		}
	}
	return filepath.Join(tempDir, repositoryName, relativePath), nil
}

func GetModifiedFiles() ([]string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, errors.New("Couldn't get current directory: " + err.Error())
	}
	defer os.Chdir(currentDir)

	os.Chdir(filepath.Join("..", "..", "plugins"))

	// Add upstream remote
	if err := RunCommand("git", "remote", "add", "upstream", jfrogCliPluginRegUrl); err != nil {
		return nil, errors.New("Failed to add git remote for " + jfrogCliPluginRegUrl + ": " + err.Error())
	}

	// Fetch from upsream
	if err := RunCommand("git", "fetch", "upstream"); err != nil {
		return nil, errors.New("Failed to fetch from " + jfrogCliPluginRegUrl + ": " + err.Error())
	}

	// Checkout to a new JFrog branch
	if err := RunCommand("git", "checkout", "-b", "jfrog-master", "upstream/"+jfrogCliPluginRegBranch); err != nil {
		return nil, errors.New("Checkout failed to 'upstream/" + jfrogCliPluginRegBranch + "': " + err.Error())
	}

	// Merge changes from JFrog branch to the current
	if err := RunCommand("git", "merge", "jfrog-master"); err != nil {
		return nil, errors.New("Failed to merge changes from '" + jfrogCliPluginRegUrl + "/" + jfrogCliPluginRegBranch + "': " + err.Error())
	}

	return runGitDiff()
}

func runGitDiff() ([]string, error) {
	cmd := exec.Command("git", "diff", "--no-commit-id", "--name-only", "-r", "jfrog-master", "HEAD")
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
