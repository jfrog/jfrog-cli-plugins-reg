package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jfrog/jfrog-cli-plugins-reg/utils"
)

// This program validate a new jfrog cli plugin register
func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: Wrong number of arguments.")
		utils.PrintUsageAndExit()
	}

	command := os.Args[1]
	var err error
	switch strings.ToLower(command) {
	case string(utils.Extension):
		err = validateStructure()
	case string(utils.Structure):
		err = validateDescriptor()
	case string(utils.Tests):
		err = runTests()
	default:
		err = errors.New("Unknown command: " + command)
	}
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}
}

// In order to add a plugin to the registry,
// the maintainer should create a pull request to the registry.
// The pull request should include the plugin(s) YAML.
// If other files extentions exists, return an error.
func validateStructure() error {
	prFiles, err := utils.GetModifiedFiles()
	if err != nil {
		return err
	}
	forbiddenFiles := ""
	for _, committedFilePath := range prFiles {
		if !strings.HasSuffix(committedFilePath, ".yml") || !strings.HasPrefix(committedFilePath, utils.PluginDescriptoPrefix) {
			forbiddenFiles += committedFilePath + "\n"
		}
	}
	if forbiddenFiles != "" {
		return errors.New("Only .yml files are permitted to be in the pull request. Please remove: " + forbiddenFiles)
	}
	return nil
}

// Check the plugin YAML file format. if one of the mandatory fields are missing, return an error.
func validateDescriptor() error {
	files, err := utils.GetModifiedFiles()
	if err != nil {
		return err
	}
	for _, yamlPath := range files {
		log.Print("Validating:" + yamlPath)

		descriptor, err := utils.ReadDescriptor(yamlPath)
		if err != nil {
			return err
		}

		if err := validateContent(descriptor); err != nil {
			return err
		}
	}
	return nil
}

// Verifies the plugin and run the plugin tests using 'go test ./...'.
func runTests() error {
	files, err := utils.GetModifiedFiles()
	if err != nil {
		return err
	}
	for _, yamlPath := range files {
		fmt.Println("Analyzing:" + yamlPath)

		descriptor, err := utils.ReadDescriptor(yamlPath)
		if err != nil {
			return err
		}
		tempDir, err := ioutil.TempDir("", "pluginRepo")
		if err != nil {
			return errors.New("Fail to create temp dir: " + err.Error())
		}
		defer os.RemoveAll(tempDir)
		projectPath, err := utils.CloneRepository(tempDir, descriptor.Repository, descriptor.RelativePath, descriptor.Branch, descriptor.Tag)
		if err != nil {
			return err
		}
		if err := runProjectTests(projectPath); err != nil {
			return err
		}
	}
	return nil
}

func runProjectTests(projectPath string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return errors.New("Fail to get current directory: " + err.Error())
	}
	defer os.Chdir(currentDir)
	err = os.Chdir(projectPath)
	if err != nil {
		return errors.New("Fail to get change directory to" + projectPath + ": " + err.Error())
	}
	cmd := exec.Command("go", "vet", "-v", "./...")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.New("Failed plugin lint for " + projectPath + ": " + err.Error())
	}

	cmd = exec.Command("go", "test", "-v", "./...")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.New("Failed plugin tests for " + projectPath + ": " + err.Error())
	}
	return nil
}

func validateContent(descriptor *utils.PluginDescriptor) error {
	missingfields := ""
	if descriptor.PluginName == "" {
		missingfields += "* 'name' is missing\n"
	}
	if descriptor.Version == "" {
		missingfields += "* 'version' is missing\n"
	}
	if descriptor.Repository == "" {
		missingfields += "* 'repository' is missing\n"
	}
	if len(descriptor.Maintainers) == 0 {
		missingfields += "* 'maintainers' is missing\n"
	}
	if descriptor.Tag != "" && descriptor.Branch != "" {
		missingfields += "* Plugin descriptor yml cannot include both 'tag' and 'branch'.\n"
	}
	if missingfields != "" {
		return errors.New("Errors detected in the yml descriptor file:\n" + missingfields)
	}
	return nil
}
