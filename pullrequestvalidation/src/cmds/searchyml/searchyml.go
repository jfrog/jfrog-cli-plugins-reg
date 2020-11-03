package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if !containOnlySingleYaml() {
		os.Exit(1)
	}
}

func containOnlySingleYaml() bool {
	pathToResource, commitSha := os.Getenv("res_validatePluginCriteria_resourcePath"), os.Getenv("res_validatePluginCriteria_commitSha")
	if pathToResource == "" || commitSha == "" {
		return false
	}
	os.Chdir(pathToResource)
	cmd := exec.Command("git", "diff-tree", "--no-commit-id", "--name-only", commitSha)
	output, err := cmd.Output()
	fmt.Println(string(output))
	if err != nil {
		return false
	}
	outputStr := strings.Trim(string(output), "\n")
	for i, commitFile := range strings.Split(outputStr, "\n") {
		if i > 0 || !strings.HasSuffix(commitFile, ".yml") {
			return false
		}
	}
	return true
}
