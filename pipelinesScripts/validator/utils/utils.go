package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ValidationType string

const (
	Extension ValidationType = "extension"
	Structure                = "structure"
	Tests                    = "unit_tests"
)

func PrintUsageAndExit() {
	fmt.Printf("Usage: `go run validator.go <command>`\nPossible commands: '%s', '%s' or '%s'\n", Extension, Structure, Tests)
	os.Exit(1)
}

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
	cmd := exec.Command("git", "clone", flags, gitRepository+".git")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return "", errors.New("Failed to run git clone for " + gitRepository + ", error:" + err.Error())
	}
	repositoryName := gitRepository[strings.LastIndex(gitRepository, "/")+1:]
	if tag != "" {
		err = os.Chdir(repositoryName)
		if err != nil {
			return "", errors.New("Fail to get change directory to" + repositoryName + ", error:" + err.Error())
		}
		cmd := exec.Command("git", "checkout", "tags/"+tag)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			return "", errors.New("Failed to checkout tag" + tag + ", error:" + err.Error())
		}
	}
	return filepath.Join(tempDir, repositoryName, relativePath), nil
}

// Return the paths to the modified files for all affected files since master's commit.
func GetModifiedFiles() ([]string, error) {
	pathToResource, commitSha := os.Getenv("res_jfrog_cli_plugins_reg_resourcePath"), os.Getenv("res_jfrog_cli_plugins_reg_commitSha")
	if pathToResource == "" || commitSha == "" {
		return nil, errors.New("Failed to parse env vars: res_jfrog_cli_plugins_reg_resourcePath & res_jfrog_cli_plugins_reg_commitSha")
	}
	os.Chdir(pathToResource)
	cmd := exec.Command("git", "diff", "--no-commit-id", "--name-only", "-r", "master..."+commitSha)
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.New("Failed to run git cmd, error:" + err.Error())
	}
	var fullPathCommittedFiles []string
	for _, committedFile := range strings.Split(string(output), "\n") {
		if committedFile != "" {
			fullPathCommittedFiles = append(fullPathCommittedFiles, pathToResource+"/"+committedFile)
		}
	}
	return fullPathCommittedFiles, nil
}
