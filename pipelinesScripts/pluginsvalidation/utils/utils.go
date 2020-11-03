package utils

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getGitCloneFlags(branch string) (flags string) {
	if branch != "" {
		flags = flags + "--branch=" + branch
	}
	return
}

// Clone the plugin's project to a local temp dir.
// 'tempDir' Temporary folder to which the project will be copiedץ
// 'branch' override the default pointing branch in the cloned project.
// 'tag' override the defalt pointing commit in the cloned project.
// 'relativePath' adds a specific path with in the cloned repo e.g:
// repo=my-plugins, relativePath=pluginA -> my-plugins/pluginA as the project-path.
// returns: (project-path, error)
func CloneProject(tempDir, projectRepository, relativePath, branch, tag string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", errors.New("Fail to get current directory, error:" + err.Error())
	}
	err = os.Chdir(tempDir)
	if err != nil {
		return "", errors.New("Fail to get change directory to" + tempDir + ", error:" + err.Error())
	}
	defer os.Chdir(currentDir)
	projectRepository = strings.TrimSuffix(projectRepository, ".git")
	flags := getGitCloneFlags(branch)
	cmd := exec.Command("git", "clone", flags, projectRepository+".git")
	if _, err := cmd.Output(); err != nil {
		return "", errors.New("Failed to run git clone for " + projectRepository + ", error:" + err.Error())
	}
	repositoryName := projectRepository[strings.LastIndex(projectRepository, "/")+1:]
	if tag != "" {
		err = os.Chdir(repositoryName)
		if err != nil {
			return "", errors.New("Fail to get change directory to" + repositoryName + ", error:" + err.Error())
		}
		cmd := exec.Command("git", "checkout", "tags/"+tag)
		if _, err := cmd.Output(); err != nil {
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
